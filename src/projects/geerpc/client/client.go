package client

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"learn-go/src/projects/geerpc/codec"
	"learn-go/src/projects/geerpc/server"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

// Call 表示客户端的一次RPC调用
type Call struct {
	Seq           uint64     // 序列号
	ServiceMethod string     // 调用服务和方法 格式为 "<service>.<method>"
	Args          any        // 参数
	Reply         any        // 返回值
	Error         error      // 错误
	Done          chan *Call // 异步调用完成标记
}

func (c *Call) done() {
	c.Done <- c
}

// Client RPC客户端
// 可能存在多个未完成的RPC调用
// 使用单个客户端，可以由多go程处理
type Client struct {
	cc       codec.Codec      // 消息的编解码器 用于序列化发送的请求和反序列化服务的响应
	opt      *server.Option   // 协议的选项位 用于表示通信的使用的协议和编码方式
	sending  sync.Mutex       // 保证并非环境下请求的有序发送
	header   codec.Header     // 请求的消息头
	mu       sync.Mutex       // 互斥锁
	seq      uint64           // 请求的唯一编号
	pending  map[uint64]*Call // 存储未处理完的请求 key是编号value指向Call实例
	closing  bool             // 用户主动关闭标识
	shutdown bool             // 内部错误关闭标识
}

var _ io.Closer = (*Client)(nil)

var ErrShutdown = errors.New("connection is shut down")

// NewClient 创建客户端实例 完成协议的交换和编码方式
func NewClient(conn net.Conn, opt *server.Option) (*Client, error) {
	f := codec.NewCodecFuncMap[opt.CodecType]
	if f == nil {
		err := fmt.Errorf("invalid codec type %s", opt.CodecType)
		log.Println("rpc client: codec error:", err)
		return nil, err
	}
	// 编码并发送选项位给服务端
	if err := json.NewEncoder(conn).Encode(opt); err != nil {
		log.Println("rpc client: options error:", err)
		_ = conn.Close()
		return nil, err
	}
	return newClientCodec(f(conn), opt), nil
}

// 更据此次协议交换商定的编解码器创建客户端实例
func newClientCodec(cc codec.Codec, opt *server.Option) *Client {
	client := &Client{
		cc:      cc,
		opt:     opt,
		seq:     1,
		pending: make(map[uint64]*Call),
	}
	// 开启go程处理服务端响应
	go client.receive()
	return client
}

// 选项位解析 默认值填充
func parseOptions(opts ...*server.Option) (*server.Option, error) {
	if len(opts) == 0 || opts[0] == nil {
		return server.DefaultOption, nil
	}
	if len(opts) != 1 {
		return nil, errors.New("number of option is more than 1")
	}
	opt := opts[0]
	opt.MagicNumber = server.DefaultOption.MagicNumber
	if opt.CodecType == "" {
		opt.CodecType = server.DefaultOption.CodecType
	}
	return opt, nil
}

// Dial 客户端拨号函数 便于创建客户端实例
func Dial(network, address string, opts ...*server.Option) (client *Client, err error) {
	return DialTimeout(NewClient, network, address, opts...)
}

type clientResult struct {
	client *Client
	err    error
}

type newClientFunc func(conn net.Conn, opt *server.Option) (client *Client, err error)

func DialTimeout(f newClientFunc, network, address string, opts ...*server.Option) (Client *Client, err error) {
	opt, err := parseOptions(opts...)
	if err != nil {
		return nil, err
	}

	// 具备超时能力的客户端拨号
	conn, err := net.DialTimeout(network, address, opt.ConnectTimeout)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			_ = conn.Close()
		}
	}()

	ch := make(chan clientResult)
	go func() {
		// go程处理客户端创建
		client, err := f(conn, opt)
		// 创建完后通过信道发送结果
		ch <- clientResult{
			client: client,
			err:    err,
		}
	}()

	if opt.ConnectTimeout == 0 {
		result := <-ch
		return result.client, result.err
	}

	select {
	// 若是time.After先收到消息，说明客户端创建超时
	case <-time.After(opt.ConnectTimeout):
		return nil, fmt.Errorf("rpc client: connect timeout: expect within %s", opt.ConnectTimeout)
	// 接收客户端创建结果，正常返回
	case result := <-ch:
		return result.client, result.err
	}
}

func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closing {
		return ErrShutdown
	}
	c.closing = true
	return c.cc.Close()
}

func (c *Client) IsAvailable() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return !c.shutdown && !c.closing
}

// 向客户端注册一次RPC调用
func (c *Client) registerCall(call *Call) (uint64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closing || c.shutdown {
		return 0, ErrShutdown
	}
	call.Seq = c.seq
	c.pending[call.Seq] = call
	c.seq++
	return call.Seq, nil
}

