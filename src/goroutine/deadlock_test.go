package goroutine

import (
	"fmt"
	"testing"
)

// 死锁1: 单go程无缓冲管道死锁
func TestDeadLock1(t *testing.T) {
	ch := make(chan int)
	// 此处需要有读端存在，否则死锁
	ch <- 789
	num := <-ch
	fmt.Printf("num: %v\n", num)
}

// 死锁2: go程建管道的访问顺序倒置死锁
func TestDeadLock2(t *testing.T) {
	ch := make(chan int)
	// 此处的读端会阻塞，没有机会执行下面的代码创建新go程做为写端
	num := <-ch
	fmt.Printf("num: %v\n", num)

	go func() {
		ch <- 789
	}()
}

// 死锁3: 多go程多管道交叉死锁
func TestDeadLock3(t *testing.T) {
	ch1 := make(chan int)
	ch2 := make(chan int)

	go func() {
		for {
			select {
			// 从ch1获取数据到ch2
			case num := <-ch1:
				ch2 <- num
			}
		}
	}()

	for {
		select {
		// 从ch2获取数据到ch1，两端同时执行都没有数据都会阻塞
		case num := <-ch2:
			ch1 <- num
		}
	}
}
