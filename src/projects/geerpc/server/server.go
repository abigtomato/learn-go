package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"learn-go/src/projects/geerpc/codec"
	"learn-go/src/projects/geerpc/service"
	"log"
	"net"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"time"
)

const MagicNumber = 0x3bef5c

// Option 表示GeeRPC协议选项位的抽象
// GeeRPC 客户端固定采用 JSON 编码 Option，后续的 header 和 body 的编码方式由 Option 中的 CodeType 指定
// 服务端首先使用 JSON 解码 Option，然后通过 Option 的 CodeType 解码剩余的内容
// GeeRPC报文的格式：
// | Option{MagicNumber: xxx, CodecType: xxx} | Header{ServiceMethod ...} | Body interface{} |
// | <------      固定 JSON 编码      ------>  | <-------   编码方式由 CodeType 决定   ------->|
// 在一次连接中，Option 固定在报文的最开始，Header 和 Body 可以有多个：
// | Option | Header1 | Body1 | Header2 | Body2 | ...
type Option struct {
	MagicNumber    int           // 魔数 GeeRPC请求的标记
	CodecType      codec.Type    // GeeRPC的消息编码器类型
	ConnectTimeout time.Duration // 连接超时
	HandleTimeout  time.Duration // 处理超时
}

// DefaultOption 默认选项位
var DefaultOption = &Option{
	MagicNumber:    MagicNumber,
	CodecType:      codec.GobType,
	ConnectTimeout: time.Second * 10,
}

// 存储一次RPC调用的所有信息
type request struct {
	h            *codec.Header       // 请求头
	argv, replyv reflect.Value       // 请求参数和返回值
	mType        *service.MethodType // 方法类型
	svc          *service.Service    // 相关服务
}

// 无效的请求 是发生错误时响应 argv 的占位符
var invalidRequest = struct{}{}

// Server RPC服务器抽象
type Server struct {
	serviceMap    sync.Map      // 服务列表
	handleTimeout time.Duration // 请求处理超时
}

func (s *Server) ServiceMap() *sync.Map {
	return &s.serviceMap
}

// NewServer 创建服务器实例
func NewServer(handleTimeout time.Duration) *Server {
	return &Server{handleTimeout: handleTimeout}
}

// DefaultServer 默认服务器
var DefaultServer = NewServer(time.Second * 10)

// Accept 监听传入的请求并建立连接
func (s *Server) Accept(lis net.Listener) {
	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Println("rpc server: accept error:", err)
		}
		go s.ServeConn(conn)
	}
}

// Accept 默认服务器的连接请求处理
func Accept(lis net.Listener) {
	DefaultServer.Accept(lis)
}

// ServeConn 处理连接事件
func (s *Server) ServeConn(conn io.ReadWriteCloser) {
	defer func(conn io.ReadWriteCloser) {
		_ = conn.Close()
	}(conn)

	// 接收并解码客户端发送的选项位
	var opt Option
	if err := json.NewDecoder(conn).Decode(&opt); err != nil {
		log.Println("rpc server: option error: ", err)
		return
	}

	// 检查魔数
	if opt.MagicNumber != MagicNumber {
		log.Printf("rpc server: invalid magic number %x", opt.MagicNumber)
	}

	// 获取对应的消息编解码器
	f := codec.NewCodecFuncMap[opt.CodecType]
	if f == nil {
		log.Printf("rpc server: invalid codec type %s", opt.CodecType)
	}

	s.serveCodec(f(conn), opt.HandleTimeout)
}

// 请求的编解码
func (s *Server) serveCodec(cc codec.Codec, handleTimeout time.Duration) {
	// 请求的处理是并发的，但请求响应是逐个的，需要使用锁保证
	sending := new(sync.Mutex)
	// 等待所有请求都得到处理
	wg := new(sync.WaitGroup)
	// 在一次连接中，允许接收多个请求，即多个 request header 和 request body
	// 直到发生错误（例如连接被关闭，接收到的报文有问题等）
	for {
		// 读取请求
		req, err := s.readRequest(cc)
		if err != nil {
			if req == nil {
				break
			}
			req.h.Error = err.Error()
			// 响应错误
			s.sendResponse(cc, req.h, invalidRequest, sending)
			continue
		}
		wg.Add(1)
		// 处理请求
		go s.handleRequest(cc, req, sending, wg, handleTimeout)
	}
	wg.Wait()
	_ = cc.Close()
}

// 读取请求
func (s *Server) readRequest(cc codec.Codec) (*request, error) {
	// 读取header
	h, err := s.readRequestHeader(cc)
	if err != nil {
		return nil, err
	}

	// 获取对应的服务和方法
	req := &request{h: h}
	req.svc, req.mType, err = s.findService(h.ServiceMethod)
	if err != nil {
		return req, nil
	}

	// 创建出入参和返回值的实例
	req.argv = req.mType.NewArgv()
	req.replyv = req.mType.NewReplyv()

	// 读取body 参数解析
	argv := req.argv.Interface()
	if req.argv.Type().Kind() != reflect.Pointer {
		// 指针类型的处理方式
		argv = req.argv.Addr().Interface()
	}
	if err = cc.ReadBody(argv); err != nil {
		log.Println("rpc server: read body err:", err)
		return req, err
	}
	return req, nil
}

