package failover

import (
	"context"
	"jike_gin/internal/service/sms"
	"sync/atomic"
)

type TimeoutFailoverSMSService struct {
	svcs      []sms.Service // 多个服务商
	idx       int32         // 取第几个服务商
	cnt       int32         //联系超时的个数
	threshold int32         // 超过几个就需要切换
}

func NewTimeoutFailoverSMSService(svcs []sms.Service, idx int32, threshold int32) *TimeoutFailoverSMSService {
	return &TimeoutFailoverSMSService{
		svcs:      svcs,
		cnt:       0,
		threshold: threshold,
		idx:       idx,
	}
}

func (svc *TimeoutFailoverSMSService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	idx := atomic.LoadInt32(&svc.idx)
	cnt := atomic.LoadInt32(&svc.cnt)
	if cnt >= svc.threshold {
		// 需要切换
		newIdx := (idx + 1) % int32(len(svc.svcs))
		if atomic.CompareAndSwapInt32(&svc.idx, idx, newIdx) {
			atomic.StoreInt32(&svc.cnt, 0)
		}
		idx = atomic.LoadInt32(&svc.cnt)
	}
	s := svc.svcs[idx]

	err := s.Send(ctx, tpl, args)
	switch err {
	case context.DeadlineExceeded:
		atomic.AddInt32(&svc.cnt, 1)
		return err
	case nil:
		// 联系状态被打断
		atomic.StoreInt32(&svc.cnt, 0)
		return nil
	default:
		return err
	}
}
