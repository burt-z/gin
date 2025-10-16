package service

import (
	"context"
	"jike_gin/internal/domain"
	"jike_gin/internal/repository"
)

type ArticleService interface {
	Save(ctx context.Context, article domain.Article) (id int64, err error)
	Update(ctx context.Context, article domain.Article) error
}

type articleService struct {
	repo repository.ArticleRepository
}

func NewArticleService(repo repository.ArticleRepository) ArticleService {
	return &articleService{
		repo: repo,
	}
}

func (a *articleService) Save(ctx context.Context, art domain.Article) (int64, error) {
	return a.repo.Create(ctx, art)
}

func (a *articleService) Update(ctx context.Context, art domain.Article) error {
	return a.repo.Update(ctx, art)
}
