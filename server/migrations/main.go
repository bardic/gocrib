package migrations

import (
	"fmt"

	"github.com/uptrace/bun/migrate"
)

var Migrations = migrate.NewMigrations()

func New() {
	fmt.Print()
	if err := Migrations.DiscoverCaller(); err != nil {
		panic(err)
	}
}
