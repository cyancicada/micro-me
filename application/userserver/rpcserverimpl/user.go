package rpcserverimpl

import (
	"context"
	"errors"

	"micro-me/application/userserver/models"
	userpb "micro-me/application/userserver/protos"
)

type (
	UserRpcServer struct {
		useModel *models.MembersModel
	}
)

var (
	ErrNotFound = errors.New("用户不存在")
)

func NewUserRpcServer(useModel *models.MembersModel) *UserRpcServer {
	return &UserRpcServer{useModel: useModel}
}
func (s *UserRpcServer) FindByToken(ctx context.Context, req *userpb.FindByTokenRequest, rsp *userpb.UserResponse) error {
	member, err := s.useModel.FindByToken(req.Token)
	if err != nil {
		return ErrNotFound
	}
	rsp.Token = member.Token
	rsp.Id = member.Id
	rsp.Username = member.Username
	rsp.Password = member.Password
	return nil
}
func (s *UserRpcServer) FindById(ctx context.Context, req *userpb.FindByIdRequest, rsp *userpb.UserResponse) error {
	member, err := s.useModel.FindById(req.Id)
	if err != nil {
		return ErrNotFound
	}
	rsp.Token = member.Token
	rsp.Id = member.Id
	rsp.Password = member.Password
	return nil
}
