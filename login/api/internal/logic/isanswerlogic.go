package logic

import (
	"context"
	"errors"

	"api/internal/svc"
	"api/internal/types"
	"pkg/conf"

	"log"

	"github.com/zeromicro/go-zero/core/logx"
)

const (
	TOKEN_VALUE         = "OK"
	TOKEN_LENGHT        = 32
	EXPIRE_LENGHT_TOKEN = 10 * 60
)

type IsAnswerLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewIsAnswerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IsAnswerLogic {
	return &IsAnswerLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *IsAnswerLogic) IsAnswer(req *types.IsAnswerReq) (resp *types.IsAnswerReply, err error) {

	//简单参数校验
	if req == nil || req.Answer == "" || len(req.Token) != TOKEN_LENGHT {
		return nil, errors.New(conf.GlobalError[conf.ILLEGAL_REQUEST])
	}

	//和缓存中的answer对比
	conn := l.svcCtx.RedisPool.Get()
	defer conn.Close()
	answer, err := conn.Do("GET", KEY_PREFIX_TOKEN+req.Token)

	if err != nil {
		return nil, errors.New(conf.GlobalError[conf.ILLEGAL_REQUEST])
	}
	if value, ok := answer.([]byte); !ok || string(value) != req.Answer {
		return nil, errors.New(conf.GlobalError[conf.ANSWER_ERROR])
	}

	//续期缓存并更改token在缓存中的值
	if _, err = conn.Do("SETEX", KEY_PREFIX_TOKEN+req.Token, EXPIRE_LENGHT_TOKEN, TOKEN_VALUE); err != nil {
		return nil, err
	}
	log.Printf("======================================================================\n\n")
	// resp.Result = true //太特么恶心人了
	log.Printf("======================================================================\n\n")
	return &types.IsAnswerReply{Result: true}, nil
}
