package logic

import (
	"context"
	"rpc/model"
	"strconv"

	"rpc/internal/svc"
	"rpc/types/edit"

	"github.com/zeromicro/go-zero/core/logx"
)

type SubmitItemLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSubmitItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SubmitItemLogic {
	return &SubmitItemLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SubmitItemLogic) SubmitItem(in *edit.ItemRequest) (*edit.ItemResponse, error) {
	aim := model.NewDefaultAuditedItemsModel(l.svcCtx)
	q, err := strconv.Atoi(in.QuestionType)
	if err != nil {
		return nil, err
	}
	var item = &model.AuditedItem{
		Producer:      in.MobileNum,
		QuestionType:  uint8(q),
		Question:      in.Question,
		Answer:        in.Answer,
		DisturbAnswer: in.DisturbAnswer,
	}
	if err = aim.Insert(item); err != nil {
		return nil, err
	}
	return &edit.ItemResponse{}, nil
}
