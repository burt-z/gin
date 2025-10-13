package failover

import (
	"context"
	"fmt"
	"jike_gin/internal/service/sms"
)

type FailoverSMSService struct {
	svcs []sms.Service // 多个服务商
}

func NewFailoverSMSService(svcs []sms.Service) *FailoverSMSService {
	return &FailoverSMSService{
		svcs: svcs,
	}
}

func (svc *FailoverSMSService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	for _, s := range svc.svcs {
		err := s.Send(ctx, tpl, args)
		if err == nil {
			return nil
		}
	}
	return fmt.Errorf("全部失败")
}
