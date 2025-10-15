package domain

type Article struct {
	Title   string
	Content string
	Author  Author
}

type Author struct {
	Name string
	Id   int64
}
