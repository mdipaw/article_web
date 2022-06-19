package server

import (
	"article_web/redis"
	"article_web/tests"
	"article_web/worker"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type fixture struct {
	server *Server
}
type testFunc func(f fixture)

func runTest(testFunc testFunc) {
	tests.NewDBTest(func(dbReader, dbWorker *gorm.DB) {
		workerClient := worker.NewWorkerClient(redis.RedisAddress)
		defer workerClient.Client.Close()
		server := NewServer(dbWorker, dbReader, workerClient)

		testFunc(fixture{
			server: server,
		})
	})
}

func TestStartServer(t *testing.T) {
	runTest(func(f fixture) {
		isCalled := false
		f.server.engine.GET("/test", func(ctx *gin.Context) {
			isCalled = true
			ctx.Status(200)
			return
		})

		w := httptest.NewRecorder()
		r, _ := http.NewRequest("GET", "/test", nil)
		f.server.ServeHTTP(w, r)

		assert.Equal(t, 200, w.Code)
		assert.True(t, isCalled)
	})
}
