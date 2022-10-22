package utils

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"math/big"
	"os"
)

//取得ECC公钥
func GetECCPublicKey(path string) (*ecdsa.PublicKey, error) {
	//读取公钥
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	info, _ := file.Stat()
	buf := make([]byte, info.Size())
	file.Read(buf)
	//pem解密
	block, _ := pem.Decode(buf)
	//x509解密
	publicInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	publicKey := publicInterface.(*ecdsa.PublicKey)
	return publicKey, nil
}

//用ECC公钥解密，优于RSA
func VerifySignECC(hash, rbytes, sbytes []byte, path string) error {
	//读取公钥
	publicKey, err := GetECCPublicKey(path)
	if err != nil {
		return err
	}

	var r, s big.Int
	r.UnmarshalText(rbytes)
	s.UnmarshalText(sbytes)
	//验证密文
	verify := ecdsa.Verify(publicKey, hash, &r, &s)
	if !verify {
		return errors.New("验证失败")
	}
	return nil
}

// 用RSA公钥加密
// func Encrypt(plain string, publicKey string) (encrypted []byte, e error) {
// 	log.Println("进入加密函数")
// 	msg := []byte(plain)
// 	// 解码公钥
// 	pubBlock, _ := pem.Decode([]byte(publicKey))
// 	// 读取公钥
// 	pubKeyValue, err := x509.ParsePKIXPublicKey(pubBlock.Bytes)
// 	if err != nil {
// 		log.Println("读取公钥报错: ", err)
// 		return nil, err
// 	}
// 	pub := pubKeyValue.(*rsa.PublicKey)
// 	// 加密数据
// 	encryptToAEP, err := rsa.EncryptOAEP(sha1.New(), rand.Reader, pub, msg, nil)
// 	if err != nil {
// 		log.Println("加密过程报错: ", err)
// 		return nil, err
// 	}
// 	return encryptToAEP, nil
// }
