package basic

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"
)

// for循环（10进制转2进制例子）
func TestConvertToBin(t *testing.T) {
	n := 5

	var result string

	// 无起始式，有迭代式形式的for循环
	for ; n > 0; n /= 2 {
		result = strconv.Itoa(n%2) + result
	}

	fmt.Println(result)
}

// for循环（无起始和迭代式，按行读取文件）
func TestPrintFile(t *testing.T) {
	filename := "demo.txt"

	// os包提供操作系统相关的函数库
	// Open()提供打开文件的功能
	file, _ := os.Open(filename)

	// bufio包提供缓冲流的函数库
	// NewScanner()按照文件创建扫描器
	scanner := bufio.NewScanner(file)
	// Scan()每次扫描文件一行并向后移动行标记，扫描到末尾则返回nil
	for scanner.Scan() {
		// Text()将当前标记指向的行转换为文本返回
		fmt.Println(scanner.Text())
	}
}

// for循环模拟while
func TestFor2while(t *testing.T) {
	// 使用Unix时间戳做为随机数种子，使产生的随机数不会重复
	rand.Seed(time.Now().Unix())

	j := 1
	for j <= 10 {
		// 生成[0, 100)间的整数
		// +1 是为了生成[1, 100]的整数
		n := rand.Intn(100) + 1
		fmt.Printf("Docker/Kubernetes [%d]\n", n)
		j++
	}

	// 内部带条件的死循环
	k := 1
	for {
		if k > 10 {
			break
		}
		k++
	}
	fmt.Printf("k = %d\n", k)
}

// 打印空心金字塔案例
func TestPyramid(t *testing.T) {
	totalLevel := 9

	// 层数控制
	for i := 1; i <= totalLevel; i++ {
		// 打印空格 => 空格的规律: 总层数-当前层数
		for j := 1; j <= totalLevel-i; j++ {
			fmt.Print(" ")
		}
		// 打印* => *的规律: 2*当前层数-1
		for k := 1; k <= 2*i-1; k++ {
			// 控制中间空出
			if k == 1 || k == 2*i-1 || i == totalLevel {
				fmt.Print("*")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
}
