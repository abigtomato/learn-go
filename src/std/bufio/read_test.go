package bufio

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

func TestNewReader(t *testing.T) {
	file, _ := os.Open("text.txt")
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	// 创建缓冲区读取器
	fr := bufio.NewReader(file)
	str, _ := fr.ReadString('\n')
	fmt.Println(str)
}

func TestBufRead(t *testing.T) {
	// 控制每次读取的缓冲区大小
	sr := strings.NewReader("Hello World")
	br := bufio.NewReader(sr)
	buf := make([]byte, 10)
	for {
		n, err := br.Read(buf)
		if err == io.EOF {
			break
		} else {
			fmt.Println(string(buf[0:n]))
		}
	}
}

func TestReadByte(t *testing.T) {
	sr := strings.NewReader("Hello World")
	br := bufio.NewReader(sr)

	// 读取一个字节
	b, _ := br.ReadByte()
	fmt.Println(string(b))

	// 吐出一个
	_ = br.UnreadRune()

	// 继续读取一个
	b, _ = br.ReadByte()
	fmt.Println(string(b))

	// 读取多个字节，到指定字符截止
	bytes, _ := br.ReadBytes('d')
	fmt.Println(string(bytes))
}

func TestReadRune(t *testing.T) {
	sr := strings.NewReader("你好，世界")
	br := bufio.NewReader(sr)

	// 读取中文、日文等UTF-8编码
	r, s, _ := br.ReadRune()
	fmt.Println(string(r), s)
}

func TestReadLine(t *testing.T) {
	sr := strings.NewReader("ABC\nDEF\r\nGHI\r\nGHI")
	br := bufio.NewReader(sr)
	// 读取一行
	line, prefix, _ := br.ReadLine()
	fmt.Println(string(line), prefix)
}

func TestReadString(t *testing.T) {
	sr := strings.NewReader("ABC DEF GHI JKL")
	br := bufio.NewReader(sr)

	// 读取字符串 到指定字符截止
	str, _ := br.ReadString(' ')
	fmt.Println(str)
}
