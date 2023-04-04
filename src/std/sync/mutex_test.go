package sync

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// 可以进行原子操作的Int
type AtomicInt struct {
	value int
	// 互斥锁
	lock sync.Mutex
}

// 原子自增
func (i *AtomicInt) increment() {
	// 若想在函数中使用defer为一段代码加锁，可以使用匿名函数实现
	func() {
		i.lock.Lock()
		defer i.lock.Unlock()
		i.value++
	}()
}

// 原子获取
func (i *AtomicInt) get() int {
	// 为当前进入到此的go程加锁，其他所有go程执行到此不能获取锁进入阻塞
	// Lock() 建议锁：由操作系统提供，建议在编程时使用的锁
	i.lock.Lock()
	// 函数结束释放锁（自动唤醒阻塞在这把锁上的所有go程，让他们去争抢锁）
	defer i.lock.Unlock()
	return i.value
}

// 互斥锁
func TestMutex(t *testing.T) {
	var a AtomicInt

	// 并发读写可能会发生同步问题
	a.increment()
	go func() {
		a.increment()
	}()

	time.Sleep(time.Millisecond)
	fmt.Println(a.get())
}
