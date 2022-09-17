package utils

import (
	"math/rand"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type MyClaim struct {
	MobileNumber string `json:"mobileNumber"`
	jwt.StandardClaims
}

var Token *jwt.Token
var Claim *MyClaim

var (
	//MySignedKey = []byte("dongyongwei&dve")
	MySignedKey = make([]string, LenOfMySignedKey)
)

const (
	// 最多几天以前的token有效
	LenOfMySignedKey = 2
	// 平台密钥长度
	LenOfKey = 50
)

func InitFirstKey() {
	// 项目初始化时生成第一天的平台密钥
	MySignedKey[0] = GetRandomString(LenOfKey)
}

func CreateToken(mobileNumber string) {
	// 首先初始化MyClaim
	Claim = &MyClaim{
		MobileNumber: mobileNumber,
		StandardClaims: jwt.StandardClaims{
			Issuer:  "dve",
			Subject: "custom man-machine verify plat",
		},
	}
	// 生成token对象
	Token = jwt.NewWithClaims(jwt.SigningMethodES256, Claim)
}

func GenerateTokenString(mobileNumber string) (res []string, e error) {
	CreateToken(mobileNumber)
	for i := 0; i < len(MySignedKey); i++ {
		// 例如此时是项目运行第一天，第二天的位置还没有密钥
		if MySignedKey[i] == "" {
			break
		}
		// 根据平台密钥获取token
		ss, err := Token.SignedString([]byte(MySignedKey[i]))
		if err != nil {
			// 如果其中一个出错，可以将之前成功的部分返回
			return res, err
		}
		res = append(res, ss)
	}
	return res, nil
}

// 随机字符串的来源
const randomStringSource = "0123456789abcdefghigklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// 生成随机字符串
func GetRandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	res := make([]byte, length)
	for i := 0; i < length; i++ {
		Index := rand.Intn(len(randomStringSource))
		res[i] = randomStringSource[Index]
	}
	return string(res)
}
