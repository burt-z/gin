package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	redismock "jike_gin/internal/repository/cache/redismocks"
	"testing"
)

func TestRedisCodeCache_Set(t *testing.T) {
	testCase := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) redis.Cmdable
		wantErr error
	}{
		{
			name: "成功",
			//方法里面调用的 eva
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismock.NewMockCmdable(ctrl)
				//result, err := c.client.Eval(ctx, luaSetCode, []string{c.Key(ctx, biz, phone)}, code).Int()
				// eval 返回的是一个 Cmd,先调用了 set后又调用了 int
				res := redis.NewCmd(context.Background())
				// 因为后面调用了 int
				res.SetVal(int64(1))
				cmd.EXPECT().Eval(
					gomock.Any(),
					luaSetCode,
					[]string{"phone_code:login:1534444:v1"}, []any{"123456"}).Return(res)
				return cmd
			},
			wantErr: ErrCodeVerifyErr,
		},
	}
	for _, tt := range testCase {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			cache := NewCodeCache(tt.mock(ctrl))
			err := cache.Set(context.Background(), "login", "1534444", "123456")
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
