package article

import (
	"context"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"jike_gin/internal/domain"
	"jike_gin/internal/repository/dao/article"
)

type ArticleRepository interface {
	Create(ctx context.Context, article domain.Article) (int64, error)
	Sync(ctx context.Context, article domain.Article) (int64, error)
	Update(ctx context.Context, article domain.Article) error
}

type CachedArticleRepository struct {
	dao       article.ArticleDAO
	readerDao article.ArticleReaderDao
	authorDao article.ArticleAuthorDao

	db *gorm.DB
}

func NewArticleRepository(dao article.ArticleDAO, readerDao article.ArticleReaderDao, authorDao article.ArticleAuthorDao, db *gorm.DB) ArticleRepository {
	return &CachedArticleRepository{
		dao:       dao,
		readerDao: readerDao,
		authorDao: authorDao,
		db:        db,
	}
}

func (c *CachedArticleRepository) Create(ctx context.Context, art domain.Article) (int64, error) {
	return c.dao.Insert(ctx, article.Article{
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
	})
}

// SyncV2 开启事务确保都成功.第二种写法
func (c *CachedArticleRepository) SyncV2(ctx context.Context, art domain.Article) (int64, error) {
	tx := c.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		zap.L().Error("SyncV2", zap.Error(tx.Error))
		return 0, tx.Error
	}
	//author := dao.NewGORMArticleDAO(tx)
	//reader := dao.NewGORMArticleDAO(tx)
	// .....
	return 0, nil
}

//func (c *CachedArticleRepository) SyncV2_1(ctx context.Context, art domain.Article) (int64, error) {
//	c.dao.Transaction()
//}

func (c *CachedArticleRepository) Sync(ctx context.Context, art domain.Article) (int64, error) {
	var (
		id  = art.Id
		err error
	)
	artn := c.toEntity(ctx, art)
	if id > 0 {
		err = c.authorDao.UpdateById(ctx, artn)
	} else {
		id, err = c.authorDao.Insert(ctx, artn)
	}
	if err != nil {
		zap.L().Error("Sync err", zap.Error(err))
		return id, err
	}
	c.readerDao.Upsert(ctx, artn)
	return id, err
}

func (c *CachedArticleRepository) toEntity(ctx context.Context, art domain.Article) article.Article {
	return article.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
	}
}

func (c *CachedArticleRepository) Update(ctx context.Context, art domain.Article) error {
	return c.dao.UpdateById(ctx, article.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
	})
}
