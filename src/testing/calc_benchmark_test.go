package testing

import (
	"bytes"
	"fmt"
	"html/template"
	"testing"
	"time"
)

// 1. 基准测试函数名必须以 Benchmark 开头，后面一般跟着待测试的函数名
// 2. 参数为 b *testing.B
// 3. 执行基准测试时，需要添加 -bench 参数，如 $ go test -benchmem -bench .
// 4. 基准测试报告对应的列值：迭代次数、基准测试花费时间、一次迭代处理的字节数、总的分配内存的次数、总的分配内存的字节数
func BenchmarkHello(b *testing.B) {
	// 模拟耗时操作
	time.Sleep(3 * time.Second)
	// 如果在运行前基准测试需要一些耗时的配置，则可以使用 b.ResetTimer() 先重置定时器
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fmt.Println("hello")
	}
}

// 使用 RunParallel 测试并发性能
func BenchmarkParallel(b *testing.B) {
	temp := template.Must(template.New("test").Parse("Hello, {{.}}!"))
	b.RunParallel(func(pb *testing.PB) {
		var buf bytes.Buffer
		for pb.Next() {
			// 所有 goroutine 一起，循环一共执行 b.N 次
			buf.Reset()
			_ = temp.Execute(&buf, "world")
		}
	})
}
