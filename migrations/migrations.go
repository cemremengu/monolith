package migrations

import (
	"database/sql"
	"embed"

	"github.com/pressly/goose/v3"
)

//go:embed *.sql
var embedMigrations embed.FS

func Up(db *sql.DB) {
	if err := goose.Up(db, "."); err != nil {
		panic(err)
	}
}

func Down(db *sql.DB) {
	if err := goose.Down(db, "."); err != nil {
		panic(err)
	}
}

// SetDialect sets the dialect for the migrations
// default is postgres
func SetDialect(dialect string) error {
	return goose.SetDialect(dialect)
}

func init() {
	goose.SetBaseFS(embedMigrations)
}
