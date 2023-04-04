package std

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func TestContains(t *testing.T) {
	b1 := []byte("Hello Golang")
	b2 := []byte("Hello Java")
	b3 := []byte("Hello")
	// 是否包含
	fmt.Println(bytes.Contains(b1, b2))
	fmt.Println(bytes.Contains(b2, b3))
}

func TestCount(t *testing.T) {
	b := []byte("hello")
	sep1 := []byte("h")
	sep2 := []byte("l")
	sep3 := []byte("o")
	// 存在多少
	fmt.Println(bytes.Count(b, sep1))
	fmt.Println(bytes.Count(b, sep2))
	fmt.Println(bytes.Count(b, sep3))
}

func TestRepeat(t *testing.T) {
	b := []byte("hello")
	// 重复复制
	fmt.Println(string(bytes.Repeat(b, 1)))
	fmt.Println(string(bytes.Repeat(b, 3)))
}

func TestReplace(t *testing.T) {
	s := []byte("hello world")
	old := []byte("o")
	news := []byte("ee")
	// 替换
	fmt.Println(string(bytes.Replace(s, old, news, 0)))  // 0次 不替换
	fmt.Println(string(bytes.Replace(s, old, news, 1)))  // 1次 替换一个o
	fmt.Println(string(bytes.Replace(s, old, news, 2)))  // 2次 替换2个o
	fmt.Println(string(bytes.Replace(s, old, news, -1))) // -1次 替换若干个个o
}

func TestRunes(t *testing.T) {
	s := []byte("你好世界")
	r := bytes.Runes(s)
	fmt.Println(len(s)) // 12 ASCII编码 汉字算3个字节
	fmt.Println(len(r)) // 4 Unicode编码 汉字算1个字节
}

func TestJoin(t *testing.T) {
	b := [][]byte{[]byte("你好"), []byte("世界")}
	sep1, sep2 := []byte(","), []byte("#")
	// 通过分隔符连接
	fmt.Println(string(bytes.Join(b, sep1)))
	fmt.Println(string(bytes.Join(b, sep2)))
}

func TestBytesReader(t *testing.T) {
	data := "123456789"
	// 通过[]byte创建Reader
	br := bytes.NewReader([]byte(data))
	// 返回未读取部分的长度 9
	fmt.Println(br.Len())
	// 返回底层数据总长度 9
	fmt.Println(br.Size())

	// 缓冲区大小，即每次读取2字节
	buf := make([]byte, 2)
	for {
		// 通过缓冲区读取数据
		n, err := br.Read(buf)
		if err != nil {
			break
		}
		// 12 34 56 78 9
		fmt.Println(string(buf[:n]))
	}

	// 重置偏移量，因为上面操作已经修改了读取位置等信息
	_, _ = br.Seek(0, 0)
	for {
		// 按字节读取数据
		b, err := br.ReadByte()
		if err != nil {
			break
		}
		// 1 2 3 4 5 6 7 8 9
		fmt.Println(string(b))
	}

	_, _ = br.Seek(0, 0)
	off := int64(0)
	for {
		// 指定偏移量读取数据
		n, err := br.ReadAt(buf, off)
		if err != nil {
			break
		}
		off += int64(n)
		// 2 12 4 34 6 56 8 78 1
		fmt.Println(off, string(buf[:n]))
	}
}

func TestBuffer(t *testing.T) {
	var buf bytes.Buffer
	fmt.Println(buf)

	bufStr := bytes.NewBufferString("hello")
	fmt.Println(bufStr)

	buf2 := bytes.NewBuffer([]byte("hello"))
	fmt.Println(buf2)
}

func TestBufferWriteString(t *testing.T) {
	var buf bytes.Buffer
	// 写入缓冲区
	n, _ := buf.WriteString("hello")
	fmt.Println(n, string(buf.Bytes()))
}

func TestBufferRead(t *testing.T) {
	bufStr := bytes.NewBufferString("hello world")
	buf := make([]byte, 2)
	for {
		// 按缓冲区读取
		n, err := bufStr.Read(buf)
		if err == io.EOF {
			break
		}
		fmt.Println(n, string(buf))
	}
}
