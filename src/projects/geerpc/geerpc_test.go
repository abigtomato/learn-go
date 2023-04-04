package geerpc

import (
	"Golearn/src/projects/geerpc/client"
	"Golearn/src/projects/geerpc/server"
	"Golearn/src/projects/geerpc/xclient"
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"
)

func _assert(condition bool, msg string, v ...any) {
	if !condition {
		panic(fmt.Sprintf("assertion failed: "+msg, v...))
	}
}

func startServer(addr chan string) {
	listen, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal("network error:", err)
	}
	log.Println("start rpc server on", listen.Addr())
	addr <- listen.Addr().String()
	server.Accept(listen)
}

func TestServer(t *testing.T) {
	log.SetFlags(0)

	addr := make(chan string)
	go startServer(addr)

	c, _ := client.Dial("tcp", <-addr)
	defer func(c *client.Client) {
		_ = c.Close()
	}(c)

	time.Sleep(time.Second)

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			args := fmt.Sprintf("geerpc req %d", i)
			var reply string
			if err := c.Call("Foo.Sum", args, &reply); err != nil {
				log.Fatal("call Foo.Sum error:", err)
			}
			log.Println("reply:", reply)
		}(i)
	}
	wg.Wait()
}

type Foo int

type Args struct {
	Num1, Num2 int
}

func (f Foo) Sum(args Args, reply *int) error {
	*reply = args.Num1 + args.Num2
	return nil
}

func startRpcServer(addr chan string) {
	var foo Foo
	if err := server.Register(&foo); err != nil {
		log.Fatal("register error:", err)
	}
	listen, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal("network error:", err)
	}
	log.Println("start rpc server on", listen.Addr())
	addr <- listen.Addr().String()
	server.Accept(listen)
}

func TestRpcCall(t *testing.T) {
	log.SetFlags(0)
	addr := make(chan string)
	go startRpcServer(addr)
	c, _ := client.Dial("tcp", <-addr)
	defer func() { _ = c.Close() }()

	time.Sleep(time.Second)

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			args := &Args{Num1: i, Num2: i * i}
			var reply int
			if err := c.Call("Foo.Sum", args, &reply); err != nil {
				log.Fatal("call Foo.Sum error:", err)
			}
			log.Printf("%d + %d = %d", args.Num1, args.Num2, reply)
		}(i)
	}
	wg.Wait()
}

func TestClientDialTimeout(t *testing.T) {
	t.Parallel()
	listen, _ := net.Listen("tcp", ":0")

	f := func(conn net.Conn, opt *server.Option) (client *client.Client, err error) {
		_ = conn.Close()
		time.Sleep(time.Second * 2)
		return nil, nil
	}

	t.Run("timeout", func(t *testing.T) {
		_, err := client.DialTimeout(f, "tcp", listen.Addr().String(), &server.Option{ConnectTimeout: time.Second})
		_assert(err != nil && strings.Contains(err.Error(), "connect timeout"), "expect a timeout error")
	})

	t.Run("0", func(t *testing.T) {
		_, err := client.DialTimeout(f, "tcp", listen.Addr().String(), &server.Option{ConnectTimeout: 0})
		_assert(err == nil, "0 means no limit")
	})
}

type Bar int

func (b Bar) Timeout(argv int, reply *int) error {
	time.Sleep(time.Second * 2)
	return nil
}

func startDefaultTimeoutServer(addr chan string) {
	var b Bar
	_ = server.Register(b)
	listen, _ := net.Listen("tcp", ":0")
	addr <- listen.Addr().String()
	server.Accept(listen)
}

func TestTimeoutCall(t *testing.T) {
	t.Parallel()
	addrCh := make(chan string)
	go startDefaultTimeoutServer(addrCh)
	addr := <-addrCh

	time.Sleep(time.Second)

	t.Run("client timeout", func(t *testing.T) {
		c, _ := client.Dial("tcp", addr)
		ctx, _ := context.WithTimeout(context.Background(), time.Second)
		var reply int
		err := c.CallTimeout(ctx, "Bar.Timeout", 1, &reply)
		_assert(err != nil && strings.Contains(err.Error(), ctx.Err().Error()), "expect a timeout error")
	})

	t.Run("server handle timeout", func(t *testing.T) {
		c, _ := client.Dial("tcp", addr, &server.Option{HandleTimeout: time.Second})
		var reply int
		err := c.CallTimeout(context.Background(), "Bar.Timeout", 1, &reply)
		_assert(err != nil && strings.Contains(err.Error(), "handle timeout"), "expect a timeout error")
	})
}

func TestXDial(t *testing.T) {
	if runtime.GOOS == "linux" {
		ch := make(chan struct{})
		addr := "/tmp/geerpc.sock"
		go func() {
			_ = os.Remove(addr)
			listen, err := net.Listen("tcp", addr)
			if err != nil {
				t.Fatal("failed to listen unix socket")
			}
			ch <- struct{}{}
			server.Accept(listen)
		}()
		<-ch
		_, err := client.XDial("unix@" + addr)
		_assert(err == nil, "failed to connect unix socket")
	}
}

func startHTTPServer(addrCh chan string) {
	var foo Foo
	listen, _ := net.Listen("tcp", ":9999")
	_ = server.Register(&foo)
	server.HandleHTTP()
	addrCh <- listen.Addr().String()
	_ = http.Serve(listen, nil)
}

func httpCall(addrCh chan string) {
	c, _ := client.DialHTTP("tcp", <-addrCh)
	defer func() { _ = c.Close() }()

	time.Sleep(time.Second)
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			args := &Args{Num1: i, Num2: i * i}
			var reply int
			if err := c.CallTimeout(context.Background(), "Foo.Sum", args, &reply); err != nil {
				log.Fatal("call Foo.Sum error:", err)
			}
			log.Printf("%d + %d = %d", args.Num1, args.Num2, reply)
		}(i)
	}
	wg.Wait()
}

func TestHTTPCall(t *testing.T) {
	log.SetFlags(0)
	ch := make(chan string)
	go httpCall(ch)
	startHTTPServer(ch)
}

func (f Foo) Sleep(args Args, reply *int) error {
	time.Sleep(time.Second * time.Duration(args.Num1))
	*reply = args.Num1 + args.Num2
	return nil
}

func startXServer(addrCh chan string) {
	var foo Foo
	listen, _ := net.Listen("tcp", ":0")
	s := server.NewServer(0)
	_ = s.Register(&foo)
	addrCh <- listen.Addr().String()
	s.Accept(listen)
}

func foo(xc *xclient.XClient, ctx context.Context, typ, serviceMethod string, args *Args) {
	var reply int
	var err error
	switch typ {
	case "call":
		err = xc.Call("", ctx, serviceMethod, args, &reply)
	case "broadcast":
		err = xc.Broadcast(ctx, serviceMethod, args, &reply)
	}
	if err != nil {
		log.Printf("%s %s error: %v", typ, serviceMethod, err)
	} else {
		log.Printf("%s %s success: %d + %d = %d", typ, serviceMethod, args.Num1, args.Num2, reply)
	}
}
