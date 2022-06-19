package rest

import (
	"article_web/redis"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CacheMiddlewareFunc func(redis *redis.RedisClient) gin.HandlerFunc

func CacheMiddleware(redis *redis.RedisClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		data, err := redis.Get("getAllArticle")
		if err != nil {
			c.Next()
			return
		}
		c.AbortWithStatusJSON(http.StatusOK, data)
		return
	}
}
