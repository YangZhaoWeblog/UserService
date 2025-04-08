package biz

import (
	"context"
)

// User 是用户模型
type User struct {
	ID        int64
	Username  string
	Nickname  string
	AvatarURL string
	Phone     string
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
	repo UserRepo
}

// NewUserUsecase 创建用户用例
func NewUserUsecase(repo UserRepo) *UserUsecase {
	return &UserUsecase{repo: repo}
}

// CreateUser 创建用户
func (uc *UserUsecase) CreateUser(ctx context.Context, u *User) (*User, error) {
	return uc.repo.Save(ctx, u)
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
