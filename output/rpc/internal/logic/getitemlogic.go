package logic

import (
	"context"

	"rpc/internal/svc"
	"rpc/model"
	"rpc/types/output"

	"github.com/zeromicro/go-zero/core/logx"
)

const ANY_TYPE = 0

type GetItemLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetItemLogic {
	return &GetItemLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetItemLogic) GetItem(in *output.OutputRequest) (*output.OutputResponse, error) {
	var auditedItem *model.AuditedItem
	var err error
	db := model.NewDefaultAuditedItemsModel(l.svcCtx)

	if in == nil {
		auditedItem, err = db.FindRandomOneByType(ANY_TYPE)
	} else {
		auditedItem, err = db.FindRandomOneByType(in.QuestionType)
	}
	if err != nil {
		return nil, err
	}

	return &output.OutputResponse{
		Question:      auditedItem.Question,
		Answer:        auditedItem.Answer,
		DisturbAnswer: auditedItem.DisturbAnswer,
	}, nil
}
