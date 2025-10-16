package repository

import (
	"context"
	"jike_gin/internal/domain"
	"jike_gin/internal/repository/dao"
)

type ArticleRepository interface {
	Create(ctx context.Context, article domain.Article) (int64, error)
	Update(ctx context.Context, article domain.Article) error
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

func (a *articleService) Update(ctx context.Context, article domain.Article) error {
	return a.dao.UpdateById(ctx, dao.Article{
		Id:       article.Id,
		AuthorId: article.Author.Id,
		Title:    article.Title,
		Content:  article.Content,
	})
}
