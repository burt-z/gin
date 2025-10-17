package article

import "context"

type ArticleReaderDao interface {
	Upsert(ctx context.Context, art Article) error
}

// PublishArticle 业务上分为线上表和创作表,同库不同表,可以使用下面的方式使用PublishArticle 表示线上表
type PublishArticle struct {
	Article
}
