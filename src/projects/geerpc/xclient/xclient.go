package xclient

import (
	"context"
	"io"
	"learn-go/src/projects/geerpc/client"
	"learn-go/src/projects/geerpc/server"
	"reflect"
	"sync"
)

// XClient 负载均衡客户端
type XClient struct {
	d       Discovery                 // 服务发现
	mode    SelectMode                // 负载均衡策略
	opt     *server.Option            // 协议选项
	mu      sync.Mutex                // 互斥锁
	clients map[string]*client.Client // 创建成功的Client实例
}

var _ io.Closer = (*XClient)(nil)

func NewXClient(d Discovery, mode SelectMode, opt *server.Option) *XClient {
	return &XClient{
		d:       d,
		mode:    mode,
		opt:     opt,
		clients: make(map[string]*client.Client),
	}
}

func (xc *XClient) Close() error {
	xc.mu.Lock()
	defer xc.mu.Unlock()
	for key, c := range xc.clients {
		_ = c.Close()
		delete(xc.clients, key)
	}
	return nil
}

func (xc *XClient) dial(rpcAddr string) (*client.Client, error) {
	xc.mu.Lock()
	defer xc.mu.Unlock()
	// 检查是否有缓存的client
	c, ok := xc.clients[rpcAddr]
	if ok && !c.IsAvailable() {
		// 不可以，从缓存中删除
		_ = c.Close()
		delete(xc.clients, rpcAddr)
		c = nil
	}
	// 创建新的client并缓存
	if c == nil {
		var err error
		c, err := client.XDial(rpcAddr, xc.opt)
		if err != nil {
			return nil, err
		}
		xc.clients[rpcAddr] = c
	}
	return c, nil
}

func (xc *XClient) Call(rpcAddr string, ctx context.Context, serviceMethod string, args, reply any) error {
	c, err := xc.dial(rpcAddr)
	if err != nil {
		return err
	}
	return c.CallTimeout(ctx, serviceMethod, args, reply)
}

// Broadcast 请求广播
func (xc *XClient) Broadcast(ctx context.Context, serviceMethod string, args, reply any) error {
	servers, err := xc.d.GetAll()
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	var mu sync.Mutex
	var e error
	replyDone := reply == nil
	// 借助 context.WithCancel 确保有错误发生时，快速失败
	ctx, cancel := context.WithCancel(ctx)
	for _, rpcAddr := range servers {
		wg.Add(1)
		go func(rpcAddr string) {
			defer wg.Done()
			var cloneReply any
			if reply != nil {
				cloneReply = reflect.New(reflect.ValueOf(reply).Elem().Type()).Interface()
			}
			err := xc.Call(rpcAddr, ctx, serviceMethod, args, cloneReply)
			mu.Lock()
			if err != nil && e == nil {
				e = err
				cancel()
			}
			if err == nil && !replyDone {
				reflect.ValueOf(reply).Elem().Set(reflect.ValueOf(cloneReply).Elem())
				replyDone = true
			}
			mu.Unlock()
		}(rpcAddr)
	}
	wg.Wait()
	cancel()
	return e
}
