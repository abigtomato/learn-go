package bufio

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

func TestWriteToBuffer(t *testing.T) {
	sr := strings.NewReader("Hello World")
	br := bufio.NewReader(sr)

	// 写入缓冲区
	buf := bytes.NewBuffer(make([]byte, 0))
	_, _ = br.WriteTo(buf)
	result, _ := io.ReadAll(buf)
	fmt.Println(string(result))
}

func TestWriteToFile(t *testing.T) {
	sr := strings.NewReader("Hello World")
	br := bufio.NewReader(sr)

	// 写入文件
	file, _ := os.OpenFile("text.txt", os.O_RDWR, 0777)
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	n, _ := br.WriteTo(file)
	result, _ := io.ReadAll(file)
	fmt.Println(n, string(result))
}

func TestWriteString(t *testing.T) {
	file, _ := os.OpenFile("text.txt", os.O_RDWR, 0777)
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	fbw := bufio.NewWriter(file)
	// 写入字符串
	_, _ = fbw.WriteString("Hello Golang")
	// 刷新缓冲区
	_ = fbw.Flush()
}

func TestReset(t *testing.T) {
	buf := bytes.NewBuffer(make([]byte, 0))
	buf2 := bytes.NewBuffer(make([]byte, 0))

	bw := bufio.NewWriter(buf)
	_, _ = bw.WriteString("Hello World")

	// 重置缓冲区
	bw.Reset(buf2)
	_, _ = bw.WriteString("Hello Golang")
	_ = bw.Flush()

	result, _ := io.ReadAll(buf2)
	fmt.Println(string(result))
}

func TestAvailableAndBuffered(t *testing.T) {
	buf := bytes.NewBuffer(make([]byte, 0))
	bw := bufio.NewWriter(buf)
	// 4096默认缓冲区大小 0
	fmt.Println(bw.Available(), bw.Buffered())

	_, _ = bw.WriteString("Hello Golang")
	// 4084 缓冲区剩余大小 12 写入缓冲区大小
	fmt.Println(bw.Available(), bw.Buffered())

	_ = bw.Flush()
	// 4096 0
	fmt.Println(bw.Available(), bw.Buffered())

	all, _ := io.ReadAll(buf)
	fmt.Println(string(all))
}

func TestWriteByteAndWriteRune(t *testing.T) {
	buf := bytes.NewBuffer(make([]byte, 0))
	bw := bufio.NewWriter(buf)

	// 写入缓存 一个一个写 byte等同于int8
	_ = bw.WriteByte('H')
	_ = bw.WriteByte('e')
	_ = bw.WriteByte('l')
	_ = bw.WriteByte('l')
	_ = bw.WriteByte('o')
	_ = bw.WriteByte(' ')
	_, _ = bw.WriteRune('世')
	_, _ = bw.WriteRune('界')
	_, _ = bw.WriteRune('！')

	_ = bw.Flush()
	fmt.Println(fmt.Sprintf("%s", buf))
}

func TestReadWriter(t *testing.T) {
	buf := bytes.NewBuffer(make([]byte, 0))
	bw := bufio.NewWriter(buf)

	br := bufio.NewReader(strings.NewReader("Hello Golang"))

	rw := bufio.NewReadWriter(br, bw)
	str, _ := rw.ReadString('\n')
	fmt.Println(str)

	_, _ = rw.WriteString("ABC")
	_ = rw.Flush()
	fmt.Println(buf)
}

func TestSplitAndScanWords(t *testing.T) {
	sr := strings.NewReader("ABC DEF GHI JKL")
	bs := bufio.NewScanner(sr)

	// 通过空格来分割
	bs.Split(bufio.ScanWords)

	// 扫描
	for bs.Scan() {
		fmt.Println(bs.Text())
	}
}

func TestScanBytesAndScanRunes(t *testing.T) {
	sr := strings.NewReader("Hello 世界")
	bs := bufio.NewScanner(sr)

	bs.Split(bufio.ScanRunes)

	for bs.Scan() {
		fmt.Println(bs.Text())
	}
}
