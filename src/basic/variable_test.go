package basic

import (
	"fmt"
	"math"
	"testing"
	"unsafe"
)

// 声明变量（变量存在初值）
func TestVariableZeroValue(t *testing.T) {
	// 64位有符号整数类型，占用8字节存储空间
	// 表数范围-2^63 ~ 2^63-1，63次方是因为第一位字节做为符号位，-1是因为+0和-0只需要表示一种情况即可
	// 32位则占4字节，表数-2^31 ~ 2^31 - 1
	var a int64

	// 8位无符号整数类型，占用1字节存储空间
	// 表数范围0 ~ 255，8位二进制位全为0则为0，全为1则为255
	var b uint8

	fmt.Printf("go中所有变量均存在初值: int64 = %d, uint8 = %d\n", a, b)
}

// 显式创建变量（定义时注明变量类型）
func TestVariableInitialValue(t *testing.T) {
	// 64位浮点型，浮点数 = 符号位 + 指数位 + 尾数位
	// 尾数部分可能会丢失，造成精度损失
	var a, b = 128.0375001, -127.0468001

	// 仅小数部分表示，忽略整数部分的0
	var c = .512 // 表示0.512

	// 科学计数法表示
	var d = 5.12e2  // 5.12 * 10^2
	var e = 5.12e-2 // 5.12 / 10^2

	fmt.Printf("%f, %f, %f, %f, %f\n", a, b, c, d, e)

	// 存储: 字符 --> 对应码值 --> 二进制 --> 存储
	// 读取: 二进制 --> 码值 --> 字符 --> 读取
	var f byte = 'a' // 本质是存储整数，直接打印的是该字符对应的utf8码
	var g byte = '0' // 字符'0'
	var h int = '北'  // 汉字在计算机中占3个字节（utf8编码存储），超过255的范可以使用int存储（存储其utf-8编码）

	// %c格式化输出时会输出该数字对应的unicode对应字符
	fmt.Printf("字符'a' = %c, 字符'a'对应的码 = %d; 字符'0' = %c\n", f, f, g)
	fmt.Printf("字符'0'对应的码 = %d; 字符'北' = %c, 字符'北'对应的码 = %d", g, h, h)
}

// 隐式创建变量（go编译器自行推断类型）
func TestVariableTypeDeduction(t *testing.T) {
	// 等号左边的代表变量名指向的内存空间（左值写操作）
	// 等号右边的代表内存空间中存储的数据（右值读操作）
	var a, b, c = 233, 127.00275, "SparkStreaming"

	fmt.Printf("数据类型 = %T, 占用的字节数 = %dbyte\n", a, unsafe.Sizeof(a))
	fmt.Println(a, b, c)
}

// 简化变量创建（定义 + 赋值）
func TestVariableShorter(t *testing.T) {
	// := 表示首次定义变量并赋值
	a, b, c, s := 3, 4, false, "etcd"
	// = 表示为变量重新赋值
	b = 5

	fmt.Printf("批量创建: a = %v, b = %v, c = %v, s = %v\n", a, b, c, s)
}

// 强制类型转换（go规定所有类型都必须显式转换）
func TestTriangle(t *testing.T) {
	var a, b = 3, 4
	var c int
	// 通过内置函数 float64() 或库函数 math.Sqrt() 显示转换
	c = int(math.Sqrt(float64(a*a + b*b)))

	fmt.Printf("类型强转: a = %v, b = %v, c = %v\n", a, b, c)
}

// 常量定义（常量可以当做任意类型使用）
func TestCosts(t *testing.T) {
	const filename = "ElasticSearch"
	const a, b = 3, 4

	var c int
	c = int(math.Sqrt(a*a + b*b))

	fmt.Println(filename, c)
}

// 枚举定义（一组常量）
func TestEnums(t *testing.T) {
	// iota自增，组内常量依次递增
	const (
		cpp = iota
		java
		python
		golang
	)

	const (
		b = 1 << (10 * iota)
		kb
		mb
		gb
		tb
		pb
		zb
	)

	fmt.Printf("iota常量递增: %v, %v, %v, %v\n", cpp, java, python, golang)
	fmt.Printf("1Bit = %v, 1KB = %v, 1MB = %v, 1GB = %v, 1TB = %v, 1PB = %v, 1ZB = %v\n", b, kb, mb, gb, tb, pb, zb)
}
