package server

import (
	"article_web/rest"
	"article_web/server/public"
	"article_web/service"
	"article_web/worker"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Server struct {
	engine *gin.Engine
}

func NewServer(dbWorker, dbReader *gorm.DB, workerClient worker.WorkerClient) *Server {
	di := service.New(dbReader, dbWorker, &workerClient)
	r := gin.New()
	baseRoute := r.Group("/api")
	publicHandler := public.NewPublicHandler(&di)
	publicHandler.RegisterRoutes(baseRoute, rest.CacheMiddleware)

	return &Server{engine: r}
}

func (s *Server) Run() {
	s.engine.Run()
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.engine.ServeHTTP(w, r)
}
