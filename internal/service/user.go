package service

import (
	"context"
	"fmt"

	"github.com/YangZhaoWeblog/GoldenTakin/takin_log"
	v1 "github.com/YangZhaoWeblog/UserService/api/user/v1"
	"github.com/YangZhaoWeblog/UserService/internal/biz"
	"github.com/go-kratos/kratos/v2/errors"
)

// UserService 是用户服务
type UserService struct {
	v1.UnimplementedUserServer
	uc        *biz.UserUsecase
	logHelper *takin_log.TakinLogger
}

// NewUserService 创建用户服务
func NewUserService(uc *biz.UserUsecase, log *takin_log.TakinLogger) *UserService {
	return &UserService{uc: uc,
		logHelper: log,
	}
}

func getAuthTypeString(user *biz.User, req *v1.RegisterRequest) *biz.User {
	user.AuthType = biz.AuthTypeNone
	if req.GetPhone() != nil {
		user.AuthType = biz.AuthTypePhone
		user.Phone = biz.Phone{
			Number:           req.GetPhone().GetPhoneNumber(),
			VerificationCode: req.GetPhone().GetVerificationCode(),
		}
	} else if req.GetGoogle() != nil {
		user.AuthType = biz.AuthTypeGoogle
		user.Phone = biz.Phone{
			Number:           req.GetPhone().GetPhoneNumber(),
			VerificationCode: req.GetPhone().GetVerificationCode(),
		}
	}
	return user
}

// Register 实现注册接口
func (s *UserService) Register(ctx context.Context, req *v1.RegisterRequest) (*v1.RegisterReply, error) {
	user := biz.User{
		Nickname: req.GetNickname(),
	}
	getAuthTypeString(&user, req)

	// 1. 注册
	createdUser, err := s.uc.CreateUser(ctx, &user)
	if err != nil {
		return nil, nil
	}

	return &v1.RegisterReply{
		UserInfo: &v1.UserInfo{
			UserId:   fmt.Sprintf("%d", user.ID),
			Nickname: createdUser.Nickname,
		},
		AuthToken: &v1.AuthToken{
			AccessToken:  createdUser.AuthToken.AccessToken,
			RefreshToken: createdUser.AuthToken.RefreshToken,
			ExpiresIn:    createdUser.AuthToken.ExpiresIn,
			TokenType:    "Bearer",
		},
	}, nil
}

// Login 实现登录接口
func (s *UserService) Login(ctx context.Context, req *v1.LoginRequest) (*v1.LoginReply, error) {
	// TODO: 实现登录逻辑

	//s.logHelper.Info("测试:", "Login", "是否", "成功")
	return nil, errors.New(521, "LOGIN_FAILED", "登录失败")

	// 不可达代码
	//return &v1.LoginReply{
	//	Success: true,
	//	Message: "登录成功",
	//}, nil
}

// Info 实现获取用户信息接口
func (s *UserService) Info(ctx context.Context, req *v1.InfoRequest) (*v1.InfoReply, error) {
	// TODO: 实现获取用户信息逻辑
	return &v1.InfoReply{}, nil
}
