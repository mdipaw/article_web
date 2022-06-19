package tests

import (
	"article_web/database"
	"article_web/migrations"
	"sync"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var doOnce sync.Once

type testFunc func(*gorm.DB, *gorm.DB)

func NewDBTest(testFunc testFunc) {
	dbReader, dbWorker := setUpDB()
	db1, _ := dbReader.DB()
	db2, _ := dbWorker.DB()

	defer dbWorker.Rollback()
	defer dbReader.Rollback()
	defer db1.Close()
	defer db2.Close()

	testFunc(dbReader, dbWorker)

}

func setUpDB() (tx, tr *gorm.DB) {
	db := database.ConnectDatabaseReader(database.DsnReader, logger.Default.LogMode(logger.Error))
	db2 := database.ConnectDatabaseWriter(database.DsnWriter, logger.Default.LogMode(logger.Info))

	doOnce.Do(func() {
		migrations.Up(db)
		migrations.Up(db2)
	})
	tx = db.Begin()
	tr = db2.Begin()
	return
}
