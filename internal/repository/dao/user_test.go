package dao

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func TestGormUserDao_Insert(t *testing.T) {

	testCase := []struct {
		name    string
		ctx     context.Context
		user    User
		wantErr error
		mock    func(t *testing.T) *sql.DB
	}{
		{
			name: "success",
			ctx:  context.Background(),
			user: User{Password: "wer3fer3", Email: sql.NullString{String: "124@qq,com", Valid: true}},
			mock: func(t *testing.T) *sql.DB {
				mockDb, mock, err := sqlmock.New()
				sqlRes := sqlmock.NewResult(1, 1)
				mock.ExpectExec("INSERT INTO `users` .*").WillReturnResult(sqlRes)
				require.NoError(t, err)
				return mockDb
			},
		},
	}
	for _, tt := range testCase {
		t.Run(tt.name, func(t *testing.T) {

			db, err := gorm.Open(mysql.New(mysql.Config{
				Conn: tt.mock(t),
				//禁用选择版本
				SkipInitializeWithVersion: true,
			}), &gorm.Config{
				// gorm 默认开启事务,禁用事务
				SkipDefaultTransaction: true,
			})
			require.NoError(t, err)
			userDao := NewUserDao(db)
			err = userDao.Insert(tt.ctx, tt.user)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
