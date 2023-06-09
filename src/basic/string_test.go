package basic

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	"unicode/utf8"
)

// 字符串遍历
func TestStringIterator(t *testing.T) {
	str := "Docker/Kubernetes云原生技术"

	// 通过下标取值遍历（每次遍历一个字符，也就是一个字节的数据，无法处理多字节的中文）
	for i := 0; i < len(str); i++ {
		fmt.Printf("(%d, %c)", i, str[i])
	}
	fmt.Println()

	// 通过[]byte()强转成字节数组，遍历字节数组获取的是每个字节的utf-8码
	for i, b := range []byte(str) {
		fmt.Printf("(%d, %X)", i, b)
	}
	fmt.Println()

	// range遍历字符串，获取的是每个字符的utf-8码（中文占多个字节）
	for i, ch := range str {
		fmt.Printf("(%d, %X)", i, ch)
	}
	fmt.Println()

	// 转码遍历
	bytes := []byte(str)
	for len(bytes) > 0 {
		// DecodeRune() utf-8 解码函数，返回解码后的字符和字符长度
		ch, size := utf8.DecodeRune(bytes)
		bytes = bytes[size:]
		fmt.Printf("%c", ch)
	}
	fmt.Println()

	// 使用[]rune()转换为切片，遍历切片会自动处理多字节的中文
	for i, ch := range []rune(str) {
		fmt.Printf("(%d, %c)", i, ch)
	}
	fmt.Println()
}

// 字符串转换
func TestStringInternal(t *testing.T) {
	// 基本类型转字符串
	str := fmt.Sprintf("%d", 666)
	fmt.Printf("str type=%T, str=%q\n", str, str)

	// 基本类型转字符串 'f'格式，10位小数，64位
	str = strconv.FormatFloat(127.00245, 'f', 10, 64)
	fmt.Printf("str type=%T, str=%q\n", str, str)

	// string转各种基本数据类型
	if num, err := strconv.ParseInt("233", 10, 64); err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("num type=%T, num=%d\n", num, num)
	}

	// string转整数
	if num, err := strconv.Atoi("233"); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(num)
	}

	// 整数转string
	str = strconv.Itoa(666)
	fmt.Printf("str=%v, type=%T\n", str, str)

	// 判断字符串中是否包含指定的子串
	b := strings.Contains("Docker/Kubernetes", "Docker")
	fmt.Printf("b=%v\n", b)

	// 判断字符串中包含多少指定的子串
	n := strings.Count("Spark/SparkSQL/SparkStreaming", "Spark")
	fmt.Printf("n=%v\n", n)

	// 不区分大小写比较字符串是否相等
	b = strings.EqualFold("abc", "ABC")
	fmt.Printf("b=%v\n", b)

	// 判断子串在字符串中第一次出现的位置，没有返回-1
	i := strings.Index("SparkMLlib", "Spark")
	fmt.Printf("i=%v\n", i)

	// 判断子串在字符串中最后一次出现的位置，没有返回-1
	i = strings.LastIndex("Go Golang", "Go")
	fmt.Printf("i=%v\n", i)

	// 将字符串中的指定子串替换，-1表示全部替换，替换后返回一个新串
	str = strings.Replace("go go hello", "go", "golang", -1)
	fmt.Printf("str=%v\n", str)

	// 字符串切割，返回新数组
	arr := strings.Split("hello,world", ",")
	fmt.Printf("arr=%v\n", arr)

	// 字符串大小写转换
	str = strings.ToLower("GOPATH")
	fmt.Printf("str=%v\n", str)
	str = strings.ToUpper("goroot")
	fmt.Printf("str=%v\n", str)

	// 去除字符串左右两端的空格
	str = strings.TrimSpace(" tn a lone gopher ntrn ")
	fmt.Printf("str=%v\n", str)

	// 去除字符串左右两边的指定字符
	str = strings.Trim("! hello !", " !")
	fmt.Printf("str=%v\n", str)

	// 判断字符串是否以指定的字符串开头
	b = strings.HasPrefix("ftp://192.168.25.130", "ftp")
	fmt.Printf("b=%v\n", b)

	// 判断字符串是否以指定的字符串结尾
	b = strings.HasSuffix("03795E.jpg", ".jpg")
	fmt.Printf("b=%v\n", b)
}
