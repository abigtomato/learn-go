package context

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestWithValue(t *testing.T) {
	type Options struct{ Interval time.Duration }

	ctx, cancel := context.WithCancel(context.Background())
	// context.WithValue 创建了一个基于 ctx 的子 Context，并携带值 options，用于向子 Goroutine 传递数据
	context.WithValue(ctx, "options", &Options{1})

	go func(ctx context.Context, name string) {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("stop", name)
				return
			default:
				fmt.Println(name, "send request")
				// 子 Goroutine 获取传递的值
				op := ctx.Value("options").(*Options)
				time.Sleep(op.Interval * time.Second)
			}
		}
	}(ctx, "work")

	time.Sleep(3 * time.Second)
	cancel()
	time.Sleep(3 * time.Second)
}
