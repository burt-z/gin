package domain

type Article struct {
	Title   string
	Content string
	Author  Author
	Id      int64
}

type Author struct {
	Name string
	Id   int64
}
