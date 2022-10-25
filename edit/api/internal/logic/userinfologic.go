package logic

import (
	"context"
	"encoding/base64"
	"errors"
	"log"
	"pkg/conf"
	"regexp"
	"rpc/types/edit"
	"strconv"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserInfoLogic) UserInfo(req *types.UserInfoReq) (resp *types.UserInfoReply, err error) {

	log.Printf("UserInfoReq///////////////////////////////////%+v", req)
	mobile, err := base64.StdEncoding.DecodeString(req.MobileNum)
	if err != nil {
		return nil, err
	}

	if matched, err := regexp.MatchString(REGEXP_MOBILENUM, string(mobile)); err != nil || !matched {
		return nil, errors.New(conf.GlobalError[conf.ILLEGAL_REQUEST])
	}
	uir, err := l.svcCtx.EditRpc.UserInfo(l.ctx, &edit.UserInfoRequest{MobileNum: string(mobile)})
	if err != nil {
		return nil, err
	}
	level, err := strconv.Atoi(uir.EmployeeLevel)
	if err != nil {
		return nil, err
	}
	as, err := strconv.Atoi(uir.AuditScore)
	if err != nil {
		return nil, err
	}
	cs, err := strconv.Atoi(uir.ContributionScore)
	if err != nil {
		return nil, err
	}
	return &types.UserInfoReply{
		EmployeeLevel:     level,
		AuditScore:        as,
		ContributionScore: cs,
		RegistrationTime:  uir.RegistrationTime,
	}, nil
}
