package sync

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

// 条件变量
var cond sync.Cond

// 生产者
func cProducer(out chan<- int, idx int) {
	for {
		func() {
			// 使用条件变量添加互斥锁
			cond.L.Lock()
			defer cond.L.Unlock()

			// 使用for循环判断条件变量是否满足（而不是使用if）
			// 1. 设想使用if的场景：第一个进来的go程通过了if判断，还未执行wait就失去了cpu的使用权
			// 2. 其他go程抢占到cpu需要在此处循环判断是否满足条件，若是if其他go程会直接向下执行（因为第一个进来的go程已经通过了判断）
			for len(out) == cap(out) {
				// cond.Wait()
				// 1. 使当前go程在条件变量cond上阻塞并等待该条件变量的唤醒
				// 2. 释放当前go程持有的互斥锁（相当于进行cond.L.Unlock()），和第1步同为1个原子操作
				// 3. 当阻塞在此的go程被唤醒，Wait()函数返回时，该go程重新获取互斥锁（相当于进行cond.L.Lock()操作）
				cond.Wait()
			}

			// 具体的生产操作
			num := rand.Intn(1000)
			out <- num
			fmt.Printf("生产者%d号 -> %d\n", idx, num)
		}()

		// cond.Signal() 给一个在条件变量cond上阻塞的go程发送唤醒通知
		cond.Signal()
		time.Sleep(time.Second)
	}
}

// 消费者
func cConsumer(in <-chan int, idx int) {
	for {
		func() {
			cond.L.Lock()
			defer cond.L.Unlock()

			// 管道没有数据可以消费了，就使当前消费go程wait
			for len(in) == 0 {
				cond.Wait()
			}

			// 具体的消费操作
			num := <-in
			fmt.Printf("消费者%d号 <- %d\n", idx, num)
		}()

		cond.Signal()
		time.Sleep(time.Second)
	}
}

//  1. 互斥锁 sync.Mutex 通常用来保护临界区和共享资源，条件变量 sync.Cond 用来协调想要访问共享资源的 goroutine
//  2. sync.Cond 经常用在多个 goroutine 等待，一个 goroutine 通知（事件发生）的场景
//  3. 场景举例：
//     3.1. 有一个协程在异步地接收数据，剩下的多个协程必须等待这个协程接收完数据，才能读取到正确的数据
//     3.2. 第一种方案是定义一个全局的变量来标志第一个协程数据是否接受完毕，剩下的协程，反复检查该变量的值，直到满足要求
//     3.3. 第二种方案是创建多个 channel，每个协程阻塞在一个 channel 上，由接收数据的协程在数据接收完毕后，逐个通知
//     3.4. 第三种方案是使用一个带缓冲的 channel，缓冲区大小和go程数量一致，数据处理完毕后发送指定数量的消息通知所有go程
//     3.5. 以上的方法都会带来而外的复杂度和局限性
func TestCondition(t *testing.T) {
	// 随机数种子
	rand.Seed(time.Now().UnixNano())

	// 使用互斥锁初始化条件变量的锁字段
	cond.L = new(sync.Mutex)

	// 数据管道
	numChan := make(chan int, 3)
	// 退出标记管道
	quitChan := make(chan bool)

	for i := 0; i < 5; i++ {
		go cProducer(numChan, i)
	}

	for i := 0; i < 5; i++ {
		go cConsumer(numChan, i)
	}

	for {
		if _, ok := <-quitChan; !ok {
			break
		}
	}
}
