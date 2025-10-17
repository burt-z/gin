package domain

type ArticleStatus uint8

const (
	// ArticleStatusUnknown 为了避免零值之类的问题
	ArticleStatusUnknown ArticleStatus = iota
	ArticleStatusUnpublished
	ArticleStatusPublished
	ArticleStatusPrivate
)

func (s ArticleStatus) ToUint8() uint8 {
	return uint8(s)
}

func (s ArticleStatus) NonPublished() bool {
	return s != ArticleStatusPublished
}

func (s ArticleStatus) String() string {
	switch s {
	case ArticleStatusPrivate:
		return "private"
	case ArticleStatusUnpublished:
		return "unpublished"
	case ArticleStatusPublished:
		return "published"
	default:
		return "unknown"
	}
}

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
