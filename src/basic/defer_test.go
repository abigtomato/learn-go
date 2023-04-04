package basic

import (
	"bufio"
	"fmt"
	"os"
	"testing"
)

func fib(total int) int {
	if total == 1 || total == 2 {
		return 1
	} else {
		return fib(total-1) + fib(total-2)
	}
}

func writeFibToFile(filename string) {
	file, err := os.OpenFile(filename, os.O_EXCL|os.O_CREATE, 0666)
	if err != nil {
		if pathError, ok := err.(*os.PathError); !ok {
			// 停止当前函数执行；一直向上返回，执行每一层的defer；遇到recover会执行
			panic(err)
		} else {
			fmt.Printf("%s, %s, %s\n", pathError.Op, pathError.Path, pathError.Err)
			return
		}
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	writer := bufio.NewWriter(file)
	defer func(writer *bufio.Writer) {
		_ = writer.Flush()
	}(writer)

	for i := 1; i <= 20; i++ {
		_, _ = fmt.Fprintln(writer, fib(i))
	}
}

// defer延迟执行
func TestDefer(t *testing.T) {
	for i := 1; i <= 10; i++ {
		fmt.Println(fib(i))
	}
	writeFibToFile("./data/fib.txt")
}
