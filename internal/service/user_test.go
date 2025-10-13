package service

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"jike_gin/internal/domain"
	"jike_gin/internal/repository"
	repositorymock "jike_gin/internal/repository/mocks"
	"testing"
)

// 测试 login 方法,因为 Login 里面需要一个 repo(UserRepository),所以我们先要 mock 这个UserRepository
// 因为 这个UserRepository 调用了 FindByEmail ,我们用生成的调用这个方法
// 因为测试的是成功,所以 mock里面返回的就是 wantUser

func Test_userService_Login(t *testing.T) {
	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) repository.UserRepository
		ctx      context.Context
		email    string
		password string
		wantErr  error
		wantUser domain.User
	}{
		{
			name:     "success",
			email:    "123@qq.com",
			password: "123",
			wantErr:  nil,
			wantUser: domain.User{
				Email:    "123@qq.com",
				Password: "$2a$10$t65yaqhlbigyNFbQTc9Wu.YCWrSJ0.AIOoRpjXi8DmjB2C5BB9SI6",
			},
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repositorymock.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), gomock.Any()).Return(domain.User{
					Email:    "123@qq.com",
					Password: "$2a$10$t65yaqhlbigyNFbQTc9Wu.YCWrSJ0.AIOoRpjXi8DmjB2C5BB9SI6"}, nil)
				return repo
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			svc := NewUserService(tc.mock(ctrl))
			u, err := svc.Login(tc.ctx, tc.email, tc.password)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, u)
		})
	}
}

func TestGenerateFromPassword(t *testing.T) {
	res, err := bcrypt.GenerateFromPassword([]byte("123"), bcrypt.DefaultCost)
	if err != nil {
		t.Log("err:", err)
		return
	}
	t.Log(string(res))
}
