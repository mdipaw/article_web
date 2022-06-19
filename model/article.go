package model

import (
	"time"

	"gorm.io/gorm"
)

type Article struct {
	ID      int `gorm:"primaryKey"`
	Author  string
	Title   string
	Body    string
	Created time.Time
}

type ArticleFilter struct {
	Query  string
	Author string
}

func (a *Article) BeforeCreate(tx *gorm.DB) error {
	if a.Created.IsZero() {
		a.Created = time.Now()
	}
	return nil
}

func (Article) TableName() string {
	return "articles"
}
