package std

import (
	"fmt"
	"testing"
)

// 追加内置函数将元素追加到切片的末尾:
// 1. 如果它有足够的容量，则会重新切片目标以容纳新元素。
// 2. 如果没有，将分配一个新的基础阵列，追加返回更新的切片。因此，有必要将追加的结果存储在保存切片本身的变量中。
func TestAppend(t *testing.T) {
	slice := []int{1, 2, 3, 4}
	slice = append(slice, 5)
	fmt.Println(slice)

	slice1 := []int{6, 7, 8, 9}
	slice = append(slice, slice1...)
	fmt.Println(slice)

	// 作为特殊情况，将字符串附加到字节片是合法的，如下所示
	byteSlice := []byte("hello ")
	byteSlice = append(byteSlice, "world"...)
}

// 根据传入的类型返回v的长度:
// 1. Array: 返回v中元素的数量
// 2. Pointer to array: 返回v中元素的数量（即使v为nil）
// 3. Slice, or map: 返回v中的元素数量，如果v为nil，那么len(v)为0
// 4. String: 返回v中的字节数
// 5. Channel: 返回通道缓冲区中排队（未读）的元素数，如果v为零，则len(v)为零
func TestLen(t *testing.T) {
	str := "Hello World"
	fmt.Println(len(str))

	slice := []int{1, 2, 3}
	fmt.Println(len(slice))
}

// new 和 make 区别:
// 1. make 只能用来分配及初始化类型为 slice、map、chan 的数据；new 可以分配任意类型的数据
// 2. new 分配返回的是指针，即类型 *T；make 返回引用，即 T
// 3. new 分配的空间被清零；make 分配后，会进行初始化
func TestNewAndMake(t *testing.T) {
	b := new(bool)
	fmt.Printf("b: %T\n", b)  // *bool
	fmt.Printf("b: %v\n", b)  // 0xc00000a098
	fmt.Printf("b: %v\n", *b) // false

	i := new(int)
	fmt.Printf("i: %T\n", i)  // *int
	fmt.Printf("i: %v\n", i)  // 0xc00000a0b8
	fmt.Printf("i: %v\n", *i) // 0

	s := new(string)
	fmt.Printf("s: %T\n", s)  // *string
	fmt.Printf("s: %v\n", s)  // 0xc000050260
	fmt.Printf("s: %v\n", *s) // ""

	var p = new([]int)
	fmt.Printf("p: %v\n", p) // &[]

	v := make([]int, 10)
	fmt.Printf("v: %v\n", v) // [0 0 0 0 0 0 0 0 0 0]
}
