package article_test

import (
	"article_web/article"
	"article_web/model"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateArticle(t *testing.T) {
	runTest(func(f fixture) {
		data := model.Article{
			Author: "john doe",
			Title:  "Awesome post",
			Body:   "this is awesome post",
		}
		task := article.CreateArticlePostTask(data)
		workers := article.NewWorkerArticle(f.dbWriter, f.dbReader)
		err := workers.WorkerArticlePost()(context.Background(), task)
		assert.NoError(t, err)
	})
}
