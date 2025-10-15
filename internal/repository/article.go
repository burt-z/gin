package repository

import (
	"context"
	"jike_gin/internal/domain"
	"jike_gin/internal/repository/dao"
)

type ArticleRepository interface {
	Create(ctx context.Context, article domain.Article) (int64, error)
}

func NewArticleRepository(dao dao.ArticleDAO) ArticleRepository {
	return &articleService{
		dao: dao,
	}
}

type articleService struct {
	dao dao.ArticleDAO
}

func (a *articleService) Create(ctx context.Context, article domain.Article) (int64, error) {
	return a.dao.Insert(ctx, dao.Article{Title: article.Title, Content: article.Content})
}
