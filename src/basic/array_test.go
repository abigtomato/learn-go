package basic

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestArrayCreate(t *testing.T) {
	// 显式声明（指定类型和下标）
	var arr1 = [3]int{0: 1, 1: 2, 2: 3}
	// 声明并预创建初值
	var arr2 = [3]int{1, 2, 3}
	// [...]声明不定长数组
	arr3 := [...]int{1, 2, 3, 4, 5}
	// 二维数组，值都默认为0
	var grid [4][5]int
	fmt.Println(arr1, arr2, arr3, grid)
}

func TestArrayIterate(t *testing.T) {
	arr := [5]int{1, 2, 3, 4, 5}
	tdArr := [5][5]int{{1, 2, 3}, {3, 2, 1}}

	// 遍历一维数组
	for i, v := range arr {
		fmt.Printf("(%v, %v) ", i, v)
	}
	fmt.Println()

	// 遍历二维数组
	for i := range tdArr {
		// 若无需第一个返回值，可用_代替
		for _, v := range tdArr[i] {
			fmt.Printf("%v ", v)
		}
		fmt.Println()
	}
}

func TestArrayDeliver(t *testing.T) {
	arr := [5]int{1, 2, 3, 4, 5}

	// 若是值传递，则会在函数中拷贝出一份新数组
	fmt.Println("原始数组: ", arr)
	func(arr [5]int) {
		// 数组值传递（[5]int 表示接收长度为5类型为int的数组）
		arr[0] = 100
	}(arr)
	fmt.Println("值传递后: ", arr)

	// 若是指针传递，则在函数中操作则会改变原数组
	func(arr *[5]int) {
		// 数组指针传递（*[5]int 表示接收长度为5类型为int的数组的内存地址）
		arr[0] = 100
	}(&arr)
	fmt.Println("指针传递后: ", arr)
}

// 示例1：利用编码的顺序，让字符参与运算存储26个大写字母到数组中
func TestArrayDemo1(t *testing.T) {
	var myChars [26]byte

	for i := 0; i < len(myChars); i++ {
		myChars[i] = 'A' + byte(i)
	}

	for i := 0; i < 26; i++ {
		fmt.Printf("%c", myChars[i])
	}

	fmt.Println()
}

// 示例2：随机生成数组，反转打印
func TestArrayDemo2(t *testing.T) {
	var intArr [5]int
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < len(intArr); i++ {
		intArr[i] = rand.Intn(100)
	}
	fmt.Println("交换前: ", intArr)

	for i := 0; i < len(intArr)/2; i++ {
		intArr[i], intArr[len(intArr)-1-i] = intArr[len(intArr)-1-i], intArr[i]
	}
	fmt.Println("交换后: ", intArr)
}

func TestArrayMemory(t *testing.T) {
	// 数组的内存布局分析:
	// 1. 数组是值类型，数组的地址就是数组首元素的地址；
	// 2. 数组是一段连续的内存空间，每个元素的间隔是由元素大小决定的。
	var arr [3]int
	fmt.Printf("数组的地址=%p\n", &arr)
	fmt.Printf("首元素的地址=%p\n", &arr[0])
	fmt.Printf("第二个元素的地址=%p\n", &arr[1])
	fmt.Printf("第三个元素的地址=%p\n", &arr[2])
	fmt.Println()

	// 二维数组内存布局分析:
	// 1. 存储多个指针，分别指向底层的一维数组地址；
	// 2. 数组在内存中的空间是连续的。
	var matrix = [2][3]int{{1, 2, 3}, {4, 5, 6}}
	fmt.Printf("二维数组地址=%p\n", &matrix)
	fmt.Printf("二维数组第一个指针指向=%p\n", &matrix[0])
	fmt.Printf("二维数组第二个指针指向=%p\n", &matrix[1])
	fmt.Printf("内部第一个一维数组地址=%p\n", &matrix[0][0])
	fmt.Printf("内部第二个一维数组地址=%p\n", &matrix[1][0])
	fmt.Println()
}