// 从客户端中移除某个RPC调用并返回
func (c *Client) removeCall(seq uint64) *Call {
	c.mu.Lock()
	defer c.mu.Unlock()
	call := c.pending[seq]
	delete(c.pending, seq)
	return call
}

// 发生错误时调用
func (c *Client) terminateCalls(err error) {
	c.sending.Lock()
	defer c.sending.Unlock()
	c.mu.Lock()
	defer c.mu.Unlock()
	c.shutdown = true
	// 向客户端中所有待定的调用发送完成标记
	for _, call := range c.pending {
		call.Error = err
		call.done()
	}
}

// 接收服务端响应
func (c *Client) receive() {
	var err error
	for err == nil {
		var h codec.Header
		if err = c.cc.ReadHeader(&h); err != nil {
			break
		}
		call := c.removeCall(h.Seq)
		switch {
		// call 不存在，可能是请求没有发送完整，或者因为其他原因被取消，但是服务端仍旧处理了
		case call == nil:
			err = c.cc.ReadBody(nil)
		// call 存在，但服务端处理出错，即 h.Error 不为空
		case h.Error != "":
			call.Error = fmt.Errorf(h.Error)
			err = c.cc.ReadBody(nil)
			call.done()
		// call 存在，服务端处理正常，那么需要从 body 中读取 Reply 的值
		default:
			err = c.cc.ReadBody(call.Reply)
			if err != nil {
				call.Error = errors.New("reading body " + err.Error())
			}
			call.done()
		}
	}
	// 发生错误，终止调用
	c.terminateCalls(err)
}

// 发送RPC请求
func (c *Client) send(call *Call) {
	// 确保客户端发送完整的请求
	c.sending.Lock()
	defer c.sending.Unlock()

	// 向客户端注册此次的RPC调用
	seq, err := c.registerCall(call)
	if err != nil {
		call.Error = err
		call.done()
		return
	}

	c.header.ServiceMethod = call.ServiceMethod
	c.header.Seq = seq
	c.header.Error = ""

	// 发送数据
	if err := c.cc.Write(&c.header, call.Args); err != nil {
		call := c.removeCall(seq)
		if call != nil {
			call.Error = err
			call.done()
		}
	}
}

// Go 用于异步调用的发送函数
func (c *Client) Go(serviceMethod string, args, reply any, done chan *Call) *Call {
	if done == nil {
		done = make(chan *Call, 10)
	} else if cap(done) == 0 {
		log.Panic("rpc client: done channel is unbuffered")
	}
	call := &Call{
		ServiceMethod: serviceMethod,
		Args:          args,
		Reply:         reply,
		Done:          done,
	}
	c.send(call)
	return call
}

func (c *Client) Call(serviceMethod string, args, reply any) error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelFunc()
	return c.CallTimeout(ctx, serviceMethod, args, reply)
}

// CallTimeout 用于同步调用的发送函数
// 调用者可以使用context.WithTimeout 创建具备超时检测能力的 context 对象来控制
// 如：ctx, _ := context.WithTimeout(context.Background(), time.Second)
func (c *Client) CallTimeout(ctx context.Context, serviceMethod string, args, reply any) error {
	call := c.Go(serviceMethod, args, reply, make(chan *Call, 1))
	select {
	// 同步调用超时检查
	case <-ctx.Done():
		c.removeCall(call.Seq)
		return errors.New("rpc client: call failed: " + ctx.Err().Error())
	case call := <-call.Done:
		return call.Error
	}
}

// NewHTTPClient 支持HTTP协议
func NewHTTPClient(conn net.Conn, opt *server.Option) (*Client, error) {
	// CONNECT请求
	_, _ = io.WriteString(conn, fmt.Sprintf("CONNECT %s HTTP/1.0\n\n", server.DefaultRPCPath))

	resp, err := http.ReadResponse(bufio.NewReader(conn), &http.Request{Method: "CONNECT"})
	if err == nil && resp.Status == server.Connected {
		// 通过 HTTP CONNECT 请求建立连接后，后续的通信过程就交给 NewClient
		return NewClient(conn, opt)
	}
	if err == nil {
		err = errors.New("unexpected HTTP response: " + resp.Status)
	}

	return nil, err
}

func DialHTTP(network, address string, opts ...*server.Option) (*Client, error) {
	return DialTimeout(NewHTTPClient, network, address, opts...)
}

// XDial 简化调用的统一入口
func XDial(rpcAddr string, opts ...*server.Option) (*Client, error) {
	parts := strings.Split(rpcAddr, "@")
	if len(parts) != 2 {
		return nil, fmt.Errorf("rpc client err: wrong format '%s', expect protocol@addr", rpcAddr)
	}
	protocol, addr := parts[0], parts[1]
	switch protocol {
	case "http":
		return DialHTTP("tcp", addr, opts...)
	default:
		return Dial(protocol, addr, opts...)
	}
}
