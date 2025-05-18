package biz

import (
	"context"
	"github.com/YangZhaoWeblog/UserService/internal/pkg"
)

// User 是用户领域模型
type User struct {
	ID       int64
	Username string
	Nickname string

	AuthType string // 通过什么方式注册的
	Phone    Phone

	AuthToken AuthToken
}

type AuthToken struct {
	TokenType    string
	ExpiresIn    int64
	AccessToken  string
	RefreshToken string
}

type Phone struct {
	Number           string
	VerificationCode string
}

// UserRepo 是用户仓库接口
type UserRepo interface {
	Save(context.Context, *User) (*User, error)
	Update(context.Context, *User) (*User, error)
	FindByID(context.Context, int64) (*User, error)
	FindByPhone(context.Context, string) (*User, error)
	FindByUsername(context.Context, string) (*User, error)
}

// UserUsecase 是用户用例
type UserUsecase struct {
	repo   UserRepo
	jwtCli *pkg.JwtClient
}

// NewUserUsecase 创建用户用例
func NewUserUsecase(repo UserRepo, jwtClient *pkg.JwtClient) *UserUsecase {
	return &UserUsecase{
		repo:   repo,
		jwtCli: jwtClient,
	}
}

const (
	AuthTypePhone  string = "phone"
	AuthTypeGoogle string = "google"
	AuthTypeNone   string = ""
)

// CreateUser 创建用户
func (uc *UserUsecase) CreateUser(ctx context.Context, u *User) (*User, error) {
	var err error
	var createdUser *User

	// 1. 创建用户
	switch u.AuthType {
	case AuthTypePhone:
		createdUser, err = uc.repo.Save(ctx, u)
	case AuthTypeGoogle:
		createdUser, err = uc.repo.Save(ctx, u)
		return nil, nil
	}
	if err != nil || createdUser == nil {
		return nil, err
	}

	// 2. 生成 JWT 令牌
	accessToken, refreshToken, err := uc.jwtCli.GenerateToken(u.ID, createdUser.Nickname)
	if err != nil {
		return nil, err
	}

	return &User{
		AuthToken: AuthToken{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}, nil
}

// GetUser 获取用户信息
func (uc *UserUsecase) GetUser(ctx context.Context, id int64) (*User, error) {
	return uc.repo.FindByID(ctx, id)
}

// GetUserByPhone 通过手机号获取用户
func (uc *UserUsecase) GetUserByPhone(ctx context.Context, phone string) (*User, error) {
	return uc.repo.FindByPhone(ctx, phone)
}

// UpdateUser 更新用户信息
func (uc *UserUsecase) UpdateUser(ctx context.Context, u *User) (*User, error) {
	return uc.repo.Update(ctx, u)
}
