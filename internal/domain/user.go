package domain

// User 业务意义上的概念
type User struct {
	Id       int64
	Email    string
	Password string
	NickName string
	Birthday string
	AboutMe  string
	Keys     []string
}
