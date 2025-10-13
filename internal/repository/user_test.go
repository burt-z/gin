package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"jike_gin/internal/domain"
	"jike_gin/internal/repository/cache"
	cachemock "jike_gin/internal/repository/cache/mocks"
	"jike_gin/internal/repository/dao"
	daomock "jike_gin/internal/repository/dao/mocks"
	"testing"
)

// 测试FindById,因为用到了CacheUserRepository,所以需要先构造这个,CacheUserRepository里面用到了UserDao,UserCache
// 所以先生成这两个接口的 mock文件

func TestCacheUserRepository_FindById(t *testing.T) {
	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache)
		wantUser dao.User
		wantErr  error
		id       int64
	}{
		{
			name: "成功",
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {
				// 最长路径缓存未匹配
				c := cachemock.NewMockUserCache(ctrl)
				c.EXPECT().Get(gomock.Any(), int64(123)).Return(domain.User{}, errors.New("不存在"))
				d := daomock.NewMockUserDao(ctrl)
				d.EXPECT().FindById(gomock.Any(), int64(123)).Return(dao.User{
					Email: sql.NullString{String: "123@qq.com", Valid: true},
				}, nil)

				c.EXPECT().Set(gomock.Any(), domain.User{
					Email: "123@qq.com",
				}).Return(nil)

				return d, c
			},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			d, c := tt.mock(ctrl)
			repo := NewUserRepository(d, c)
			user, err := repo.FindById(context.Background(), tt.id)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.wantUser, user)
		})
	}
}
