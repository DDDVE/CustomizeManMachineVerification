package utils

// 全局错误码
const (
	SUCCESS              = 0
	RUFUSE_LEGAL_REQUEST = 1
	UNKNOWN_ERROR        = -1

	FORWARD_REQUEST_TO_API = 4001

	BLACK_USER            = 4100
	TOKEN_CHECK_ERROR     = 4101
	REQUEST_METHOD_ERROR  = 4102
	REQUEST_PATH_ERROR    = 4103
	REQUEST_PARAM_ERROR   = 4104
	SIGNATURE_CHECK_ERROR = 4105
	REQUEST_TIMEOUT       = 4106
	REPEATED_REQUEST      = 4107
	ANSWER_ERROR          = 4108
	ILLEGAL_REQUEST       = 4109

	SEND_MSG_ERROR = 4201
	MSG_CODE_ERROR = 4202

	NO_DATA_WAS_QUERIED = 4301
	INSERT_DATA_ERROR   = 4302
)

var GlobalError = map[int]string{
	0: "成功",

	//代表拒绝合法请求
	1:  "服务器繁忙,请稍后重试",
	-1: "未定义错误/未知错误",

	//40xx：代表转发请求出错，请求超时、网路错误等
	4001: "转发给api网关请求出错",

	//41xx:代表请求校验错误
	4100: "抱歉，您频繁访问被认定为恶意用户，已加入临时黑名单，一天后解封",
	4101: "token校验失败,请重新登录",
	4102: "请求方法错误",
	4103: "请求路径不存在",
	4104: "请求参数（字段）错误",
	4105: "签名校验错误",
	4106: "请求超时",
	4107: "此请求已处理，勿重复发送",
	4108: "人机验证答案错误",
	4109: "非法请求",

	// 42xx:代表请求通过校验，业务处理时出错
	4201: "短信发送失败",
	4202: "验证码错误",

	// 43xx:代表查询数据库错误/未查到数据
	4301: "未查询到相关数据",
	4302: "新增数据失败",

	// 44xx:其他
}
