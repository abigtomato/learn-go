package basic

import (
	"errors"
	"fmt"
	"log"
	"testing"
)

// 使用defer + recover()来捕获处理异常
func TestTryRecover(t *testing.T) {
	// defer将匿名函数调用语句压入defer栈中，等待函数执行结束调用
	defer func() {
		// recover()内置函数可以捕获到异常
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				log.Fatalln("load config err: ", err)
			}
			log.Fatalln("load config panic: ", r)
		}

		//r := recover()
		//if err, ok := r.(error); ok {
		//	fmt.Printf("Error occurred: %s\n", err.Error())
		//}
	}()

	num1 := 10
	num2 := 0

	//if num2 == 0 {
	//	panic("panic err")
	//	//log.Fatalln("fatalln err")
	//}

	fmt.Printf("result: %v\n\n", num1/num2)
}

// 抛出自定义错误测试
func readConf() (err error) {
	// 返回一个自定义错误
	return errors.New("配置文件读取错误")
}

// 测试自定义错误
func TestError(t *testing.T) {
	if err := readConf(); err != nil {
		// 若出错，则打印自定义错误的信息，终止程序
		fmt.Println(err.Error())
	}
}
