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

const (
	KEY_PREFIX_MSG_CODE = "cache:login:msgCode:mobileNum:"
	REGEXP_MSG_CODE     = "^[0-9]{6}$"
)

type EmployeeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEmployeeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EmployeeLogic {
	return &EmployeeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EmployeeLogic) Employee(req *types.EmployeeReq) (resp *types.EmployeeReply, err error) {
	log.Printf("%v//////%v//////%v", req.MobileNum, req.MsgCode, req.Token)
	mobile, err := base64.StdEncoding.DecodeString(req.MobileNum)
	if err != nil {
		return nil, err
	}

	//简单参数校验
	if matched, err := regexp.MatchString(REGEXP_MOBILENUM, string(mobile)); err != nil || !matched || len(req.Token) != TOKEN_LENGHT {
		return nil, errors.New(conf.GlobalError[conf.ILLEGAL_REQUEST])
	}
	if matched, err := regexp.MatchString(REGEXP_MSG_CODE, req.MsgCode); err != nil || !matched {
		return nil, errors.New(conf.GlobalError[conf.ILLEGAL_REQUEST])
	}

	//缓存中token是否存在
	conn := l.svcCtx.RedisPool.Get()
	defer conn.Close()
	t, err := conn.Do("GET", KEY_PREFIX_TOKEN+req.Token)
	if value, ok := t.([]byte); !ok || string(value) != TOKEN_VALUE || err != nil {
		return nil, errors.New(conf.GlobalError[conf.ILLEGAL_REQUEST])
	}

	//与缓存中的msgCode对比
	msgCode, err := conn.Do("GET", KEY_PREFIX_MSG_CODE+string(mobile))
	log.Println("缓存中的msgCode为//////", msgCode)
	if value, ok := msgCode.([]byte); !ok || string(value) != req.MsgCode || err != nil {
		return nil, errors.New(conf.GlobalError[conf.MSG_CODE_ERROR])
	}
	if _, err = l.svcCtx.LoginRpc.AddEmployee(l.ctx, &loginclient.EmployeeRequest{MobileNum: string(mobile)}); err != nil {
		return nil, err
	}
	return &types.EmployeeReply{Token: ""}, nil
}
