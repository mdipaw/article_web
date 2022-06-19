package public

import (
	"article_web/service"

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

	}
}

func (p *PublicHandler) GetArticle() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}
