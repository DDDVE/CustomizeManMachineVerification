// Code generated by goctl. DO NOT EDIT!
// Source: edit.proto

package server

import (
	"context"

	"rpc/internal/logic"
	"rpc/internal/svc"
	"rpc/types/edit"
)

type EditServer struct {
	svcCtx *svc.ServiceContext
	edit.UnimplementedEditServer
}

func NewEditServer(svcCtx *svc.ServiceContext) *EditServer {
	return &EditServer{
		svcCtx: svcCtx,
	}
}

func (s *EditServer) SubmitItem(ctx context.Context, in *edit.ItemRequest) (*edit.ItemResponse, error) {
	l := logic.NewSubmitItemLogic(ctx, s.svcCtx)
	return l.SubmitItem(in)
}

func (s *EditServer) UserInfo(ctx context.Context, in *edit.UserInfoRequest) (*edit.UserInfoResponse, error) {
	l := logic.NewUserInfoLogic(ctx, s.svcCtx)
	return l.UserInfo(in)
}