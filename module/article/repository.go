package article

import (
	"context"
	"fmt"
	"github.com/uptrace/bun"
	"go-bunrouter-example/module/primitive"
	"strings"
)

type RepositoryInterface interface {
	CreateArticle(ctx context.Context, payload primitive.Article) (primitive.Article, error)
	CountArticle(ctx context.Context, param primitive.ParameterFindArticle) (int64, error)
	FindListArticle(ctx context.Context, param primitive.ParameterFindArticle) ([]primitive.Article, error)
	FindArticleByID(ctx context.Context, articleID int64) (primitive.Article, error)
	SetParamQueryToOrderByQuery(orderBy string) string
}

type Repository struct {
	db *bun.DB
}

func NewRepository(db *bun.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreateArticle(ctx context.Context, payload primitive.Article) (primitive.Article, error) {
	if _, err := r.db.NewInsert().Model(payload).Table("articles").Exec(ctx); err != nil {
		return payload, err
	}
	return payload, nil
}

func (r *Repository) CountArticle(ctx context.Context, param primitive.ParameterFindArticle) (int64, error) {
	query := r.db.NewSelect().Table("articles")
	query.Where(`"deleted_at" is null`)
	if param.Author != "" {
		query.Where(`"author" ILIKE ?`, "%"+param.Author+"%")
	}
	if param.Query != "" {
		query.Where(`"title" ILIKE ? or "body" ILIKE ?`, "%"+param.Query+"%", "%"+param.Query+"%")
	}
	count, err := query.Count(ctx)
	if err != nil {
		return 0, err
	}
	return int64(count), nil
}

func (r *Repository) FindListArticle(ctx context.Context, param primitive.ParameterFindArticle) ([]primitive.Article, error) {
	var listData []primitive.Article

	query := r.db.NewSelect().Model(&listData)
	query.Where(`"deleted_at" is null`)
	if param.Author != "" {
		query.Where(`"author" ILIKE ?`, "%"+param.Author+"%")
	}
	if param.Query != "" {
		query.Where(`"title" ILIKE ? or "body" ILIKE ?`, "%"+param.Query+"%", "%"+param.Query+"%")
	}

	err := query.Offset(param.Offset).
		Limit(param.PageSize).
		Order(strings.Join([]string{param.SortBy, param.SortOrder}, " ")).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	return listData, nil
}

func (r *Repository) SetParamQueryToOrderByQuery(orderBy string) string {
	var result string
	switch orderBy {
	case "id":
		result = fmt.Sprintf(`id`)
	case "author":
		result = fmt.Sprintf(`author`)
	case "title":
		result = fmt.Sprintf(`title`)
	case "body":
		result = fmt.Sprintf(`body`)
	case "created":
		result = fmt.Sprintf(`created_at`)
	default:
		result = fmt.Sprintf(`created_at`)
	}
	return result
}

func (r *Repository) FindArticleByID(ctx context.Context, articleID int64) (primitive.Article, error) {
	var data primitive.Article
	err := r.db.NewSelect().
		Model(&data).
		Where(`"deleted_at" is null and "id" = ?`, articleID).
		Scan(ctx)
	if err != nil {
		return primitive.Article{}, err
	}

	return data, nil
}
