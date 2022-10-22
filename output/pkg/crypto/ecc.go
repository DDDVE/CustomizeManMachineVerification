package crypto

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"math/big"
	"os"
	"pkg/conf"
)

const (
	PUBLIC_KEY_PATH = "../pkg/crypto/eccpublic.pem"
)

//取得ECC公钥
func GetECCPublicKey() (*ecdsa.PublicKey, error) {
	//读取公钥
	file, err := os.Open(PUBLIC_KEY_PATH)
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

//验证数字签名
func VerifySignECC(hash, rbytes, sbytes []byte) error {
	//读取公钥
	publicKey, err := GetECCPublicKey()
	if err != nil {
		return err
	}
	//计算哈希值
	// hash := sha256.New()
	// _, err = hash.Write([]byte(msg))
	// if err != nil {
	// 	return false
	// }
	// bytes := hash.Sum(nil)

	var r, s big.Int
	r.UnmarshalText(rbytes)
	s.UnmarshalText(sbytes)
	//验证数字签名
	verify := ecdsa.Verify(publicKey, hash, &r, &s)
	if !verify {
		return errors.New(conf.GlobalError[conf.SIGNATURE_CHECK_ERROR])
	}
	return nil
}
