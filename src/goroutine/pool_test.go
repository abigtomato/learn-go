package goroutine

import (
	"errors"
	"fmt"
	"log"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// 表示pool关闭状态的常量
const CLOSED = 1

var (
	chanSize = func() int {
		// 如果GOMAXPROCS为1时，使用阻塞channel
		if runtime.GOMAXPROCS(0) == 1 {
			return 0
		}
		// 如果GOMAXPROCS大于1时，使用非阻塞channel
		return 1
	}()
	ClosedError          = errors.New("this pool has been closed")
	InvalidPoolSizeError = errors.New("invalid size for pool")
)

// Go程池（调度器）
type Pool struct {
	cap          int32             // 池容量
	closed       int32             // 关闭标记
	jobQueue     chan Job          // 总任务队列
	workerQueue  chan *Worker      // worker队列
	once         sync.Once         // 保证某些操作只执行一次
	PanicHandler func(interface{}) // 用户自定义错误处理
}

// 创建新go池
func NewPool(size int) (*Pool, error) {
	if size <= 0 {
		return nil, InvalidPoolSizeError
	}

	pool := &Pool{
		cap:      int32(size),
		jobQueue: make(chan Job, chanSize),
	}

	if chanSize != 0 {
		pool.workerQueue = make(chan *Worker, size)
	} else {
		pool.workerQueue = make(chan *Worker)
	}

	return pool, nil
}

// 启动pool
func (p *Pool) Run() {
	// 根据最大pool容量创建worker
	for i := 0; i < int(p.cap); i++ {
		// 创建worker实例
		worker := NewWorker(p)
		// 开启go程分支执行worker逻辑
		go worker.start()
		// worker入队
		p.workerQueue <- worker
	}

	go p.scheduler()
}

// 开启调度
func (p *Pool) scheduler() {
	for {
		select {
		// 监听任务队列的数据
		case job := <-p.jobQueue:
			// 若有任务需要处理，出队一个worker进行处理
			worker := <-p.workerQueue
			// 存入worker的专属任务队列中
			worker.task <- job
		}
	}
}

// 关闭pool
func (p *Pool) Close() {
	// 1. Once包含一个bool变量和一个互斥量，bool变量记录逻辑初始化是否完成，互斥量负责保护bool变量和客户端的数据结构
	// 2. Once的唯一方法Do以需要执行的初始化函数作为参数
	// 3. 每次调用Do时会先锁定互斥量并检查里边的bool变量，第一次调用时bool变量为false
	// 4. Do会调用初始化函数并将变量置为true，后续的再次调用相当于空操作
	p.once.Do(func() {
		defer func() {
			close(p.jobQueue)
			close(p.workerQueue)
		}()

		atomic.StoreInt32(&p.closed, 1)

		p.jobQueue = nil
		p.workerQueue = nil
	})
}

// 提交任务给pool
func (p *Pool) Submit(job Job) error {
	// 判断pool是否已经关闭
	if atomic.LoadInt32(&p.closed) == CLOSED {
		return ClosedError
	}

	// 任务入队
	p.jobQueue <- job

	return nil
}

type Job func()

// 工作节点
type Worker struct {
	pool *Pool     // 所属的池
	task chan Job  // 每个worker专属的任务队列
	quit chan bool // 退出标记管道
}

// 实例新工作节点
func NewWorker(pool *Pool) *Worker {
	return &Worker{
		pool: pool,
		task: make(chan Job, chanSize),
		quit: make(chan bool),
	}
}

// 工作节点开始干活
func (w *Worker) start() {
	for {
		select {
		case job := <-w.task:
			// 执行任务的具体逻辑
			job()

			// worker执行完毕后重新入队pool
			w.pool.workerQueue <- w

			// go程执行中的错误处理
			if p := recover(); p != nil {
				// 若用户自定义了错误处理函数则执行
				if w.pool.PanicHandler != nil {
					w.pool.PanicHandler(p)
				} else {
					// 否则默认错误处理
					log.Printf("worker exits from a panic: %v", p)
				}
			}
		case <-w.quit:
			return
		}
	}
}

// 停止worker
func (w *Worker) stop() {
	go func() {
		// 存入结束标记
		w.quit <- true
	}()
}

func TestPool(t *testing.T) {
	var wg sync.WaitGroup

	goPool, _ := NewPool(100)
	goPool.Run()
	defer goPool.Close()

	for i := 0; i < 10; i++ {
		// go程加锁
		wg.Add(1)
		_ = goPool.Submit(func() {
			time.Sleep(10 * time.Millisecond)
			// 这里写任务逻辑
			fmt.Println("Hello Goroutine Pool!")
			// go程解锁
			wg.Done()
		})
	}

	// 主go程等待
	wg.Wait()
}
