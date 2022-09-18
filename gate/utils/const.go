package utils

// 统一记录常量

/**
0：查询成功
1：拒绝请求，返回服务器繁忙
41xx:代表请求校验错误
42xx:代表请求通过校验，业务处理时出错
43xx:代表查询数据库错误/未查到数据
44xx:其他
*/
const (
	HttpSucceed = 0
	HttpRefuse  = 1
	// token校验错误
	HttpTokenCheckFalse = 4101
	// 请求方法错误
	HttpMethodCheckFalse = 4102
	// 请求路径错误，例如不存在
	HttpUrlCheckFalse = 4103
)
