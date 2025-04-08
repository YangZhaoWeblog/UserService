package data

import (
	"context"

	"github.com/YangZhaoWeblog/UserService/internal/biz"
)

// UserRepo 实现 biz.UserRepo 接口
type userRepo struct {
	data *Data
}

// NewUserRepo 创建用户仓库实例
func NewUserRepo(data *Data) biz.UserRepo {
	return &userRepo{
		data: data,
	}
}

// Save 保存用户
func (r *userRepo) Save(ctx context.Context, u *biz.User) (*biz.User, error) {
	// TODO: 实现用户保存逻辑
	return u, nil
}

// Update 更新用户
func (r *userRepo) Update(ctx context.Context, u *biz.User) (*biz.User, error) {
	// TODO: 实现用户更新逻辑
	return u, nil
}

// FindByID 通过ID查找用户
func (r *userRepo) FindByID(ctx context.Context, id int64) (*biz.User, error) {
	// TODO: 实现通过ID查找用户逻辑
	return &biz.User{
		ID:       id,
		Username: "user_" + string(id),
		Nickname: "用户" + string(id),
	}, nil
}

// FindByPhone 通过手机号查找用户
func (r *userRepo) FindByPhone(ctx context.Context, phone string) (*biz.User, error) {
	// TODO: 实现通过手机号查找用户逻辑
	return &biz.User{
		Phone:    phone,
		Username: "user_" + phone,
		Nickname: "用户" + phone,
	}, nil
}

// FindByUsername 通过用户名查找用户
func (r *userRepo) FindByUsername(ctx context.Context, username string) (*biz.User, error) {
	// TODO: 实现通过用户名查找用户逻辑
	return &biz.User{
		Username: username,
		Nickname: "用户" + username,
	}, nil
}
