package service

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"jike_gin/internal/domain"
	"jike_gin/internal/repository"
)

type UserService interface {
	Signup(ctx context.Context, user domain.User) error
	Login(ctx context.Context, email, password string) (domain.User, error)
	ProfileEdit(ctx context.Context, user domain.User) error
	Profile(ctx context.Context, id int64) (domain.User, error)
	FindOrCreate(ctx context.Context, phone string) (domain.User, error)
	FindOrCreateByWechat(ctx context.Context, wechatInfo domain.WechatInfo) (domain.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(resp repository.UserRepository) UserService {
	return &userService{
		repo: resp,
	}
}

// Signup 用户注册, 定义一个 user,保证数据是向下传递,不使用handler,数据也出现不一致了
func (u *userService) Signup(ctx context.Context, user domain.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hash)
	return u.repo.Create(ctx, user)
}

// Login 登录
func (u *userService) Login(ctx context.Context, email, password string) (domain.User, error) {
	row, err := u.repo.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(row.Password), []byte(password))
	if err != nil {
		return domain.User{}, errors.New("账号/邮箱或密码错误")
	}
	return row, nil
}

// ProfileEdit 编辑
func (u *userService) ProfileEdit(ctx context.Context, user domain.User) error {
	err := u.repo.Update(ctx, user)
	return err
}

func (u *userService) Profile(ctx context.Context, id int64) (domain.User, error) {
	domainUser, err := u.repo.FindById(ctx, id)
	return domainUser, err
}

func (u *userService) FindOrCreate(ctx context.Context, phone string) (domain.User, error) {
	domainUser, err := u.repo.FindByPhone(ctx, phone)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return domain.User{}, err
	}
	if domainUser.Id > 0 {
		return domainUser, nil
	}
	newU := domain.User{Phone: phone}
	err = u.repo.Create(ctx, newU)
	if err != nil {
		return domain.User{}, err
	}
	return u.repo.FindByPhone(ctx, phone)
}

func (u *userService) FindOrCreateByWechat(ctx context.Context, wechatInfo domain.WechatInfo) (domain.User, error) {
	domainUser, err := u.repo.FindByWeChat(ctx, wechatInfo.OpenID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return domain.User{}, err
	}
	if domainUser.Id > 0 {
		return domainUser, nil
	}
	newU := domain.User{WechatInfo: wechatInfo}
	err = u.repo.Create(ctx, newU)
	if err != nil {
		return domain.User{}, err
	}
	return u.repo.FindByWeChat(ctx, wechatInfo.OpenID)
}
