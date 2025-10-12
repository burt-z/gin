package dao

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	MysqlErrorDuplicateEmail = errors.New("邮箱已存在")
)

type UserDao interface {
	Insert(ctx context.Context, user User) error
	FindByEmail(ctx context.Context, email string) (User, error)
	FindById(ctx context.Context, id int64) (User, error)
	Save(ctx context.Context, user User) error
	FindByPhone(ctx context.Context, phone string) (User, error)
}

type GormUserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) UserDao {
	return &GormUserDao{
		db: db,
	}
}

func (d *GormUserDao) Insert(ctx context.Context, user User) error {
	user.CTime = time.Now().UnixMilli()
	user.UTime = time.Now().UnixMilli()
	err := d.db.WithContext(ctx).Create(&user).Error
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		fmt.Println("mysqlErr.Number", mysqlErr.Number == 1062)
		if mysqlErr.Number == 1062 {
			return MysqlErrorDuplicateEmail
		}
	}
	return err
}

func (d *GormUserDao) FindByEmail(ctx context.Context, email string) (User, error) {
	user := User{}
	err := d.db.WithContext(ctx).Table("users").Where("email = ?", email).Find(&user).Error
	return user, err
}

func (d *GormUserDao) FindById(ctx context.Context, id int64) (User, error) {
	user := User{}
	err := d.db.WithContext(ctx).Table("users").Where("id = ?", id).Find(&user).Error
	return user, err
}

// User 数据库上的概念,对应表内容
type User struct {
	Id       int64 `gorm:"primaryKey,autoIncrement"`
	Password string
	Email    sql.NullString `gorm:"unique"`
	Nickname string         `gorm:"nickname"`
	Birthday string         `gorm:"birthday"`
	AboutMe  string         `gorm:"about_me"`
	CTime    int64          //毫秒数
	UTime    int64          //毫秒数

	Phone sql.NullString `gorm:"unique"`
}

func (d *GormUserDao) Save(ctx context.Context, user User) error {
	return d.db.WithContext(ctx).Table("users").Save(&user).Error
}

func (d *GormUserDao) FindByPhone(ctx context.Context, phone string) (User, error) {
	user := User{}
	err := d.db.WithContext(ctx).Table("users").Where("phone = ?", phone).Find(&user).Error
	return user, err
}
