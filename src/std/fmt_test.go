package std

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestPrintf(t *testing.T) {
	type User struct {
		Id int64
	}
	user := &User{Id: 1}
	// 值的默认格式
	fmt.Printf("%v\n", user)
	// 类似%v，但会添加字段名
	fmt.Printf("%+v\n", user)
	// 值的Go语法表示
	fmt.Printf("%#v\n", user)
	// 值的类型
	fmt.Printf("%T\n", user)
	// 百分号
	fmt.Printf("%%\n")
	// 布尔值
	fmt.Printf("%t\n", true)

	n := 180
	// 二进制
	fmt.Printf("%b\n", n)
	// Unicode码
	fmt.Printf("%c\n", n)
	// 十进制
	fmt.Printf("%d\n", n)
	// 八进制
	fmt.Printf("%o\n", n)
	// 十六进制，使用a-f
	fmt.Printf("%x\n", n)
	// 十六进制，使用A-F
	fmt.Printf("%X\n", n)
	// Unicode格式：U+1234
	fmt.Printf("%U\n", n)
	// 单引号括起来的Go语法字符字面值
	fmt.Printf("%q\n", n)

	f := 3.14
	// 无小数部分、二进制指数的科学计数法，如-123456p-78
	fmt.Printf("%b\n", f)
	// 科学计数法，如-1234.456e+78
	fmt.Printf("%e\n", f)
	// 科学计数法，如-1234.456E+78
	fmt.Printf("%E\n", f)
	// 有小数部分但无指数部分，如123.456
	fmt.Printf("%f\n", f)
	// 等价于%f
	fmt.Printf("%F\n", f)
	// 根据实际情况采用%e或%f格式
	fmt.Printf("%g\n", f)
	// 根据实际情况采用%E或%F格式
	fmt.Printf("%G\n", f)

	s := "Hello Golang"
	// 直接输出字符串或[]byte
	fmt.Printf("%s\n", s)
	// 该值对应的双引号括起来的Go语法字符串字面量，必要时会安全转义
	fmt.Printf("%q\n", s)
	// 每个字节用两字符十六进制数表示，使用a-f
	fmt.Printf("%x\n", s)
	// 每个字节用两字符十六进制数表示，使用A-F
	fmt.Printf("%X\n", s)

	// 宽度9，默认精度
	fmt.Printf("%10f\n", f)
	// 默认宽度，精度2
	fmt.Printf("%.2f\n", f)
	// 宽度9，精度2
	fmt.Printf("%10.2f\n", f)
	// 宽度9，精度0
	fmt.Printf("%10f.\n", f)
}

func TestFPrint(t *testing.T) {
	file, _ := os.OpenFile("test.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	n, _ := fmt.Fprintf(file, "Hello %s\n", "Golang")
	fmt.Println(n)
}

func TestSprint(t *testing.T) {
	host := "localhost"
	port := 6379
	addr := fmt.Sprintf("%s:%d", host, port)
	fmt.Println(addr)
}

func TestErrorf(t *testing.T) {
	_ = fmt.Errorf("error message is %s", "golang")
}

func TestFScan(t *testing.T) {
	var (
		name    string
		age     int
		married bool
	)
	reader := strings.NewReader("1:albert 2:25 3:true")
	_, _ = fmt.Fscanf(reader, "1:%s 2:%d 3:%t", &name, &age, &married)
	fmt.Printf("name = %s, age = %d, married = %t", name, age, married)
}
