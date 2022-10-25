package apiregist

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	MODULE_ID        = 2
	ADDRESS          = "cmmvplat.com:8083"
	PRIVATE_KEY_PATH = "../pkg/conf/apiprivate.pem"

	DIGITAL_SIGNATURE_CONNECTOR = "@==@"
	REQUEST_ADDR                = "https://cmmvplat.com:8888/apiRegist"
)

func ApiRegist() error {

	//读取私钥文件
	file, err := os.Open(PRIVATE_KEY_PATH)
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	bytes := make([]byte, info.Size())
	file.Read(bytes)

	//生成数字签名
	block, _ := pem.Decode(bytes)                   //pem解码
	ecc, err := x509.ParseECPrivateKey(block.Bytes) //x509解码
	if err != nil {
		return err
	}
	timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)
	msg := fmt.Sprintf("%d%v%v", MODULE_ID, ADDRESS, timestamp)
	hash := sha256.Sum256([]byte(msg))                 //sha256生成b摘要
	r, s, err := ecdsa.Sign(rand.Reader, ecc, hash[:]) //用ECC私钥加密摘要获得数字签名
	if err != nil {
		return err
	}

	rbytes, err := r.MarshalText()
	sbytes, err := s.MarshalText()
	if err != nil {
		return err
	}
	digitalSignature := string(rbytes) + DIGITAL_SIGNATURE_CONNECTOR + string(sbytes) //数字签名转字符串

	//拼装get请求url
	url := fmt.Sprintf("%v?moduleID=%v&address=%v&timestamp=%v&ciphertext=%v", REQUEST_ADDR, MODULE_ID, ADDRESS, timestamp, digitalSignature)

	//发送get请求
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if string(d) != "ok" {
		return errors.New(string(d))
	}
	log.Println("向微服务网关注册服务成功")
	return nil
}
