package service

import (
	"context"
	"fmt"
	"jike_gin/internal/repository"
	"jike_gin/internal/service/sms"
	"math/rand"
)

const codeTplId = "1877556"

type CodeService interface {
	Send(ctx context.Context, biz string, phone string) error
	Verify(ctx context.Context, biz string, phone string, inputCode string) error
}

type codeService struct {
	repo   repository.CodeRepository
	smsSvc sms.Service
}

var (
	ErrCodeVerifyTooManyTimes = repository.ErrCodeVerifyTooManyTimes
	ErrCodeSendTooMany        = repository.ErrCodeSendTooMany
)

func NewCodeService(repo repository.CodeRepository, smsSvc sms.Service) CodeService {
	return &codeService{
		repo:   repo,
		smsSvc: smsSvc,
	}
}

func (c *codeService) Send(ctx context.Context, biz string, phone string) error {
	code := c.generateCode(ctx)
	err := c.repo.Store(ctx, biz, phone, code)
	if err != nil {
		fmt.Println("Send store err:", err.Error())
		return err
	}
	err = c.smsSvc.Send(ctx, codeTplId, []string{biz}, phone)
	if err != nil {
		// redis 有,但是发送失败
		//尝试重试
		fmt.Println("Send to sms failed, err:===>", err.Error())
		return err
	}
	fmt.Println("Send to sms success:===>", code)
	return nil
}

func (c *codeService) Verify(ctx context.Context, biz string, phone string, inputCode string) error {
	return c.repo.Verify(ctx, biz, phone, inputCode)
}

func (c *codeService) generateCode(ctx context.Context) string {
	num := rand.Intn(1000000)
	return fmt.Sprintf("%6d", num)

}
