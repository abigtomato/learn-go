package basic

import (
	"fmt"
	"os"
	"testing"
)

func TestBranchIf(t *testing.T) {
	const fineName = "./data/demo.txt"

	// if的判断条件中可以先定义变量并赋值；
	// os库提供io流常用的函数；
	// ReadFile提供读取文件内容的功能，返回值为[]byte类型的文件内容和错误信息。

	if contents, err := os.ReadFile(fineName); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(contents))
	}
}

func TestBranchSwitch(t *testing.T) {
	var g string
	var score int

	switch {
	// case中的break由go编译器自动添加
	case score < 0 || score > 100:
		// panic提供报错机制
		// fmt.Sprintf返回格式化后的字符串
		panic(fmt.Sprintf("Wrong score: %d", score))
	case score < 60, score == 60, score <= 60:
		g = "D"
	case score < 80:
		g = "C"
		// fallthrough关键字默认穿透一层case
		fallthrough
	case score < 90:
		g = "B"
	case score <= 100:
		g = "A"
	// 所有case条件都不符合，执行default默认
	default:
		fmt.Println("Default ......")
	}

	fmt.Println(g)
}
