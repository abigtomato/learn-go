package context

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestWithCancel(t *testing.T) {
	// context.Background 用于创建根 Context，通常在 main 函数、初始化和测试代码中创建，作为顶层 Context
	// context.WithCancel 创建可取消的子 Context，同时返回函数 cancel
	ctx, cancel := context.WithCancel(context.Background())

	go func(ctx context.Context, name string) {
		for {
			select {
			// 在子 Goroutine 中，使用 select 监听 Context 判断是否需要退出
			case <-ctx.Done():
				fmt.Println("stop", name)
				return
			default:
				fmt.Println(name, "send request")
				time.Sleep(1 * time.Second)
			}
		}
	}(ctx, "work")

	time.Sleep(3 * time.Second)
	// 主 Goroutine 中，调用 cancel() 通知子 Goroutine 退出
	cancel()
	time.Sleep(3 * time.Second)
}
