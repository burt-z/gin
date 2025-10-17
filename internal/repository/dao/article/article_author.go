package article

import (
	"context"
)

type ArticleAuthorDao interface {
	Insert(ctx context.Context, art Article) (int64, error)
	UpdateById(ctx context.Context, art Article) error
}
