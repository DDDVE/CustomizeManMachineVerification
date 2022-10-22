package utils

// 随机字符串的来源
const RandomStringSource = "0123456789abcdefghigklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ+/"

// []byte转随机字符串
func BytestoString(b []byte) string {
	res := make([]byte, len(b))
	for i := 0; i < len(b); i++ {
		Index := int(b[i]) % len(RandomStringSource)
		res[i] = RandomStringSource[Index]
	}
	return string(res)
}

// 生成随机字符串
// func GetRandomString(length int) string {
// 	//   用crypto/rand生成随机数更随机但速度稍慢，密钥可以直接存
// 	//   []byte类型，加解密也都是用[]byte,就不用来回转化了
// 	randomBytes := make([]byte, length)
// 	_, err := rand.Read(randomBytes)
// 	if err != nil {
// 		log.Println("生成随机数报错", err)
// 	}
// 	//   return randomBytes
// 	// 返回string用下面代码
// 	res := make([]byte, length)
// 	for i := 0; i < length; i++ {
// 		Index := int(randomBytes[i]) % len(RandomStringSource)
// 		res[i] = RandomStringSource[Index]
// 	}
// 	return string(res)
// }
