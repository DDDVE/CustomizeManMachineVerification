package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/pem"
	"log"
)

// 用公钥加密
func Encrypt(plain string, publicKey string) (encrypted []byte, e error) {
	log.Println("进入加密函数")
	msg := []byte(plain)
	// 解码公钥
	pubBlock, _ := pem.Decode([]byte(publicKey))
	// 读取公钥
	pubKeyValue, err := x509.ParsePKIXPublicKey(pubBlock.Bytes)
	if err != nil {
		log.Println("读取公钥报错: ", err)
		return nil, err
	}
	pub := pubKeyValue.(*rsa.PublicKey)
	// 加密数据
	encryptToAEP, err := rsa.EncryptOAEP(sha1.New(), rand.Reader, pub, msg, nil)
	if err != nil {
		log.Println("加密过程报错: ", err)
		return nil, err
	}
	return encryptToAEP, nil
}
