package logic

import (
	"context"
	"rpc/model"
	"strconv"

	"rpc/internal/svc"
	"rpc/types/edit"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserInfoLogic) UserInfo(in *edit.UserInfoRequest) (*edit.UserInfoResponse, error) {
	em := model.NewDefaultEmployeeModel(l.svcCtx)
	e, err := em.SelectByMobile(in.MobileNum)
	if err != nil {
		return nil, err
	}
	return &edit.UserInfoResponse{
		EmployeeLevel:     strconv.Itoa(int(e.EmployeeLevel)),
		AuditScore:        strconv.Itoa(int(e.AuditScore)),
		ContributionScore: strconv.Itoa(int(e.ContributionScore)),
		RegistrationTime:  e.RegistrationTime.Format("2006-01-02"),
	}, nil
}
