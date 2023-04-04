package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"testing"
)

// 生成消息认证码
func generateHMac(plainText, key []byte) []byte {
	// 创建一个采用sha256作为底层hash接口、key作为密钥的HMAC算法的hash接口
	myHash := hmac.New(sha256.New, key)
	// 向hash中添加明文数据
	myHash.Write(plainText)
	// 计算hash结果
	return myHash.Sum(nil)
}

// 验证消息认证码
func verifyHMac(plainText, key, hashText []byte) bool {
	// 创建一个采用sha256作为底层hash接口、key作为密钥的HMAC算法的hash接口
	myHash := hmac.New(sha256.New, key)
	// 向hash中添加明文数据
	myHash.Write(plainText)
	// 计算hmac
	hmacBytes := myHash.Sum(nil)
	// 比较两个hmac是否相同
	return hmac.Equal(hashText, hmacBytes)
}

func TestHmac(t *testing.T) {
	src := []byte("Hello, Golang")
	key := []byte("Hello World")
	hashText := generateHMac(src, key)
	final := verifyHMac(src, key, hashText)
	if final {
		fmt.Println("消息认证码认证成功!!!")
	} else {
		fmt.Println("消息认证码认证失败...")
	}
}
