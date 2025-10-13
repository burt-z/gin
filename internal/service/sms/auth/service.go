package auth

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"jike_gin/internal/service/sms"
)

// 装饰器模式,

type SmsService struct {
	svcs sms.Service
	key  string
}

func NewSmsService(svcs sms.Service, key string) *SmsService {
	return &SmsService{
		svcs: svcs,
		key:  key,
	}
}

// 给请求提供一个 token,通过 token来控制业务方是否有权限,使用 jwttoken,biz 是业务方请求的 token
func (s *SmsService) Send(ctx context.Context, biz string, args []string, numbers ...string) error {

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(biz, claims, func(token *jwt.Token) (interface{}, error) {
		return s.key, nil
	})
	if err != nil {
		fmt.Println("装饰器 send===>", err)
		return err
	}
	if !token.Valid {
		return err
	}

	return s.svcs.Send(ctx, claims.Tpl, args, numbers...)
}

type Claims struct {
	jwt.RegisteredClaims
	Tpl string //模版 id
}
