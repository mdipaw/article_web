package public

import (
	"article_web/article"
	"article_web/model"
	"article_web/rest"
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

		query := p.di.ArticleReader.GetQuery(model.ArticleFilter{
			Query:  paramQuery,
			Author: paramAuthor,
		}).Query

		paginator := rest.PaginatorFromContext[[]model.Article, []ArticleResponse](c)
		paginatedResponse := paginator.PaginateQuery(query).
			Map(func(x []model.Article) []ArticleResponse {
				result := []ArticleResponse{}
				for _, X := range x {
					result = append(result, ArticleResponse{
						ID:      X.ID,
						Author:  X.Author,
						Title:   X.Title,
						Body:    X.Body,
						Created: X.Created,
					})
				}
				return result
			})

		if paginatedResponse.Error != nil {
			c.JSON(http.StatusBadRequest,
				gin.H{"error": paginatedResponse.Error.Error()})
			return
		}
		// TODO:implement caching
		c.JSON(http.StatusOK, paginatedResponse)
		return
	}
}
