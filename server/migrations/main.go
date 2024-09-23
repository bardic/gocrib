package migrations

import (
	"github.com/uptrace/bun/migrate"
)

var Migrations = migrate.NewMigrations()

func New() {
	if err := Migrations.DiscoverCaller(); err != nil {
		panic(err)
	}
}
