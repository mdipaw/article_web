package worker

import (
	"article_web/article"

	"github.com/hibiken/asynq"
	"gorm.io/gorm"
)

type WorkerServer struct {
	srv    *asynq.Server
	router *asynq.ServeMux
}

func NewServer(redisAddr string, dbWorker, dbReader *gorm.DB) WorkerServer {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{
			Concurrency: 10,
		},
	)
	articleWorker := article.NewWorkerArticle(dbWorker, dbReader)
	router := asynq.NewServeMux()
	router.HandleFunc(article.ArticlePost, articleWorker.WorkerArticlePost())

	return WorkerServer{
		srv,
		router,
	}
}

func (w *WorkerServer) Run() error {
	return w.srv.Run(w.router)
}

func (w *WorkerServer) Shutdown() {
	w.srv.Shutdown()
}

func (w *WorkerServer) Start(handler asynq.Handler) error {
	return w.srv.Start(handler)
}
