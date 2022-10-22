package utils

// 统一记录常量，分为需要修改的和基本不变的常量

/**
0：查询成功
1：拒绝合法请求，返回服务器繁忙，稍后重试
41xx:代表请求校验错误
42xx:代表请求通过校验，业务处理时出错
43xx:代表查询数据库错误/未查到数据
44xx:其他
*/
// const (
// 	HttpSucceed = 0
// 	HttpRefuse  = 1
// 	// token校验错误
// 	HttpTokenCheckFalse = 4101
// 	// 请求方法错误
// 	HttpMethodCheckFalse = 4102
// 	// 请求路径错误，例如不存在
// 	HttpUrlCheckFalse = 4103
// 	// 请求参数（字段）错误
// 	HttpParamCheckFalse = 4104

// 	// 请求通过校验，api注册业务出错
// 	HttpApiRegistFalse = 4201
// )

// /**
// api网关相关常量
// */
// const (
// 	// 该类型在切片中的位置
// 	PositionOfApiLogin    = 0
// 	PositionOfApiEdit     = 1
// 	PositionOfApiAudit    = 2
// 	PositionOfApiFeedback = 3

// 	// 该类型的名字
// 	TypeOfApiLogin    = "login"
// 	TypeOfApiEdit     = "edit"
// 	TypeOfApiAudit    = "audit"
// 	TypeOfApiFeedback = "feedback"
// )

// var (
// 	// api网关类型切片
// 	ApiGateSlice = []string{TypeOfApiLogin, TypeOfApiEdit, TypeOfApiAudit, TypeOfApiFeedback}
// )

// // 统一连接过期时间
// const LinkTimeOut = 3
