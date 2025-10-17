package service

import (
	"context"
	"jike_gin/internal/domain"
	"jike_gin/internal/repository/article"
)

type ArticleService interface {
	Save(ctx context.Context, article domain.Article) (id int64, err error)
	Update(ctx context.Context, article domain.Article) error
	Publish(ctx context.Context, article domain.Article) (int64, error)
	PublishV1(ctx context.Context, article domain.Article) (int64, error)
}

type articleService struct {
	repo   article.ArticleRepository
	author article.ArticleAuthorRepository
	reader article.ArticleReaderRepository
}

func NewArticleService(repo article.ArticleRepository) ArticleService {
	return &articleService{
		repo: repo,
	}
}

func NewArticleServiceV1(author article.ArticleAuthorRepository,
	reader article.ArticleReaderRepository) ArticleService {
	return &articleService{
		author: author,
		reader: reader,
	}
}

func (a *articleService) Save(ctx context.Context, art domain.Article) (int64, error) {
	if art.Id >= 0 {
		err := a.repo.Update(ctx, art)
		return art.Id, err
	}
	return a.repo.Create(ctx, art)
}

func (a *articleService) Update(ctx context.Context, art domain.Article) error {
	return a.repo.Update(ctx, art)
}

func (a *articleService) Publish(ctx context.Context, art domain.Article) (int64, error) {
	return a.repo.Sync(ctx, art)
}

func (a *articleService) PublishV1(ctx context.Context, art domain.Article) (int64, error) {
	var (
		id  = art.Id
		err error
	)
	if art.Id > 0 {
		err = a.author.Update(ctx, art)
	} else {
		id, err = a.author.Create(ctx, art)
	}
	if err != nil {
		return 0, err
	}
	art.Id = id

	return a.reader.Save(ctx, art)
}
