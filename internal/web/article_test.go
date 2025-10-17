package web

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"jike_gin/internal/domain"
	"jike_gin/internal/service"
	svcmock "jike_gin/internal/service/mocks"
	ijwt "jike_gin/internal/web/jwt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestArticleHandler_Publish(t *testing.T) {
	tests := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) service.ArticleService
		reqBody  string
		wantCode int
		wantRes  Result
	}{
		// TODO: Add test cases.
		{
			name: "success",
			mock: func(ctrl *gomock.Controller) service.ArticleService {
				userSvc := svcmock.NewMockArticleService(ctrl)
				userSvc.EXPECT().Publish(gomock.Any(), domain.Article{
					Title:   "我的标题",
					Content: "我的内容",
					Author: domain.Author{
						Id: 123,
					},
				}).Return(int64(1), nil)
				return userSvc
			},
			reqBody: `
{
	"title":"我的标题",
	"content": "我的内容"
}
`,
			wantCode: http.StatusOK,
			wantRes: Result{
				Data: float64(1), // json 转换需要转义
				Msg:  "OK",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			server := gin.Default()
			server.Use(func(context *gin.Context) {
				context.Set("claims", &ijwt.UserClaims{
					Uid: 123,
				})
			})

			handler := NewArticleHandler(tt.mock(ctrl))
			handler.RegisterRoutes(server)

			req, err := http.NewRequest("POST", "/articles/publish", bytes.NewBuffer([]byte(tt.reqBody)))
			req.Header.Set("Content-Type", "application/json")
			//err != nil 会结束执行
			require.NoError(t, err)

			resp := httptest.NewRecorder()

			server.ServeHTTP(resp, req)

			assert.Equal(t, tt.wantCode, resp.Code)
			if resp.Code != 200 {
				return
			}

			var webRes Result
			err = json.NewDecoder(resp.Body).Decode(&webRes)
			require.NoError(t, err)
			assert.Equal(t, tt.wantRes, webRes)
		})
	}
}
