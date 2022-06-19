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
	workerDB *gorm.DB
	readerDB *gorm.DB
}

func NewWorkerArticle(workerDB, readerDB *gorm.DB) *WorkerArticle {
	return &WorkerArticle{workerDB, readerDB}
}

func (w *WorkerArticle) WorkerArticlePost() asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {
		byt := []byte(t.Payload())
		articleData := model.Article{}
		json.Unmarshal(byt, &articleData)

		if err := w.workerDB.Model(model.Article{}).
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
