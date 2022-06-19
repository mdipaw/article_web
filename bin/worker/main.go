package main

import (
	"article_web/database"
	"article_web/redis"
	"article_web/worker"
	"log"

	"gorm.io/gorm/logger"
)

func main() {
	dbWorker := database.ConnectDatabaseWriter(database.DsnWriter, logger.Default.LogMode(logger.Error))
	dbReader := database.ConnectDatabaseReader(database.DsnReader, logger.Default.LogMode(logger.Error))

	server := worker.NewServer(redis.RedisAddress, dbWorker, dbReader)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}

}
