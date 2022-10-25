package logic

import (
	"context"
	"log"

	"rpc/internal/svc"
	"rpc/model"
	"rpc/types/login"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddEmployeeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddEmployeeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddEmployeeLogic {
	return &AddEmployeeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddEmployeeLogic) AddEmployee(in *login.EmployeeRequest) (*login.EmployeeResponse, error) {
	log.Println("要注册的手机号：", in.MobileNum)
	em := model.NewDefaultEmployeeModel(l.svcCtx)
	if err := em.Insert(&model.Employee{MobileNum: in.MobileNum}); err != nil {
		log.Println("注册失败！！！")
		return nil, err
	}
	return &login.EmployeeResponse{}, nil
}
