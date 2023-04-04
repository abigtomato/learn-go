package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"testing"
)

// 生成RSA密钥对，并持久化到磁盘中
func generateRsaKey(keySize int) {
	privateKey, _ := rsa.GenerateKey(rand.Reader, keySize)

	derText := x509.MarshalPKCS1PrivateKey(privateKey)
	block := pem.Block{
		Type:  "rsa private key",
		Bytes: derText,
	}

	file, _ := os.Create("./pem/private.pem")
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	_ = pem.Encode(file, &block)

	publicKey := privateKey.PublicKey
	derStream, _ := x509.MarshalPKIXPublicKey(&publicKey)

	file, _ = os.Create("./pem/public.pem")
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	_ = pem.Encode(file, &pem.Block{
		Type:  "rsa public key",
		Bytes: derStream,
	})
}

// 使用RSA公钥进行加密
func rsaEncrypt(plainText []byte, publicKeyFileName string) []byte {
	file, _ := os.Open(publicKeyFileName)
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	fileInfo, _ := file.Stat()
	buf := make([]byte, fileInfo.Size())
	_, _ = file.Read(buf)

	block, _ := pem.Decode(buf)
	pubInterface, _ := x509.ParsePKIXPublicKey(block.Bytes)

	publicKey, _ := pubInterface.(*rsa.PublicKey)

	cipherText, _ := rsa.EncryptPKCS1v15(rand.Reader, publicKey, plainText)
	return cipherText
}

// 使用RSA私钥进行解密
func rsaDecrypt(cipherText []byte, privateKeyFileName string) []byte {
	file, _ := os.Open(privateKeyFileName)
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	fileInfo, _ := file.Stat()
	buf := make([]byte, fileInfo.Size())
	_, _ = file.Read(buf)

	block, _ := pem.Decode(buf)
	privateKey, _ := x509.ParsePKCS1PrivateKey(block.Bytes)

	plainText, _ := rsa.DecryptPKCS1v15(rand.Reader, privateKey, cipherText)
	return plainText
}

// 非对称加密 RSA
func TestRSA(t *testing.T) {
	generateRsaKey(1024)

	src := []byte("Golang才是世界上最好的语言")

	cipherText := rsaEncrypt(src, "./pem/public.pem")
	fmt.Printf("加密后: %s\n", cipherText)

	plainText := rsaDecrypt(cipherText, "./pem/private.pem")
	fmt.Printf("解密后: %s\n", plainText)
}
