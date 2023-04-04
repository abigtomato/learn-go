package rpc

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"testing"
)

// Panda 注册rpc服务的必须是自定义类型
type Panda int

// GetInfo
// rpc服务提供的远程过程必须是以下固定格式:
// func (t *T) MethodName(argType T1, replyType *T2) error
// 1. 第一个参数是服务调用时传入的参数
// 2. 第二个参数是服务调用方主机的内存块的指针
// 3. 只存在一个返回值error
func (p *Panda) GetInfo(argType int, replyType *int) error {
	log.Println(argType)
	// 操作的是服务调用方主机的内存
	*replyType = argType + 100
	return nil
}

func TestRpcServer(t *testing.T) {
	// 注册一个页面请求的处理器
	http.HandleFunc("/panda", func(w http.ResponseWriter, r *http.Request) {
		if _, err := io.WriteString(w, "panda"); err != nil {
			log.Println(err.Error())
		}
	})

	// 注册rpc服务
	if err := rpc.Register(new(Panda)); err != nil {
		log.Println(err.Error())
		return
	}
	rpc.HandleHTTP()

	// 创建网络监听
	ln, err := net.Listen("tcp", ":10086")
	if err != nil {
		log.Println(err.Error())
		return
	}

	// 开启http服务，接收监听器接收的http连接请求
	if err := http.Serve(ln, nil); err != nil {
		log.Println(err.Error())
		return
	}
}

func TestRpcClient(t *testing.T) {
	// 与服务端建立tcp连接
	cli, err := rpc.DialHTTP("tcp", "127.0.0.1:10086")
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer func() {
		if err := cli.Close(); err != nil {
			log.Println(err.Error())
			return
		}
	}()

	// rpc远程调用函数
	var val int
	if err := cli.Call("Panda.GetInfo", 123, &val); err != nil {
		log.Println(err.Error())
		return
	}
	fmt.Printf("rpc result: %v\n", val)
}
