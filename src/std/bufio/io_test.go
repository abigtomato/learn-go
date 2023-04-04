package bufio

import (
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

func TestIO(t *testing.T) {
	sr := strings.NewReader("hello world")

	// 读取10个字符
	buf := make([]byte, 10)
	_, _ = sr.Read(buf)
	fmt.Println(string(buf))

	_, _ = sr.Seek(0, 0)
	// 把 reader 的内容拷贝到 Stdout 标准输出，每次拷贝一个字节
	_, _ = io.Copy(os.Stdout, sr)
	fmt.Println()

	_, _ = sr.Seek(0, 0)
	// 让程序一次读取8个字节
	buf = make([]byte, 8)
	_, _ = io.CopyBuffer(os.Stdout, sr, buf)
	fmt.Println()

	file, _ := os.Open("text.txt")
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	// 读取文件中的全部内容
	bytes, _ := io.ReadAll(file)
	fmt.Println(string(bytes))

	// 当前程序目录下的所有文件
	dir, _ := os.ReadDir(".")
	for _, d := range dir {
		fmt.Println(d.Name())
	}

	content := []byte("temporary file's content")
	// 创建临时文件
	tempFile, _ := os.CreateTemp("", "example")
	fmt.Printf("tempFile.Name(): %v\n", tempFile.Name())
	defer func(name string) {
		_ = os.Remove(name)
	}(tempFile.Name())
	// 向临时文件写入内容
	_, _ = tempFile.Write(content)
}
