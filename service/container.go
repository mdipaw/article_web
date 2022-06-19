package service

import (
	"article_web/article"
	"article_web/redis"
	"article_web/worker"
	"time"

	"gorm.io/gorm"
)

type Container struct {
	ArticleReader *article.ArticleReader
	WorkerClient  *worker.WorkerClient
	RedisClient   *redis.RedisClient
}

func New(dbReader, dbWorker *gorm.DB, workerClient *worker.WorkerClient) Container {
	articleReader := article.NewArticleReader(dbReader)
	redisClient := redis.NewRedis(time.Second * 5)
	return Container{
		ArticleReader: articleReader,
		WorkerClient:  workerClient,
		RedisClient:   redisClient,
	}
}
