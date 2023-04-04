package sync

import (
	"bytes"
	"encoding/json"
	"sync"
	"testing"
)

type Student struct {
	Name   string
	Age    int32
	Remark [1024]byte
}

var buf, _ = json.Marshal(Student{Name: "abigtomato", Age: 25})

// 1. sync.Pool 是可伸缩的，同时也是并发安全的，其大小仅受限于内存的大小
// 2. 用于存储那些被分配了但是没有被使用，而未来可能会使用的值。这样就可以不用再次经过内存分配，可直接复用已有对象，减轻 GC 的压力，从而提升系统的性能
// 3. 其大小是可伸缩的，高负载时会动态扩容，存放在池中的对象如果不活跃了会被自动清理
// 4. 声明对象池只需要实现 New 函数即可。对象池中没有对象时，将会调用 New 函数创建
var stuPool = sync.Pool{
	New: func() any {
		return new(Student)
	},
}

// 使用Benchmark测试pool的性能
func BenchmarkUnmarshal(b *testing.B) {
	b.Run("BenchmarkUnmarshal", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			stu := &Student{}
			_ = json.Unmarshal(buf, stu)
		}
	})

	b.Run("BenchmarkUnmarshalWithPool", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			// 用于从对象池中获取对象，因为返回值是 interface{}，因此需要类型转换
			stu := stuPool.Get().(*Student)
			_ = json.Unmarshal(buf, stu)
			// 在对象使用完毕后，返回对象池
			stuPool.Put(stu)
		}
	})
}

var bufferPool = sync.Pool{
	New: func() any {
		return &bytes.Buffer{}
	},
}

var data = make([]byte, 10000)

func BenchmarkBuffer(b *testing.B) {
	b.Run("BenchmarkBuffer", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var buf bytes.Buffer
			buf.Write(data)
		}
	})

	b.Run("BenchmarkBufferWithPool", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			buf := bufferPool.Get().(*bytes.Buffer)
			buf.Write(data)
			buf.Reset()
			bufferPool.Put(buf)
		}
	})
}
