package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectDatabaseWriter(dsn string, l logger.Interface) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: l,
	})

	if err != nil {
		panic("Failed to connect to database writer")
	}
	return db
}

func ConnectDatabaseReader(dsn string, l logger.Interface) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: l,
	})

	if err != nil {
		panic("Failed to connect to database writer")
	}
	return db
}

type initQuery[T any, F any] struct {
	GetQuery func(filter func(f F) *gorm.DB, filterStruct F) *GetQuery[T]
}

type GetQuery[T any] struct {
	Query *gorm.DB
}

func NewQueryGeneric[T any, F any]() *initQuery[T, F] {
	return &initQuery[T, F]{GetQuery: getQuery[T, F]}
}

func getQuery[T any, F any](filter func(f F) *gorm.DB, filterStruct F) *GetQuery[T] {
	queryAfterFiltering := filter(filterStruct)
	return &GetQuery[T]{Query: queryAfterFiltering}
}

func (super *GetQuery[T]) First() (T, error) {
	var x T
	if err := super.Query.First(&x).Error; err != nil {
		return x, err
	}
	return x, nil
}

func (super *GetQuery[T]) Find() ([]T, error) {
	var x []T
	if err := super.Query.Find(&x).Error; err != nil {
		return x, err
	}
	return x, nil
}
