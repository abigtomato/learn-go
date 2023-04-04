package std

import (
	"fmt"
	"log"
	"os"
	"testing"
)

var logger *log.Logger

func initLog() {
	// 日志配置
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)

	// 添加前缀
	log.SetPrefix("MyLog: ")

	// 日志输出位置
	f, _ := os.OpenFile("a.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0664)
	log.SetOutput(f)

	// 自定义日志配置
	logger = log.New(f, "MyLog: ", log.Ldate|log.Ltime|log.Llongfile)
}

func TestLog(t *testing.T) {
	initLog()

	logger.Println("test")

	defer fmt.Println("defer...")
	log.Print("my log")
	log.Fatal("fatal") // os.exit(1)退出系统 函数不返回 不会执行defer
}
