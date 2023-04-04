package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"fmt"
	"testing"
)

// 使用AES算法，CTR模式进行加解密
func aesCrypto(text, key []byte) []byte {
	block, _ := aes.NewCipher(key)
	iv := []byte("12345678abcdefgh")
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(text, text)
	return text
}

// 对称加密 AES算法
func TestAES(t *testing.T) {
	key := []byte("12345678abcdefgh")
	src := []byte("Golang才是世界上最好的语言")
	cipherText := aesCrypto(src, key)
	fmt.Printf("加密后: %s\n", cipherText)
	plainText := aesCrypto(cipherText, key)
	fmt.Printf("解密后: %s\n", plainText)
}

// 填充函数
func paddingLastGroup(plainText []byte, blockSize int) []byte {
	padNum := blockSize - len(plainText)%blockSize
	char := []byte{byte(padNum)}
	newPlain := bytes.Repeat(char, padNum)
	return append(plainText, newPlain...)
}

// 去除填充数据
func unPaddingLastGroup(plainText []byte) []byte {
	length := len(plainText)
	lastChar := plainText[length-1]
	number := int(lastChar)
	return plainText[:length-number]
}

// 使用DES算法，CBC分组模式加密
func desEncrypt(plainText, key []byte) []byte {
	block, _ := des.NewCipher(key)
	newText := paddingLastGroup(plainText, block.BlockSize())
	iv := []byte("12345678")
	blockModel := cipher.NewCBCEncrypter(block, iv)
	cipherText := make([]byte, len(newText))
	blockModel.CryptBlocks(cipherText, newText)
	return cipherText
}

// 使用DES算法解密
func desDecrypt(cipherText, key []byte) []byte {
	block, _ := des.NewCipher(key)
	iv := []byte("12345678")
	blockModel := cipher.NewCBCDecrypter(block, iv)
	blockModel.CryptBlocks(cipherText, cipherText)
	return unPaddingLastGroup(cipherText)
}

// 对称加密 DES算法
func TestDES(t *testing.T) {
	key := []byte("1234abcd")
	src := []byte("PHP不是世界上最好的语言")
	cipherText := desEncrypt(src, key)
	fmt.Printf("加密后: %s\n", cipherText)
	plainText := desDecrypt(cipherText, key)
	fmt.Printf("解密后: %s\n", plainText)
}
