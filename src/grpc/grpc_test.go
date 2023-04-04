package grpc

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pt "learn-go/src/grpc/proto"
	"log"
	"net"
	"testing"
)

// 根据proto内定义的服务创建rpc服务
type HelloServer struct{}

// 实现pb.go内定义的rpc服务接口
func (s *HelloServer) SayHello(_ context.Context, in *pt.HelloRequest) (*pt.HelloReplay, error) {
	return &pt.HelloReplay{
		Message: "Hello " + in.Name,
	}, nil
}

func (s *HelloServer) GetHelloMsg(context.Context, *pt.HelloRequest) (*pt.HelloMessage, error) {
	return &pt.HelloMessage{
		Msg: "s is from server",
	}, nil
}

func TestServer(t *testing.T) {
	// tcp连接监听
	ln, err := net.Listen("tcp", ":18881")
	if err != nil {
		log.Println(err.Error())
		return
	}

	// 创建一个grpc服务端句柄
	srv := grpc.NewServer()
	// 将自定义rpc服务注册到grpc句柄中
	pt.RegisterHelloServerServer(srv, new(HelloServer))

	// 开启grpc服务的监听
	if err := srv.Serve(ln); err != nil {
		log.Println(err.Error())
		return
	}
}

func TestClient(t *testing.T) {
	// 客户端发起grpc请求连接服务器
	conn, err := grpc.Dial("127.0.0.1:18881", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer func(conn *grpc.ClientConn) {
		_ = conn.Close()
	}(conn)

	// 创建一个grpc客户端句柄
	cli := pt.NewHelloServerClient(conn)

	// 执行grpc远程调用
	replay, err := cli.SayHello(context.Background(), &pt.HelloRequest{Name: "panda"})
	if err != nil {
		log.Println(err.Error())
		return
	}
	fmt.Printf("HelloServer SayHello HelloReplay: %v\n", replay.Message)

	msg, err := cli.GetHelloMsg(context.Background(), &pt.HelloRequest{Name: "panda"})
	if err != nil {
		log.Println(err.Error())
		return
	}
	fmt.Printf("HelloServer GetHelloMsg HelloMessage: %v\n", msg.Msg)
}
