package utils

import (
	"errors"
	"log"
	"math/rand"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type MyClaim struct {
	MobileNumber string `json:"mobile_num"`
	jwt.StandardClaims
}

// 这个Token和Claim做成全局唯一的变量在并发环境下会导致用户的token模具被其他协程修改
// var Token *jwt.Token
// var Claim *MyClaim

var (
	//MySignedKey = []byte("dongyongwei&dve")
	MySignedKey = make([]string, LenOfMySignedKey)
)

const (
	// 最多几天以前的token有效
	LenOfMySignedKey = 2
	// 平台密钥长度
	LenOfKey = 50
	// 一天的秒数
	SecondsOfDay int64 = 86400
)

// todo: 把这些字段做成从本地文件里读取，不写入代码
const (
	JwtClaimIssuer  = "dve"
	JwtClaimSubject = "custom man-machine verify plat"
)

func InitFirstKey() {
	// 项目初始化时生成第一天的平台密钥
	MySignedKey[0] = GetRandomString(LenOfKey)
	log.Println("初始化平台密钥: ", MySignedKey[0])
}

func CreateToken(mobileNumber string) *jwt.Token {
	// 首先初始化MyClaim
	claim := &MyClaim{
		MobileNumber: mobileNumber,
		StandardClaims: jwt.StandardClaims{
			Issuer:  JwtClaimIssuer,
			Subject: JwtClaimSubject,
		},
	}
	// 生成token对象
	// 这里要采用对称加密否则会报错
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
}

// 根据最新的密钥生成token并返回
func GenerateTokenString(mobileNumber string) (string, error) {
	token := CreateToken(mobileNumber)
	index := LenOfMySignedKey - 1
	for index >= 0 && MySignedKey[index] == "" {
		index--
	}
	if index < 0 {
		log.Println("平台密钥不存在！")
		return "", errors.New("获取token失败")
	}
	s, err := token.SignedString([]byte(MySignedKey[index]))
	if err != nil {
		log.Println("生成token报错: ", err)
		return "", errors.New("获取token失败")
	}
	return s, nil
}

// 用户传入的token是否合法
func CheckTokenString(userToken string, mobileNum string) bool {
	token := CreateToken(mobileNum)
	for i := 0; i < LenOfMySignedKey; i++ {
		if MySignedKey[i] == "" {
			continue
		}
		ss, err := token.SignedString([]byte(MySignedKey[i]))
		if err != nil {
			log.Println("判断token时报错: ", err)
			return false
		}
		if userToken == ss {
			return true
		}
	}
	log.Println("token不正确")
	return false
}

// 随机字符串的来源
const randomStringSource = "0123456789abcdefghigklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// 生成随机字符串
func GetRandomString(length int) string {
	//plus := rand.Int63n(SecondsOfDay)
	rand.Seed(time.Now().UnixNano())
	res := make([]byte, length)
	for i := 0; i < length; i++ {
		Index := rand.Intn(len(randomStringSource))
		res[i] = randomStringSource[Index]
	}
	return string(res)
}
