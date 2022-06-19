package article

import (
	"article_web/database"
	"article_web/model"

	"gorm.io/gorm"
)

type ArticleReader struct {
	dbReader *gorm.DB
}

func NewArticleReader(dbReader *gorm.DB) *ArticleReader {
	return &ArticleReader{dbReader}
}

type superGetQuery = database.GetQuery[model.Article]

type thisGetQuery struct {
	*superGetQuery
}

func (a *ArticleReader) GetQuery(filter model.ArticleFilter) thisGetQuery {
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
