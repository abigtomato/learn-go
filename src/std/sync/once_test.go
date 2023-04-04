package sync

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"
)

type Config struct {
	Server string
	Port   int64
}

var (
	// 1. sync.Once 是 Go 标准库提供的使函数只执行一次的实现，常应用于单例模式，例如初始化配置、保持数据库连接等
	// 2. 可以在代码的任意位置初始化和调用，因此可以延迟到使用时再执行，并发场景下是线程安全的
	// 3. 在多数情况下，sync.Once 被用于控制变量的初始化，这个变量的读写满足如下三个条件：
	// 	3.1. 当且仅当第一次访问某个变量时，进行初始化（写）
	//  3.2. 变量初始化过程中，所有读都被阻塞，直到初始化完成
	//  3.3. 变量仅初始化一次，初始化完成后驻留在内存里
	once   sync.Once
	config *Config
)

func ReadConfig() {
	// 函数 ReadConfig 需要读取环境变量，并转换为对应的配置
	// 环境变量在程序执行前已经确定，执行过程中不会发生改变
	// ReadConfig 可能会被多个协程并发调用，为了提升性能（减少执行时间和内存占用），使用 sync.Once 是一个比较好的方式
	once.Do(func() {
		var err error
		config = &Config{Server: os.Getenv("TT_SERVER_URL")}
		config.Port, err = strconv.ParseInt(os.Getenv("TT_PORT"), 10, 0)
		if err != nil {
			config.Port = 8080
		}
		fmt.Println("init config")
	})
}

func TestOnce(t *testing.T) {
	for i := 0; i < 10; i++ {
		go func() { ReadConfig() }()
	}
	time.Sleep(3 * time.Second)
}
