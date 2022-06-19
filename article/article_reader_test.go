package article_test

import (
	"article_web/article"
	"article_web/model"
	"article_web/tests"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type fixture struct {
	dbReader *gorm.DB
	dbWriter *gorm.DB
}

type testFunc func(fixture)

func runTest(f testFunc) {
	tests.NewDBTest(func(dbReader, dbWriter *gorm.DB) {
		f(fixture{
			dbReader: dbReader,
			dbWriter: dbWriter,
		})
	})
}

func TestGetArticle(t *testing.T) {
	runTest(func(f fixture) {
		data := model.Article{
			Author: "john doe",
			Title:  "Awesome post",
			Body:   "this is awesome post",
		}
		task := article.CreateArticlePostTask(data)
		workers := article.NewWorkerArticle(f.dbWriter, f.dbReader)
		workers.WorkerArticlePost()(context.Background(), task)

		reader := article.NewArticleReader(f.dbReader)
		retrievedArticle, err := reader.GetQuery(model.ArticleFilter{
			Author: data.Author,
		}).First()
		assert.NoError(t, err)

		assert.Equal(t, data.Author, retrievedArticle.Author)
	})
}
