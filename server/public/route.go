package public

import (
	"article_web/article"
	"article_web/model"
	"article_web/rest"
	"article_web/service"
	"log"
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

func (p *PublicHandler) RegisterRoutes(engine *gin.RouterGroup, cacheMiddleWare rest.CacheMiddlewareFunc) *gin.RouterGroup {
	engine.POST("/article", p.PostArticle())
	engine.GET("/article", cacheMiddleWare(p.di.RedisClient), p.GetArticle())
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
		p.di.RedisClient.Delete("getAllArticle")
		ctx.Status(200)
		return
	}
}

type ArticleResponse struct {
	ID      int       `json:"id"`
	Author  string    `json:"author"`
	Title   string    `json:"title"`
	Body    string    `json:"body"`
	Created time.Time `json:"created"`
}

func (p *PublicHandler) GetArticle() gin.HandlerFunc {
	return func(c *gin.Context) {
		paramQuery := c.Query("query")
		paramAuthor := c.Query("author")

		article, err := p.di.ArticleReader.GetQuery(model.ArticleFilter{
			Query:  paramQuery,
			Author: paramAuthor,
		}).Find()
		if err != nil {
			return
		}

		result := []ArticleResponse{}
		for _, x := range article {
			result = append(result, ArticleResponse{
				ID:      x.ID,
				Author:  x.Author,
				Title:   x.Title,
				Body:    x.Body,
				Created: x.Created,
			})
		}
		log.Println(result)

		p.di.RedisClient.Set("getAllArticle", result)
		c.JSON(http.StatusOK, result)
		return
	}
}
