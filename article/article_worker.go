package article

import (
	"article_web/model"
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
	"gorm.io/gorm"
)

const (
	ArticlePost = "article:post"
)

type WorkerArticle struct {
	writerDB *gorm.DB
	readerDB *gorm.DB
}

func NewWorkerArticle(writerDB, readerDB *gorm.DB) *WorkerArticle {
	return &WorkerArticle{writerDB, readerDB}
}

func (w *WorkerArticle) WorkerArticlePost() asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {
		byt := []byte(t.Payload())
		articleData := model.Article{}
		json.Unmarshal(byt, &articleData)

		if err := w.writerDB.Model(model.Article{}).
			Create(&articleData).Error; err != nil {
			return err
		}
		if err := w.readerDB.Model(model.Article{}).
			Create(&articleData).Error; err != nil {
			return err
		}
		return nil
	}
}

func CreateArticlePostTask(data model.Article) *asynq.Task {
	byt, _ := json.Marshal(data)
	return asynq.NewTask(ArticlePost, byt)
}
