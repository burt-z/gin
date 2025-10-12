package tencent

import (
	"context"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"testing"
)

func TestService_Send(t *testing.T) {

	type fields struct {
		appid    *string
		signName *string
		client   *sms.Client
	}

	tests := []struct {
		ctx     context.Context
		tpl     string
		args    []string
		numbers []string
		name    string
	}{
		{
			ctx:     context.Background(),
			tpl:     "",
			name:    "发送验证码",
			args:    []string{"test"},
			numbers: []string{"143"},
		},
	}

	secretId := ""
	secretKey := ""
	c, err := sms.NewClient(common.NewCredential(secretId, secretKey), "", profile.NewClientProfile())
	if err != nil {
		t.Fatal(err)
		return
	}
	s := NewService(context.Background(), "", "", c)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := s.Send(tt.ctx, tt.tpl, tt.args, tt.numbers...); err != nil {
				t.Errorf("Send() error = %v", err)
			}
		})
	}
}
