package migrations

import (
	"embed"

	"github.com/pressly/goose/v3"
	"gorm.io/gorm"
)

//go:embed sql/*.sql
var embedMigrations embed.FS

func Up(model *gorm.DB) {
	db, err := model.DB()
	if err != nil {
		panic(err)
	}

	goose.SetBaseFS(embedMigrations)

	if err := goose.Up(db, "sql", goose.WithAllowMissing()); err != nil {
		panic(err)
	}

}
