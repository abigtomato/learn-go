package context

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestWithTimeout(t *testing.T) {
	// context.WithTimeout 用于控制子 Goroutine 执行时间的 Context，具有超时通知机制
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

	go func(ctx context.Context, name string) {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("timeout stop", name)
				return
			default:
				fmt.Println("send request")
				time.Sleep(1 * time.Second)
			}
		}
	}(ctx, "work")

	time.Sleep(3 * time.Second)
	fmt.Println("before cancel")
	// 由于超时时间设置为2s，主 Goroutine 3s后才会调用 cancel，所以在这之前子 Goroutine 已经退出了
	cancel()
	time.Sleep(3 * time.Second)
}
