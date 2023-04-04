package context

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestWithDeadline(t *testing.T) {
	// context.WithDeadline 用于控制子 Goroutine 的执行截止时间
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(1*time.Second))

	go func(ctx context.Context, name string) {
		for {
			select {
			case <-ctx.Done():
				// ctx.Err 用于获取退出原因
				fmt.Println("stop", name, ctx.Err())
				return
			default:
				fmt.Println(name, "send request")
				time.Sleep(1 * time.Second)
			}
		}
	}(ctx, "work")

	time.Sleep(3 * time.Second)
	fmt.Println("before cancel")
	// 由于截止时间被设置为1s后，cancel 调用之前子 Goroutine 已经退出
	cancel()
	time.Sleep(3 * time.Second)
}
