package basic

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"testing"
)

// const (
// 	 O_RDONLY int = syscall.O_RDONLY // 只读模式打开文件
// 	 O_WRONLY int = syscall.O_WRONLY // 只写模式打开文件
// 	 O_RDWR   int = syscall.O_RDWR   // 读写模式打开文件
//	 O_APPEND int = syscall.O_APPEND // 写操作时将数据附加到文件尾部
//	 O_CREATE int = syscall.O_CREAT  // 如果不存在将创建一个新文件
//	 O_EXCL   int = syscall.O_EXCL   // 和O_CREATE配合使用，文件必须不存在
//	 O_SYNC   int = syscall.O_SYNC   // 打开文件用于同步I/O
//	 O_TRUNC  int = syscall.O_TRUNC  // 如果可能，打开时清空文件
// )
//
// const (
//	 // 单字符是被String方法用于格式化的属性缩写。
//	 ModeDir        FileMode = 1 << (32 - 1 - iota) // d: 目录
//	 ModeAppend                                     // a: 只能写入，且只能写入到末尾
//	 ModeExclusive                                  // l: 用于执行
//	 ModeTemporary                                  // T: 临时文件（非备份文件）
//	 ModeSymlink                                    // L: 符号链接（不是快捷方式文件）
//	 ModeDevice                                     // D: 设备
//	 ModeNamedPipe                                  // p: 命名管道（FIFO）
//	 ModeSocket                                     // S: Unix域socket
//	 ModeSetuid                                     // u: 表示文件具有其创建者用户id权限
//	 ModeSetgid                                     // g: 表示文件具有其创建者组id的权限
//	 ModeCharDevice                                 // c: 字符设备，需已设置ModeDevice
//	 ModeSticky                                     // t: 只有root/创建者能删除/移动文件
//	 // 覆盖所有类型位（用于通过&获取类型位），对普通文件，所有这些位都不应被设置
//	 ModeType          = ModeDir | ModeSymlink | ModeNamedPipe | ModeSocket | ModeDevice
//	 ModePerm FileMode = 0777 // 覆盖所有Unix权限位（用于通过&获取类型位）
// )

// 缓冲区读取方式，不会一次性加载到内存中
func TestBufRead(t *testing.T) {
	// 只读方式打开文件，获取文件句柄
	file, _ := os.Open("")
	// 文件句柄就是指向为文件预分配的一块内存的指针
	fmt.Printf("指向=%p, 类型=%T\n", file, file)
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	// 创建带缓冲区(用户缓冲)的Reader指针用于读取文件
	// 预读入: 预先读取一定大小数据缓冲，之后由程序读取
	// 缓输出: 程序写入系统缓冲区后就结束操作，之后由操作系统算法决定系统缓冲区何时写入磁盘（一般会攒一段时间数据，这些数据可以是不同程序的写入）
	reader := bufio.NewReader(file)
	for {
		// 游标指向文件开头
		// 字符串的形式读取，出现换行符结束本次读取，游标移动换行符之后，下一次从此处开始读取
		str, err := reader.ReadString('\n')
		// io.EOF表示文件读取到末尾
		if err == io.EOF {
			break
		}
		fmt.Printf("%v", str)
	}
}

// 一次性加载进内存中的方式读取
func TestReadAll(t *testing.T) {
	content, _ := os.ReadFile("")
	fmt.Printf("%v\n", string(content))
}

// 缓冲区写入方式，数据先写入内存缓冲区（默认4096字节），缓冲区满了则溢写进磁盘文件
func TestBufWriter(t *testing.T) {
	// 参数1表示文件路径，参数2表示打开方式，参数3表示权限控制(适用于Linux/Unix)
	// os.O_APPEND表示追加，os.O_CREATE表示不存在则创建，|表示连接多个方式
	file, _ := os.OpenFile("", os.O_APPEND|os.O_CREATE, 0666)
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	// 创建新的缓冲区写入指针
	writer := bufio.NewWriter(file)
	for i := 0; i < 1; i++ {
		// 写入字符串到内存缓冲区
		_, _ = writer.WriteString("\r\nElasticSearch")
	}
	// 将缓冲区中的数据写入磁盘文件
	_ = writer.Flush()
}

