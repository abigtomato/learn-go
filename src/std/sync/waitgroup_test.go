package sync

import (
	"fmt"
	"sync"
	"testing"
)

type WorkerChan struct {
	in   chan int
	done func()
}

// 创建worker
func createWorker(id int, wg *sync.WaitGroup) WorkerChan {
	w := WorkerChan{
		in: make(chan int),
		done: func() {
			// 解除waitGroup中的一个go程锁定
			wg.Done()
		},
	}

	// 开启消费协程，从worker.in消费数据
	go func(id int, w WorkerChan) {
		for n := range w.in {
			fmt.Printf("Worker %d Received %c\n", id, n)
			// 消费完毕解除锁定
			w.done()
		}
	}(id, w)

	return w
}

func TestWaitGroup(t *testing.T) {
	var wg sync.WaitGroup

	// 启动10个工作go程
	var workers [10]WorkerChan
	for i := 0; i < 10; i++ {
		workers[i] = createWorker(i, &wg)
		// 往waitGroup中注册协程任务
		wg.Add(1)
	}

	// 生产数据到worker.in
	for i, worker := range workers {
		worker.in <- 'a' + i
	}

	// 使用waitGroup使主go程等待
	wg.Wait()
}
