package logic

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"time"

	"pkg/conf"
	"rpc/internal/svc"
	"rpc/types/login"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"github.com/zeromicro/go-zero/core/logx"
)

const (
	KEY_PREFIX_MSG_CODE    = "cache:login:msgCode:mobileNum:"
	EXPIRE_LENGHT_MSG_CODE = 600
)

type SendMsgLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendMsgLogic {
	return &SendMsgLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SendMsgLogic) SendMsg(in *login.SendMsgRequest) (*login.SendMsgResponse, error) {
	log.Printf("in/////////%+v", in)
	// 实例化一个认证对象，入参需要传入腾讯云账户secretId，secretKey,此处还需注意密钥对的保密
	// 密钥可前往https://console.cloud.tencent.com/cam/capi网站进行获取
	credential := common.NewCredential(
		"AKID64whQ7e0CC1fBlNq8E8jeYMzVvGVb7ub",
		"obHoIfvAMWXiopjXGWNnLkyzqeO5cEQZ",
	)
	// 实例化一个client选项，可选的，没有特殊需求可以跳过
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "sms.tencentcloudapi.com"
	// 实例化要请求产品的client对象,clientProfile是可选的
	client, _ := sms.NewClient(credential, "ap-nanjing", cpf)

	// 实例化一个请求对象,每个接口都会对应一个request对象
	request := sms.NewSendSmsRequest()

	mobileNums := []string{}
	msgCodes := []string{}
	msgCode := generateMsgCode()
	log.Println("生成的验证码是:", msgCode)
	request.PhoneNumberSet = common.StringPtrs(append(mobileNums, in.MobileNum))
	request.SmsSdkAppId = common.StringPtr("1400736497")
	request.SignName = common.StringPtr("自定义人机验证码制作网站")
	request.TemplateId = common.StringPtr("1572159")
	request.TemplateParamSet = common.StringPtrs(append(msgCodes, msgCode))

	// 返回的resp是一个SendSmsResponse的实例，与请求对象对应
	response, err := client.SendSms(request)
	if err != nil || *response.Response.SendStatusSet[0].Code != "Ok" {
		log.Printf("response:///%+v", response)
		return nil, errors.New(conf.GlobalError[conf.SEND_MSG_ERROR])
	}
	//缓存msgcode
	conn := l.svcCtx.RedisPool.Get()
	defer conn.Close()
	if _, err = conn.Do("SETEX", KEY_PREFIX_MSG_CODE+in.MobileNum, EXPIRE_LENGHT_MSG_CODE, msgCode); err != nil {
		return nil, err
	}

	log.Println("发送短信验证码成功")
	return &login.SendMsgResponse{}, nil
}

func generateMsgCode() string {
	msgCode := []byte{}
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 6; i++ {
		msgCode = append(msgCode, byte(rand.Intn(10)+48))
	}
	return string(msgCode)
}
