package service

import (
	"context"
	"errors"
	"go_project/gin/internal/domain"
	"go_project/gin/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	resp *repository.UserRepository
}

func NewUserService(resp *repository.UserRepository) *UserService {
	return &UserService{
		resp: resp,
	}
}

// Signup 用户注册, 定义一个 user,保证数据是向下传递,不使用handler,数据也出现不一致了
func (u *UserService) Signup(ctx context.Context, user domain.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hash)
	return u.resp.Create(ctx, user)
}

// Login 登录
func (u *UserService) Login(ctx context.Context, user domain.User) (domain.User, error) {
	row, err := u.resp.FindByEmail(ctx, user.Email)
	if err != nil {
		return domain.User{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(row.Password), []byte(user.Password))
	if err != nil {
		return domain.User{}, errors.New("账号/邮箱或密码错误")
	}
	return row, nil
}

// ProfileEdit 编辑
func (u *UserService) ProfileEdit(ctx context.Context, user domain.User) error {
	err := u.resp.Update(ctx, user)
	return err
}

func (u *UserService) Profile(ctx context.Context, id int64) (domain.User, error) {
	domainUser, err := u.resp.FindById(ctx, id)
	return domainUser, err
}
