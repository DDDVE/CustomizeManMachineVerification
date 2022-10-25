package logic

import (
	"context"
	"encoding/base64"
	"errors"
	"log"
	"pkg/conf"
	"regexp"
	"rpc/loginclient"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

const REGEXP_MOBILENUM = "^1[3-9][0-9]{9}$"

type SendMsgLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendMsgLogic {
	return &SendMsgLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendMsgLogic) SendMsg(req *types.SendMsgReq) (resp *types.SendMsgReply, err error) {

	mobile, err := base64.StdEncoding.DecodeString(req.MobileNum)
	log.Println("base64解码后的手机号为///////////", string(mobile))
	if err != nil {
		return nil, err
	}

	//简单参数校验
	if matched, err := regexp.MatchString(REGEXP_MOBILENUM, string(mobile)); err != nil || !matched || len(req.Token) != TOKEN_LENGHT {
		return nil, errors.New(conf.GlobalError[conf.ILLEGAL_REQUEST])
	}
	//缓存中token是否存在
	conn := l.svcCtx.RedisPool.Get()
	defer conn.Close()
	t, err := conn.Do("GET", KEY_PREFIX_TOKEN+req.Token)

	if value, ok := t.([]byte); !ok || string(value) != TOKEN_VALUE || err != nil {
		return nil, errors.New(conf.GlobalError[conf.ILLEGAL_REQUEST])
	}

	if _, err = l.svcCtx.LoginRpc.SendMsg(l.ctx, &loginclient.SendMsgRequest{MobileNum: string(mobile)}); err != nil {
		return nil, err
	}

	return
}
