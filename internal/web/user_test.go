package web

import (
	"bytes"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"jike_gin/internal/domain"
	"jike_gin/internal/service"
	svcmock "jike_gin/internal/service/mocks"
	"net/http"
	"net/http/httptest"
	"testing"
)

//测试方法里面调用svcmock.New.. 这个svcmock是命令行里面的-package后面的值
//mockgen -source=gin/internal/service/code.go -package=svcmock -destination=gin/internal/service/mocks/code.mock.go

func TestUserHandler_SingUp(t *testing.T) {
	tests := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) service.UserService
		reqBody  string
		wantCode int
		wantBody string
	}{
		// TODO: Add test cases.
		{
			name: "success",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmock.NewMockUserService(ctrl)
				userSvc.EXPECT().Signup(gomock.Any(), domain.User{
					Email:    "123@qq.com",
					Password: "hello#world123",
				}).Return(nil)
				return userSvc
			},
			reqBody: `{
				"email": "123@qq.com",
				"password": "hello#world123",
				"confirmPassword": "hello#world123"
				}`,
			wantCode: http.StatusOK,
			wantBody: "注册成功",
		},
		{
			name: "邮箱格式不对",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmock.NewMockUserService(ctrl)
				return usersvc
			},

			reqBody: `
				{
					"email": "123@q",
					"password": "hello#world123",
					"confirmPassword": "hello#world123"
				}
			`,
			wantCode: http.StatusOK,
			wantBody: "你的邮箱格式不对",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			server := gin.Default()
			handler := NewUserHandler(tt.mock(ctrl), nil)
			handler.RegisterRoutes(server)

			req, err := http.NewRequest("POST", "/users/signup", bytes.NewBuffer([]byte(tt.reqBody)))
			req.Header.Set("Content-Type", "application/json")
			//err != nil 会结束执行
			require.NoError(t, err)

			resp := httptest.NewRecorder()

			//这是 http 请求进 gin的入口
			//当你这样调用的时候,gin会处理
			//响应写回到resp 里面
			server.ServeHTTP(resp, req)

			assert.Equal(t, tt.wantCode, resp.Code)
			assert.Equal(t, tt.wantBody, resp.Body.String())
		})
	}
}

func TestMock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRvc := svcmock.NewMockUserService(ctrl)

	userRvc.EXPECT().Signup(gomock.Any(), gomock.Any()).Return(errors.New("mock error"))
	err := userRvc.Signup(context.Background(), domain.User{Email: "1433507825@qq.com"})
	t.Log(err)
}
