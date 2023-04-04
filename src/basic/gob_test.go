package basic

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"testing"
)

// Person 测试结构
type Person struct {
	Name string
	Age  int
}

// Gob 是Go语言自己以二进制形式序列化和反序列化程序数据的格式
func TestGob(t *testing.T) {
	data := Person{"Albert", 22}

	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	_ = encoder.Encode(&data)

	encode := buffer.Bytes()
	fmt.Printf("encode: %v\n", encode)

	var decode Person
	decoder := gob.NewDecoder(bytes.NewReader(encode))
	_ = decoder.Decode(&decode)

	fmt.Printf("decode: %v\n", decode)
}
