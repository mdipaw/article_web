package public_test

import (
	"article_web/article"
	"article_web/model"
	"article_web/redis"
	"article_web/server"
	"article_web/server/public"
	"article_web/service"
	"article_web/tests"
	"article_web/worker"
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type fixture struct {
	server *server.Server
	worker *article.WorkerArticle
	di     *service.Container
}

type testFunc func(fixture)

func runTest(testFunc testFunc) {
	tests.NewDBTest(func(dbWorker, dbReader *gorm.DB) {
		workers := article.NewWorkerArticle(dbWorker, dbReader)
		workerClient := worker.NewWorkerClient(redis.RedisAddress)
		defer workerClient.Client.Close()
		di := service.New(dbReader, dbWorker, &workerClient)
		server := server.NewServer(dbWorker, dbReader, workerClient)

		testFunc(fixture{
			server: server,
			worker: workers,
			di:     &di,
		})
	})
}

func TestPostArticle(t *testing.T) {
	runTest(func(f fixture) {
		params := struct {
			Author  string    `json:"author" validate:"required"`
			Title   string    `json:"title" validate:"required"`
			Body    string    `json:"body" validate:"required"`
			Created time.Time `json:"created"`
		}{
			Author:  "test",
			Title:   "test title",
			Body:    "test body",
			Created: time.Now(),
		}

		jsonPayload, _ := json.Marshal(params)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest("POST", "/api/article", bytes.NewBuffer(jsonPayload))

		f.server.ServeHTTP(w, r)

		assert.Equal(t, 200, w.Code)
		// since creating article need different server to start the worker
		// and the server worker and server api cannot be run concurrently
		// we will manually invoke the job for post article
		task := article.CreateArticlePostTask(model.Article{
			Author:  params.Author,
			Title:   params.Title,
			Body:    params.Body,
			Created: params.Created,
		})

		err := f.worker.WorkerArticlePost()(context.Background(), task)
		assert.NoError(t, err)

		retrievedArticle, err := f.di.ArticleReader.GetQuery(model.ArticleFilter{
			Author: params.Author,
		}).First()
		assert.NoError(t, err)

		assert.Equal(t, params.Author, retrievedArticle.Author)
		assert.Equal(t, params.Title, retrievedArticle.Title)
		assert.Equal(t, params.Body, retrievedArticle.Body)
		isEqual := retrievedArticle.Created.Equal(params.Created)
		assert.True(t, isEqual)

	})
}

func TestGetArticle(t *testing.T) {
	runTest(func(f fixture) {
		var data = model.Article{
			Author:  "test",
			Body:    "test body",
			Title:   "test Title",
			Created: time.Now(),
		}
		f.di.ArticleReader.
			GetQuery(model.ArticleFilter{}).Query.
			Create(&data)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest("GET", "/api/article", nil)
		f.server.ServeHTTP(w, r)

		log.Println(w.Body)
		assert.Equal(t, 200, w.Code)

		var response []public.ArticleResponse
		json.Unmarshal(w.Body.Bytes(), &response)
		log.Println(response)

		assert.Equal(t, data.Author, response[0].Author)
		assert.Equal(t, data.Title, response[0].Title)
		assert.Equal(t, data.Body, response[0].Body)
		isEqual := response[0].Created.Equal(data.Created)
		assert.True(t, isEqual)
	})
}
