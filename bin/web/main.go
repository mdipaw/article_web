package main

import (
	"article_web/database"
	"article_web/migrations"
	"article_web/redis"
	"article_web/server"
	"article_web/worker"

	"gorm.io/gorm/logger"
)

func main() {
	dbWorker := database.ConnectDatabaseWriter(database.DsnWriter, logger.Default.LogMode(logger.Error))
	dbReader := database.ConnectDatabaseReader(database.DsnReader, logger.Default.LogMode(logger.Error))
	migrations.Up(dbWorker)
	migrations.Up(dbReader)

	workerClient := worker.NewWorkerClient(redis.RedisAddress)
	defer workerClient.Client.Close()
	srv := server.NewServer(dbWorker, dbReader, workerClient)
	srv.Run()
}
