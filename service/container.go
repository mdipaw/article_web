package service

import (
	"article_web/article"
	"article_web/worker"

	"gorm.io/gorm"
)

type Container struct {
	ArticleReader *article.ArticleReader
	WorkerClient  *worker.WorkerClient
}

func New(dbReader, dbWorker *gorm.DB, workerClient *worker.WorkerClient) Container {
	articleReader := article.NewArticleReader(dbReader)
	return Container{
		ArticleReader: articleReader,
		WorkerClient:  workerClient,
	}
}
