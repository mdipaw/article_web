package article

import (
	"article_web/database"
	"article_web/model"

	"gorm.io/gorm"
)

type Article struct {
	dbReader *gorm.DB
}

func NewArticle(dbReader *gorm.DB) *Article {
	return &Article{dbReader}
}

type superGetQuery = database.GetQuery[model.Article]

type thisGetQuery struct {
	*superGetQuery
}

func (a *Article) GetQuery(filter model.ArticleFilter) thisGetQuery {
	return thisGetQuery{
		database.NewQueryGeneric[model.Article, model.ArticleFilter]().
			GetQuery(func(f model.ArticleFilter) *gorm.DB {
				query := a.dbReader.Model(model.Article{})
				if filter.Author != "" {
					query = query.Where("author = ?", filter.Author)
				}
				return query
			}, filter)}
}
