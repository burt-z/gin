package tencent

import (
	"context"
	"fmt"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

type Service struct {
	appid    *string
	signName *string
	client   *sms.Client
}

func (s *Service) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	req := sms.NewSendSmsRequest()
	req.SmsSdkAppId = s.appid
	req.SignName = s.signName
	req.TemplateId = &tpl
	req.PhoneNumberSet = s.toStringPtrSlice(numbers)
	req.TemplateParamSet = s.toStringPtrSlice(args)
	res, err := s.client.SendSms(req)
	if err != nil {
		fmt.Println("Send====>", err.Error())
		return err
	}
	for _, status := range res.Response.SendStatusSet {
		if status == nil || *(status.Code) != "OK" {
			return fmt.Errorf("failed to send sms: %s %s", *(status.Code), *(status.Message))
		}
	}
	return nil
}

func (s *Service) toStringPtrSlice(src []string) []*string {
	res := make([]*string, 0)
	for _, n := range src {
		res = append(res, &n)
	}
	return res
}

func NewService(ctx context.Context, appid string, signName string, signClient *sms.Client) *Service {
	return &Service{
		appid:    &appid,
		signName: &signName,
		client:   signClient,
	}
}
