package repository

import (
	"context"
	"go_project/gin/internal/domain"
	"go_project/gin/internal/repository/dao"
)

type UserRepository struct {
	dao *dao.UserDao
}

func NewUserRepository(dao *dao.UserDao) *UserRepository {
	return &UserRepository{
		dao: dao,
	}
}

// Create 数据层,没有 signup 的概念,所以是 create
func (r *UserRepository) Create(ctx context.Context, u domain.User) error {
	return r.dao.Insert(ctx, dao.User{Password: u.Password, Email: u.Email})
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{Email: u.Email, Password: u.Password, Id: u.Id}, nil
}

func (r *UserRepository) Update(ctx context.Context, user domain.User) error {
	dbUser, err := r.dao.FindById(ctx, user.Id)
	if err != nil {
		return err
	}

	for _, v := range user.Keys {
		switch v {
		case "nickname":
			dbUser.Nickname = user.NickName
		case "birthday":
			dbUser.Birthday = user.Birthday
		case "about_me":
			dbUser.AboutMe = user.AboutMe
		}
	}
	err = r.dao.Save(ctx, dbUser)
	if err != nil {
		return err
	}

	return nil
}

// FindById 查询
func (r *UserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
	daoUser, err := r.dao.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	domainUser := DaoUserToDomainUser(ctx, daoUser)
	return domainUser, nil
}

func DaoUserToDomainUser(ctx context.Context, u dao.User) domain.User {
	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		NickName: u.Nickname,
		Birthday: u.Birthday,
		AboutMe:  u.AboutMe,
	}

}
