package geerpc

import (
	"Golearn/src/projects/geerpc/registry"
	"Golearn/src/projects/geerpc/server"
	"Golearn/src/projects/geerpc/xclient"
	"context"
	"log"
	"net"
	"net/http"
	"sync"
	"testing"
	"time"
)

func startRegistry(wg *sync.WaitGroup) {
	listen, _ := net.Listen("tcp", ":9999")
	registry.HandleHTTP()
	wg.Done()
	_ = http.Serve(listen, nil)
}

func startRegistryServer(registryAddr string, wg *sync.WaitGroup) {
	var foo Foo
	listen, _ := net.Listen("tcp", ":0")
	s := server.NewServer(0)
	_ = s.Register(&foo)
	registry.Heartbeat(registryAddr, "tcp@"+listen.Addr().String(), 0)
	wg.Done()
	s.Accept(listen)
}

func call(registry string) {
	d := xclient.NewGeeRegistryDiscovery(registry, 0)
	xc := xclient.NewXClient(d, xclient.RandomSelect, nil)
	defer func(xc *xclient.XClient) {
		_ = xc.Close()
	}(xc)
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			foo(xc, context.Background(), "call", "Foo.Sum", &Args{Num1: i, Num2: i * i})
		}(i)
	}
	wg.Wait()
}

func broadcast(registry string) {
	d := xclient.NewGeeRegistryDiscovery(registry, 0)
	xc := xclient.NewXClient(d, xclient.RandomSelect, nil)
	defer func(xc *xclient.XClient) {
		_ = xc.Close()
	}(xc)
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			foo(xc, context.Background(), "broadcast", "Foo.Sum", &Args{Num1: i, Num2: i * i})
			ctx, _ := context.WithTimeout(context.Background(), time.Second*2)
			foo(xc, ctx, "broadcast", "Foo.Sleep", &Args{Num1: i, Num2: i * i})
		}(i)
	}
	wg.Wait()
}

func TestRegistry(t *testing.T) {
	log.SetFlags(0)
	registryAddr := "http://localhost:9999/_geerpc_/registry"
	var wg sync.WaitGroup
	wg.Add(1)
	go startRegistry(&wg)
	wg.Wait()

	time.Sleep(time.Second)
	wg.Add(2)
	go startRegistryServer(registryAddr, &wg)
	go startRegistryServer(registryAddr, &wg)
	wg.Wait()

	time.Sleep(time.Second)
	call(registryAddr)
	broadcast(registryAddr)
}
