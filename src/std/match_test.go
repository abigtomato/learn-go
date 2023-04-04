package std

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
	"time"
)

func TestMatch(t *testing.T) {
	fmt.Printf("math.MaxFloat64: %v\n", math.MaxFloat64)
	fmt.Printf("math.SmallestNonzeroFloat64: %v\n", math.SmallestNonzeroFloat64)
	fmt.Printf("math.MaxFloat32: %v\n", math.MaxFloat32)
	fmt.Printf("math.SmallestNonzeroFloat32: %v\n", math.SmallestNonzeroFloat32)
	fmt.Printf("math.MaxInt8: %v\n", math.MaxInt8)
	fmt.Printf("math.MinInt8: %v\n", math.MinInt8)
	fmt.Printf("math.MaxUint8: %v\n", math.MaxUint8)
	fmt.Printf("math.MaxInt16: %v\n", math.MaxInt16)
	fmt.Printf("math.MaxUint16: %v\n", math.MaxUint16)
	fmt.Printf("math.MaxInt32: %v\n", math.MaxInt32)
	fmt.Printf("math.MinInt32: %v\n", math.MinInt32)
	fmt.Printf("math.MaxUint32: %v\n", math.MaxUint32)
	fmt.Printf("math.MaxInt64: %v\n", math.MaxInt64)
	fmt.Printf("math.MinInt64: %v\n", math.MinInt64)
	fmt.Printf("math.Pi: %v\n", math.Pi)

	// 取绝对值
	fmt.Printf("[-3.14]的绝对值为: [%.2f]\n", math.Abs(-3.14))

	// x的y次方
	fmt.Printf("[2]的16次方为: %v\n", math.Pow(2, 16))

	// 10的三次方
	fmt.Printf("math.Pow10(3): %v\n", math.Pow10(3))

	// x的开平方
	fmt.Printf("math.Sqrt(64): %v\n", math.Sqrt(64))

	// x的开立方
	fmt.Printf("math.Cbrt(27): %v\n", math.Cbrt(27))

	// 向上取整
	fmt.Printf("math.Ceil(3.14): %v\n", math.Ceil(3.14))

	// 向下取整
	fmt.Printf("math.Floor(8.75): %v\n", math.Floor(8.75))

	// 取余数
	fmt.Printf("math.Mod(10, 3): %v\n", math.Mod(10, 3))

	// 取整数与小数部分
	Integer, Decimal := math.Modf(3.1415926)
	fmt.Printf("Integer: %v\n", Integer)
	fmt.Printf("Decimal: %.2f\n", Decimal)

	// seed种子
	rand.Seed(time.Now().UnixMicro())
	for i := 0; i < 10; i++ {
		a := rand.Int()
		fmt.Println(a)
	}
	// 指定范围 100以内
	for i := 0; i < 10; i++ {
		a := rand.Intn(100)
		fmt.Printf("a: %v\n", a)
	}
	for i := 0; i < 10; i++ {
		a := rand.Float32()
		fmt.Println(a)
	}
}
