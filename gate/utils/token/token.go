package token

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"gate/utils"
	"gate/utils/log"
	"strings"
)

var (
	MySignedKey = [NUM_OF_KEYS][LEN_OF_KEY]byte{} //token密钥
)

const (
	// 密钥数量
	NUM_OF_KEYS = 2
	// 密钥长度
	LEN_OF_KEY = 32

	TOKEN_CONNECTER = "@-@"
)

// 生成随机密钥
func GetRandomKey() {
	// 用crypto/rand生成随机数更随机但速度稍慢，密钥可以直接存
	randomBytes := make([]byte, LEN_OF_KEY)
	_, err := rand.Read(randomBytes)
	if err != nil {
		log.Println("生成随机数报错", err)
	}
	for i := NUM_OF_KEYS - 1; i >= 0; i-- {
		if i == 0 {
			copy(MySignedKey[0][:], randomBytes)
		} else {
			MySignedKey[i] = MySignedKey[i-1]
		}
	}
}

func GenerateToken(mobileNum string) string {
	if mobileNum == "" {
		return ""
	}
	// 不可逆加密
	mn := md5.Sum([]byte(mobileNum))
	k := [LEN_OF_KEY]byte{}
	index := 0
	for i := 0; i < LEN_OF_KEY; i++ {
		if index >= len(mn) {
			index = 0
		}
		k[i] = mn[index] ^ MySignedKey[0][i]
		index++
	}

	preToken := md5.Sum(k[:len(k)/2])
	sufToken := md5.Sum(k[len(k)/2:])
	s := preToken[:]
	s = append(s, sufToken[:]...)
	//转化token
	token := utils.BytestoString(s)
	//拼接手机号
	mobile := base64.StdEncoding.EncodeToString([]byte(mobileNum))

	return token + TOKEN_CONNECTER + mobile
}

//检查token并返回与第几个密钥匹配，如果未通过检验则返回-1
func CheckToken(token string) (id int, mobile_num string) {
	t := strings.Split(token, TOKEN_CONNECTER)
	if len(t) != 2 {
		return -1, ""
	}
	mobileNum, err := base64.StdEncoding.DecodeString(t[1])
	if err != nil {
		return -1, ""
	}
	for i := 0; i < NUM_OF_KEYS; i++ {
		if compareToken(t[0], mobileNum, i) {
			return i, string(mobileNum)
		}
	}
	return -1, ""
}

func compareToken(token string, mobileNum []byte, num int) bool {
	if num < 0 || num >= NUM_OF_KEYS {
		return false
	}
	mn := md5.Sum([]byte(mobileNum))
	k := [LEN_OF_KEY]byte{}
	index := 0
	for i := 0; i < LEN_OF_KEY; i++ {
		if index >= len(mn) {
			index = 0
		}
		k[i] = mn[index] ^ MySignedKey[num][i]
		index++
	}

	preToken := md5.Sum(k[:len(k)/2])
	sufToken := md5.Sum(k[len(k)/2:])
	s := preToken[:]
	s = append(s, sufToken[:]...)
	if token != utils.BytestoString(s) {
		return false
	}

	return true
}

// type MyClaim struct {
// 	MobileNumber string `json:"mobile_num"`
// 	jwt.StandardClaims
// }

// // 这个Token和Claim做成全局唯一的变量在并发环境下会导致用户的token模具被其他协程修改
// // var Token *jwt.Token
// // var Claim *MyClaim
// var (
// 	MySignedKey0 = "0123456789makestring/LEOFKE1245Y0123456789makestring/LEOFKE1245Y"
// )

// // TODO: 把这些字段做成从本地文件里读取，不写入代码

// const (
// 	JwtClaimIssuer  = "dve"
// 	JwtClaimSubject = "custom man-machine verify plat"
// )

// func CreateToken(mobileNumber string) *jwt.Token {
// 	// 首先初始化MyClaim
// 	claim := &MyClaim{
// 		MobileNumber: mobileNumber,
// 		StandardClaims: jwt.StandardClaims{
// 			Issuer:  JwtClaimIssuer,
// 			Subject: JwtClaimSubject,
// 		},
// 	}
// 	// 生成token对象
// 	// 这里要采用对称加密否则会报错
// 	return jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
// }

// // 根据最新的密钥生成token并返回
// func GenerateTokenString(mobileNumber string) (string, error) {
// 	token := CreateToken(mobileNumber)
// 	s, err := token.SignedString([]byte(MySignedKey0))
// 	if err != nil {
// 		log.Println("生成token报错: ", err)
// 		return "", errors.New("获取token失败")
// 	}
// 	return s, nil
// }

// // 用户传入的token是否合法
// func CheckTokenString(userToken string) string {
// 	token, _ := jwt.ParseWithClaims(userToken, &MyClaim{}, func(t *jwt.Token) (interface{}, error) {
// 		return []byte(MySignedKey0), nil
// 	})
// 	claims, ok := token.Claims.(*MyClaim)
// 	if ok && token.Valid {
// 		//已经验证成功，说明token没问题，拿到的手机号也是加密时传入的手机号
// 		return claims.MobileNumber
// 	}
// 	return ""
// }
