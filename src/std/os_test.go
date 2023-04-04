package std

import (
	"fmt"
	"os"
	"testing"
)

func TestOS(t *testing.T) {
	// O_RDWR|O_CREATE|O_TRUNC, 0666
	// 可读写 自动创建 清空内容 存放在 go.mod 定义的模块下
	file, _ := os.Create("test.txt")
	_, _ = file.WriteString("hello golang")

	// os.ModePerm 0o777 最高权限
	_ = os.Mkdir("demo", os.ModePerm)
	// 多级目录创建
	_ = os.MkdirAll("a/b/c", os.ModePerm)

	// 删除目录或文件
	_ = os.Remove("demo")
	// 删除目录及里面的所有内容
	_ = os.RemoveAll("a")

	// 获取工作目录
	dir, _ := os.Getwd()
	fmt.Printf("dir: %v\\n", dir)

	// 修改工作目录
	_ = os.Chdir("d:/")
	fmt.Println(os.Getwd())

	// 获取临时目录
	tempDir := os.TempDir()
	fmt.Printf("tempDir: %v\\n", tempDir)

	// 重命名文件
	_ = os.Rename("test.txt", "test1.txt")

	// 读文件
	bytes, _ := os.ReadFile("test1.txt")
	fmt.Printf("string(b[:]): %v\\n", string(bytes[:]))

	// 写文件 如果没有该文件会创建 如果文件有内容会覆盖
	_ = os.WriteFile("test.txt", []byte("Hello Golang"), os.ModePerm)

	// os.O_RDWR会把文件前面对应的长度内容覆盖 后面的内容不覆盖
	// os.O_APPEND在后面追加
	file, _ = os.OpenFile("test.txt", os.O_RDWR|os.O_APPEND, 0755)
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	// 从第三个位置开始写入 覆盖对应长度的内容 在后面的不覆盖
	n, _ := file.WriteAt([]byte("rust"), 3)
	fmt.Printf("n: %v\n", n)
}
