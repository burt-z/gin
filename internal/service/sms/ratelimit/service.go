package ratelimit

import (
	"context"
	"fmt"
	"jike_gin/internal/service/sms"
	"jike_gin/pkg/ratelimit"
)

type Service struct {
	svc     sms.Service
	limiter ratelimit.Limiter
}

func NewService(svc sms.Service, limiter ratelimit.Limiter) *Service {
	return &Service{
		svc:     svc,
		limiter: limiter,
	}
}

func (s *Service) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	limit, err := s.limiter.Limited(ctx, "sms:tencent")
	if err != nil {
		// 如果下游性能较差则进行限流
		// 不想影响业务,或者下游服务较稳定可以不限流
		return nil
	}
	if limit {
		return fmt.Errorf("%s", "限流错误")
	}
	err = s.svc.Send(ctx, tpl, args)
	return err
}
