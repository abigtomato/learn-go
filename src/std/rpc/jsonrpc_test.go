package rpc

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"testing"
)

// Panda 协议规定建立Service结构体
type PandaStruct struct{}

// Args 协议规定被远程调用函数的参数必须封装为结构体
type Args struct {
	A, B int
}

// Div 1. 协议规定可被远程调用的函数必须是结构体的方法
// 2. 协议规定传入参数必须是结构体，函数结果必须是指针类型
func (p *PandaStruct) Div(args Args, result *float64) error {
	if args.B == 0 {
		return errors.New("division by zero")
	}
	*result = float64(args.A) / float64(args.B)
	return nil
}

func TestJsonRpcServer(t *testing.T) {
	// 注册rpc服务
	if err := rpc.Register(new(PandaStruct)); err != nil {
		log.Println(err.Error())
		return
	}

	// 通过tcp协议，监听本机1234端口的连接请求
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		panic(err)
	}

	for {
		// 接收连接
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err.Error())
			continue
		}

		// 开启go程处理连接
		go jsonrpc.ServeConn(conn)
	}
}

func TestJsonRpcClient(t *testing.T) {
	// 发送连接请求
	conn, err := net.Dial("tcp", "127.0.0.1:1234")
	if err != nil {
		log.Println(err.Error())
	}
	defer func(conn net.Conn) {
		_ = conn.Close()
	}(conn)

	// 创建jsonRpc客户端
	client := jsonrpc.NewClient(conn)
	defer func(client *rpc.Client) {
		_ = client.Close()
	}(client)

	var result float64
	// 执行远程函数(通过"结构体.方法名"的方式调用)
	if err := client.Call("Panda.Div", Args{A: 10, B: 3}, &result); err != nil {
		log.Println(err.Error())
		return
	}
	fmt.Printf("result: %v\n", result)

	if err := client.Call("Panda.Div", Args{A: 10, B: 0}, &result); err != nil {
		log.Println(err.Error())
	}
	fmt.Printf("result: %v\n", result)
}