// 读取请求头部
func (s *Server) readRequestHeader(cc codec.Codec) (*codec.Header, error) {
	var h codec.Header
	if err := cc.ReadHeader(&h); err != nil {
		if err != io.EOF && err != io.ErrUnexpectedEOF {
			log.Println("rpc server: read header error:", err)
		}
		return nil, err
	}
	return &h, nil
}

// 处理请求
func (s *Server) handleRequest(cc codec.Codec, req *request, sending *sync.Mutex, wg *sync.WaitGroup, handleTimeout time.Duration) {
	defer wg.Done()

	// 方法调用阶段的信道标记
	called := make(chan struct{})
	// 发送响应阶段的信道标记
	sent := make(chan struct{})

	go func() {
		err := req.svc.Call(req.mType, req.argv, req.replyv)
		called <- struct{}{}
		if err != nil {
			req.h.Error = err.Error()
			s.sendResponse(cc, req.h, invalidRequest, sending)
			sent <- struct{}{}
			return
		}
		s.sendResponse(cc, req.h, req.replyv.Interface(), sending)
		sent <- struct{}{}
	}()

	// 若超时时间未设置，等待两个阶段执行完毕退出
	if handleTimeout == 0 {
		if s.handleTimeout == 0 {
			<-called
			<-sent
			return
		}
		handleTimeout = s.handleTimeout
	}

	select {
	// 超时则响应错误
	case <-time.After(handleTimeout):
		req.h.Error = fmt.Sprintf("rpc server: request handle timeout: expect within %s", handleTimeout)
		s.sendResponse(cc, req.h, invalidRequest, sending)
	// 未超时，则先等待调用执行完毕，再等待响应执行完后退出
	case <-called:
		<-sent
	}
}

// 发送响应
func (s *Server) sendResponse(cc codec.Codec, h *codec.Header, body any, sending *sync.Mutex) {
	sending.Lock()
	defer sending.Unlock()
	if err := cc.Write(h, body); err != nil {
		log.Println("rpc server: write response error:", err)
	}
}

// Register 服务注册
func (s *Server) Register(rcvr any) error {
	svc := service.NewService(rcvr)
	if _, dup := s.serviceMap.LoadOrStore(svc.Name, svc); dup {
		return errors.New("rpc: service already defined: " + svc.Name)
	}
	return nil
}

// Register 注册服务到默认服务器上
func Register(rcvr any) error {
	return DefaultServer.Register(rcvr)
}

// 获取服务
func (s *Server) findService(serviceMethod string) (svc *service.Service, mtype *service.MethodType, err error) {
	// 分割 "服务名.方法名"
	dot := strings.LastIndex(serviceMethod, ".")
	if dot < 0 {
		err = errors.New("rpc server: service/method request ill-formed: " + serviceMethod)
		return
	}
	serviceName, methodName := serviceMethod[:dot], serviceMethod[dot+1:]

	// 根据服务名获取服务
	svcAny, ok := s.serviceMap.Load(serviceName)
	if !ok {
		err = errors.New("rpc server: can't find service " + serviceName)
		return
	}
	svc = svcAny.(*service.Service)

	// 获取具体的方法类型
	mtype = svc.Method[methodName]
	if mtype == nil {
		err = errors.New("rpc server: can't find method " + methodName)
	}

	return
}

const (
	Connected        = "200 Connected to Gee RPC" // 连接建立后响应的内容
	DefaultRPCPath   = "/_geerpc"                 // HTTP请求RPC资源路径
	DefaultDebugPath = "/debug/geerpc"            // HTTP请求RPC调试路径
)

// 支持HTTP协议，通信过程：
// 1. 客户端向RPC服务器发送CONNECT请求
// 2. RPC服务器返回200表示连接建立
// 3. 客户端使用连接发送RPC报文，先发送Option，再发送N个请求报文，服务端处理RPC请求并响应
func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != "CONNECT" {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = io.WriteString(w, "405 must CONNECT\n")
		return
	}
	conn, _, err := w.(http.Hijacker).Hijack()
	if err != nil {
		log.Print("rpc hijacking ", req.RemoteAddr, ": ", err.Error())
		return
	}
	_, _ = io.WriteString(conn, "HTTP/1.0"+Connected+"\n\n")
	s.ServeConn(conn)
}

// HandleHTTP 注册HTTP资源路径处理器
func (s *Server) HandleHTTP() {
	// 协议转换
	http.Handle(DefaultRPCPath, s)
	// 调试路径
	http.Handle(DefaultDebugPath, debugHTTP{s})
	log.Println("rpc server debug path:", DefaultDebugPath)
}

func HandleHTTP() {
	DefaultServer.HandleHTTP()
}
