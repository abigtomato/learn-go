package goroutine

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
	"testing"
	"time"
)

// Func是用于记忆的函数类型
type Func func(key string) (interface{}, error)

// 调用Func的返回结果
type Result struct {
	value interface{}
	err   error
}

// 对Func返回结果的包装
type Entry struct {
	res   Result        // 结果
	ready chan struct{} // 当res准备好后关闭ready通道
}

// request是一条请求消息，key需要用Func来调用
type Request struct {
	key      string
	response chan<- Result // 客户端需要单个result
}

// 函数记忆
type Memo struct {
	requests chan Request
}

// New返回f的函数记忆，客户端之后需要调用Close
func New(f Func) *Memo {
	memo := &Memo{
		requests: make(chan Request),
	}

	// 开启服务
	go memo.server(f)
	return memo
}

// 获取记忆
func (memo *Memo) Get(key string) (interface{}, error) {
	response := make(chan Result)
	defer close(response)

	// requests由server消费，用于执行慢函数f(key)
	memo.requests <- Request{
		key:      key,
		response: response,
	}

	// 阻塞等待server将慢函数f的执行结果存入
	res := <-response
	return res.value, res.err
}

func (memo *Memo) Close() {
	close(memo.requests)
}

func (memo *Memo) server(f Func) {
	// cache只由一个go程所监控，不存在竞态问题
	cache := make(map[string]*Entry)

	// 消费requests通道，获取的request包含函数f的参数key和用于保存结果的通道response
	for req := range memo.requests {
		e := cache[req.key]
		if e == nil {
			// entry不存在的情况，实例新entry并调用函数f并关闭ready通道
			e = &Entry{
				ready: make(chan struct{}),
			}
			cache[req.key] = e

			// 调用f(key)，调用成功后通知数据准备完毕
			go e.call(f, req.key)
		}

		// 等待entry准备完毕发送结果给客户端
		go e.deliver(req.response)
	}
}

func (e *Entry) call(f Func, key string) {
	// 执行函数
	e.res.value, e.res.err = f(key)
	// 通知数据已准备完毕
	close(e.ready)
}

func (e *Entry) deliver(response chan<- Result) {
	// 等待该数据准备完毕
	<-e.ready
	// 向客户端发送结果
	response <- e.res
}

// 慢函数f
func httpGetBody(url string) (interface{}, error) {
	resp, err := http.Get(url)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Print(err)
		}
	}(resp.Body)
	if err != nil {
		return nil, err
	}

	return io.ReadAll(resp.Body)
}

// 生成url
func incomingURLs() []string {
	return []string{"https://www.baidu.com"}
}

// 并发非阻塞缓存示例（channel版）
func TestCacheChannel(t *testing.T) {
	m := New(httpGetBody)
	defer m.Close()

	var wg sync.WaitGroup
	for url := range incomingURLs() {
		wg.Add(1)

		go func(url string) {
			start := time.Now()

			// 获取函数缓存
			value, err := m.Get(url)
			if err != nil {
				log.Print(err)
			}
			fmt.Printf("%s, %s, %d bytes\n", url, time.Since(start), len(value.([]byte)))

			wg.Done()
		}(strconv.Itoa(url))
	}

	wg.Wait()
}
