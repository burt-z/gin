package sms

import "context"

// Service 发送短信的抽象
type Service interface {
	Send(ctx context.Context, tpl string, args []string, numbers ...string) error
}

type NamedArg struct {
	Val  string
	Name string
}
