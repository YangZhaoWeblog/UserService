package service

import (
	"context"
	v1 "user-svr/api/user/v1"
	"user-svr/internal/biz"
)

// UserService 是用户服务
type UserService struct {
	v1.UnimplementedUserServer

	uc *biz.UserUsecase
}

// NewUserService 创建用户服务
func NewUserService(uc *biz.UserUsecase) *UserService {
	return &UserService{uc: uc}
}

// Register 实现注册接口
func (s *UserService) Register(ctx context.Context, req *v1.RegisterRequest) (*v1.RegisterReply, error) {
	// TODO: 实现注册逻辑
	return &v1.RegisterReply{
		Success: true,
		Message: "注册成功",
	}, nil
}

// Login 实现登录接口
func (s *UserService) Login(ctx context.Context, req *v1.LoginRequest) (*v1.LoginReply, error) {
	// TODO: 实现登录逻辑
	return &v1.LoginReply{
		Success: true,
		Message: "登录成功",
	}, nil
}

// Info 实现获取用户信息接口
func (s *UserService) Info(ctx context.Context, req *v1.InfoRequest) (*v1.InfoReply, error) {
	// TODO: 实现获取用户信息逻辑
	return &v1.InfoReply{}, nil
}
