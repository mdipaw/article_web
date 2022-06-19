package model

import (
	"time"
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
