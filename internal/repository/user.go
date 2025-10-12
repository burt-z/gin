package repository

import (
	"context"
	"database/sql"
	"fmt"
	"jike_gin/internal/domain"
	"jike_gin/internal/repository/cache"
	"jike_gin/internal/repository/dao"
)

type UserRepository interface {
	Create(ctx context.Context, u domain.User) error
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	Update(ctx context.Context, user domain.User) error
	FindById(ctx context.Context, id int64) (domain.User, error)
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
}

type CacheUserRepository struct {
	dao   dao.UserDao
	cache cache.UserCache
}

func NewUserRepository(dao dao.UserDao, cache cache.UserCache) UserRepository {
	return &CacheUserRepository{
		dao:   dao,
		cache: cache,
	}
}

// Create 数据层,没有 signup 的概念,所以是 create
func (r *CacheUserRepository) Create(ctx context.Context, u domain.User) error {
	return r.dao.Insert(ctx, dao.User{Password: u.Password, Email: sql.NullString{Valid: u.Email != "", String: u.Email}})
}

func (r *CacheUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{Email: email, Password: u.Password, Id: u.Id}, nil
}

func (r *CacheUserRepository) Update(ctx context.Context, user domain.User) error {
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
func (r *CacheUserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
	uCache, err := r.cache.Get(ctx, id)
	if err == nil {
		return uCache, nil
	}

	daoUser, err := r.dao.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	domainUser := DaoUserToDomainUser(ctx, daoUser)
	err = r.cache.Set(ctx, domainUser)
	if err != nil {
		fmt.Println("Set===>", err.Error())
		return domainUser, nil
	}
	return domainUser, nil
}

func DaoUserToDomainUser(ctx context.Context, u dao.User) domain.User {
	return domain.User{
		Id:       u.Id,
		Email:    u.Email.String,
		NickName: u.Nickname,
		Birthday: u.Birthday,
		AboutMe:  u.AboutMe,
		Phone:    u.Phone.String,
	}
}

func (r *CacheUserRepository) domainUserToEntity(ctx context.Context, u domain.User) dao.User {
	return dao.User{
		Id:       u.Id,
		Email:    sql.NullString{String: u.Email, Valid: u.Email != ""},
		Birthday: u.Birthday,
		AboutMe:  u.AboutMe,
		Phone: sql.NullString{
			Valid:  u.Phone != "",
			String: u.Phone,
		},
	}
}

func (r *CacheUserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	daoUser, err := r.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	domainUser := DaoUserToDomainUser(ctx, daoUser)
	return domainUser, nil
}
