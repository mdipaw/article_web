package public

import (
	"article_web/article"
	"article_web/model"
	"article_web/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type PublicHandler struct {
	di *service.Container
}

func NewPublicHandler(di *service.Container) *PublicHandler {
	return &PublicHandler{di}
}

func (p *PublicHandler) RegisterRoutes(engine *gin.RouterGroup) *gin.RouterGroup {
	engine.POST("/article", p.PostArticle())
	engine.GET("/article", p.GetArticle())
	return engine
}

func (p *PublicHandler) PostArticle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var params struct {
			Author  string    `json:"author" validate:"required"`
			Title   string    `json:"title" validate:"required"`
			Body    string    `json:"body" validate:"required"`
			Created time.Time `json:"created"`
		}

		if err := ctx.BindJSON(&params); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		p.di.WorkerClient.Client.Enqueue(article.CreateArticlePostTask(model.Article{
			Author:  params.Author,
			Title:   params.Title,
			Body:    params.Body,
			Created: params.Created,
		}))
		ctx.Status(200)
		return
	}
}

func (p *PublicHandler) GetArticle() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}
