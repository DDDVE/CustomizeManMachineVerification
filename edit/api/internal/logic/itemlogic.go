package logic

import (
	"context"
	"errors"
	"pkg/conf"
	"regexp"
	"rpc/types/edit"
	"strconv"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

const REGEXP_MOBILENUM = "^1[3-9][0-9]{9}$"

type ItemLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ItemLogic {
	return &ItemLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ItemLogic) Item(req *types.ItemReq) (resp *types.ItemReply, err error) {
	//简单参数校验
	if matched, err := regexp.MatchString(REGEXP_MOBILENUM, req.MobileNum); err != nil || !matched {
		return nil, errors.New(conf.GlobalError[conf.ILLEGAL_REQUEST])
	}
	if _, err := l.svcCtx.EditRpc.SubmitItem(l.ctx, &edit.ItemRequest{
		MobileNum:     req.MobileNum,
		QuestionType:  strconv.Itoa(req.QuestionType),
		Question:      req.Question,
		Answer:        req.Answer,
		DisturbAnswer: req.DisturbAnswer,
	}); err != nil {
		return nil, err
	}
	return &types.ItemReply{}, nil
}
