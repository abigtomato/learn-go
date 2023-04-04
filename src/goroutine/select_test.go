package goroutine

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

// 创建生产者
func newProducer(pid, total int) chan int {
	pChan := make(chan int)

	go func() {
		defer close(pChan)

		x, y := 1, 1
		for i := 0; i < total; i++ {
			time.Sleep(time.Duration(rand.Intn(1500)) * time.Millisecond)
			fmt.Printf("生产者%d号 -> %d\n", pid, x)

			pChan <- x
			x, y = y, x+y
		}
	}()

	return pChan
}

// 创建工作go程
func newWorker() (chan int, chan bool) {
	intChan := make(chan int, 3)
	quitChan := make(chan bool)

	go func() {
		defer func() {
			close(intChan)
			close(quitChan)
		}()

		for num := range intChan {
			fmt.Printf("消费数据: %d\n", num)
		}
		// 无数据消费，退出
		quitChan <- true
	}()

	return intChan, quitChan
}

// 使用 select + channel 生产消费斐波那契数列
func TestSelect(t *testing.T) {
	// 原始数据管道
	pChan1, pChan2 := newProducer(1, 20), newProducer(2, 30)

	// 工作和退出管道
	workerChan, quitChan := newWorker()

	// 每隔指定时间提供数据的只读管道
	tick := time.Tick(time.Second)
	// 单select超时的只读管道，超过指定的时间会提供数据
	timeout := time.After(1500 * time.Millisecond)
	// 最大超时只读管道
	maxTimeout := time.After(100 * time.Second)

	for {
		// select语句用于监听管道的数据流动
		// 1. 按照顺序从头到尾评估每一个 case 后面的 I/O 操作
		// 2. 当任意一个 case 可执行（即管道解阻塞），则会执行 case 内的代码
		// 3. 若本次有多个 case 可执行，那么从可执行的 case 中任意选择一条执行
		// 4. 若本次没有 case 可执行（即所有case的管道都阻塞），本次执行 default，若无 default 则 select 阻塞，直到至少有一个通信可以进行下去
		select {
		case num := <-pChan1:
			workerChan <- num
		case num := <-pChan2:
			workerChan <- num
		case <-tick:
			fmt.Printf("定时状态汇报")
		case <-timeout:
			fmt.Println("select监听超时")
		case <-maxTimeout:
			fmt.Println("超时退出")
			return
		case <-quitChan:
			return
		}
	}
}