// 一次性加载内存方式拷贝文件
func TestAllCopyFile(t *testing.T) {
	data, _ := os.ReadFile("")
	_ = os.WriteFile("", data, 0666)
}

// 缓冲区方式拷贝文件
func TestBufCopyFile(t *testing.T) {
	fromFile, _ := os.Open("")
	defer func(fromFile *os.File) {
		_ = fromFile.Close()
	}(fromFile)

	toFile, _ := os.OpenFile("", os.O_WRONLY|os.O_CREATE, 0666)
	defer func(toFile *os.File) {
		_ = toFile.Close()
	}(toFile)

	reader := bufio.NewReader(fromFile)
	writer := bufio.NewWriter(toFile)

	_, _ = io.Copy(writer, reader)
}

// 以字节为单位拷贝文件
func TestByteCopyFile(t *testing.T) {
	fr, _ := os.Open("")
	defer func(fr *os.File) {
		_ = fr.Close()
	}(fr)

	fw, _ := os.Create("")
	defer func(fw *os.File) {
		_ = fw.Close()
	}(fw)

	buf := make([]byte, 4096)
	for {
		n, err := fr.Read(buf)
		if err == io.EOF {
			break
		}
		_, _ = fw.Write(buf[:n])
	}
}

// 写入文件
func TestWriteFile(t *testing.T) {
	f, _ := os.OpenFile("", os.O_WRONLY|os.O_CREATE, 0666)
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	// WriteString(): 向文件写入字符串
	// 参数: 需要写入文件的字符串，\r\n表示换行符
	// 返回值: 写入文件的字节数和error
	n, _ := f.WriteString("HelloGolang\r\n")
	fmt.Printf("WriteString to %dbyte", n)

	// Seek(): 按字节移动读写指针
	// 参数1: 偏移量，以字节为单位
	// 参数2: 偏移的起始位置，SeekEnd文件结束位置，SeekStart文件起始位置，SeekCurrent文件当前位置
	// 返回值: 表示从文件起始位置，到当前读写指针位置的偏移量和error
	off, _ := f.Seek(5, io.SeekStart)
	fmt.Printf("Seek Off: %v\n", off)

	// WriteAt(): 从偏移位置开始写入数据到文件
	// 参数1: 需要写入文件的字节数组
	// 参数2: 开始写入的偏移量
	// 返回值:
	n, _ = f.WriteAt([]byte("Python"), off)
	fmt.Printf("WriteAt to %dbyte", n)
}

// 目录操作
func TestDir(t *testing.T) {
	f, _ := os.OpenFile("", os.O_RDONLY, os.ModeDir)
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	// Readdir(): 读取目录中所有目录项
	// 参数: 读取的目录项数量，-1代表读取全部
	// 返回值: 返回fileInfo接口类型的切片
	info, _ := f.Readdir(-1)
	for _, fileInfo := range info {
		// IsDir(): fileInfo提供的判断当前目录项是否是目录的方法
		if fileInfo.IsDir() {
			// Name(): fileInfo提供的获取当前目录项名字的方法
			fmt.Printf("%v是一个目录", fileInfo.Name())
		} else {
			fmt.Printf("%v是一个文件", fileInfo.Name())
		}
	}
}

// 用于存储字符统计结果的结构体
type CharCount struct {
	ChCount    int
	NumCount   int
	SpaceCount int
	OtherCount int
}

// 用于描述CharCount的内容
func (count *CharCount) String() string {
	return fmt.Sprintf("英文字母个数=%v, 数字个数=%v, 空格制表符个数=%v, 其他字符个数=%v\n",
		count.ChCount, count.NumCount, count.SpaceCount, count.OtherCount)
}

// 统计文件中各种字符出现的次数
func Count(path string, count *CharCount) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("err=%v\n", err)
		return
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	reader := bufio.NewReader(file)
	for {
		str, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}

		for _, v := range []rune(str) {
			switch {
			case v >= 'a' && v <= 'z':
				fallthrough
			case v >= 'A' && v <= 'Z':
				count.ChCount++
			case v == ' ' || v == '\t':
				count.SpaceCount++
			case v >= 0 && v <= 9:
				count.NumCount++
			default:
				count.OtherCount++
			}
		}
	}
}

// 文件统计案例
func TestFileCharCount(t *testing.T) {
	var cc CharCount
	Count("", &cc)
	fmt.Println(cc.String())
}
