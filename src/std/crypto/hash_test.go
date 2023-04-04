package crypto

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"testing"
)

func getMD51(src []byte) string {
	// 直接通过sum函数计算数据的md5
	res := md5.Sum(src)
	fmt.Println(res)

	// 将md5结果格式化成16进制格式的字符串
	resStr := fmt.Sprintf("%x", res)
	fmt.Println(resStr)

	// 通过hex.EncodeToString函数将md5结果格式化成16进制格式的字符串
	resStr = hex.EncodeToString(res[:])
	fmt.Println(resStr)

	return resStr
}

func getMD52(src ...[]byte) string {
	myHash := md5.New()
	for _, v := range src {
		myHash.Write(v)
	}

	result := myHash.Sum(nil)
	fmt.Println(result)

	resStr := fmt.Sprintf("%x", result)
	fmt.Println(resStr)

	resStr = hex.EncodeToString(result)
	fmt.Println(resStr)

	return resStr
}

func TestMD5(t *testing.T) {
	getMD51([]byte("Hello, SparkMlLib"))
	getMD52([]byte("Hello, Docker"), []byte("Hello, kubernetes"), []byte("Hello, Kafka"))
}

func getSha1(fileName string) (result string) {
	file, _ := os.Open(fileName)
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	myHash := sha1.New()
	_, _ = io.Copy(myHash, file)

	result = hex.EncodeToString(myHash.Sum(nil))
	return
}

func getSha256(fileName string) (result string) {
	file, _ := os.Open(fileName)
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	myHash := sha256.New()
	_, _ = io.Copy(myHash, file)

	result = hex.EncodeToString(myHash.Sum(nil))
	return
}

func getSha512(fileName string) (result string) {
	file, _ := os.Open(fileName)
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	myHash := sha512.New()
	_, _ = io.Copy(myHash, file)

	result = hex.EncodeToString(myHash.Sum(nil))
	return
}

func TestSha(t *testing.T) {
	fmt.Println(getSha1("Hello, Python"))
	fmt.Println(getSha256("Hello Scala"))
	fmt.Println(getSha512("Hello Golang"))
}
