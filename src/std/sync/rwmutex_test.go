package sync

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

// 读写锁
func TestRWMutex(t *testing.T) {
	var value int
	// 读写锁
	var rwLock sync.RWMutex

	for i := 0; i < 5; i++ {
		go func(index int) {
			for {
				// RLock()：当前go程以读模式加锁
				// 读锁共享：当其他读go程执行此处，会共享读锁，执行下面的读逻辑
				rwLock.RLock()

				num := value
				fmt.Printf("读go程%d号 <- %d\n", index, num)

				// 解锁读锁
				rwLock.RUnlock()
			}
		}(i)
	}

	for i := 0; i < 5; i++ {
		go func(index int) {
			for {
				// Lock()：当前go程以写模式加锁
				// 写锁独占：执行到此处的其他任何go程在锁未释放前都会阻塞
				// 写锁优先级高：当多个读写go程要对同一块资源操作，写锁会优先被抢到
				rwLock.Lock()

				value = rand.Intn(1000)
				fmt.Printf("写go程%d号 -> %d\n", index, value)
				time.Sleep(time.Millisecond * 300)

				rwLock.Unlock()
			}
		}(i)
	}
}
