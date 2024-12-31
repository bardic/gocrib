package main

import (
	"dagger/CribService/internal/dagger"
)

func (c *CribService) postgresService(
	// +optional
	withPort bool,
) *dagger.Service {
	return c.
		postgres(withPort).
		AsService(dagger.ContainerAsServiceOpts{UseEntrypoint: true}).
		WithHostname("db")
}

func (c *CribService) startGameServer(src *dagger.Directory, migrate bool) *dagger.Service {
	c.Db = c.postgresService(true)
	s := c.
		gameServer(src).
		WithServiceBinding("db", c.Db).
		WithExposedPort(1323).
		WithDefaultArgs([]string{"go", "run", "main.go"})

	if migrate {
		s = s.WithExec([]string{
			"go",
			"install",
			"-tags",
			"'postgres'",
			"github.com/golang-migrate/migrate/v4/cmd/migrate@latest",
		})

		s = s.WithExec([]string{
			"migrate",
			"-path",
			"/src/sql/migrations",
			"-database",
			"postgres://postgres:example@db:5432/cribbage?sslmode=disable",
			"up",
		})
	}

	c.ServerContainer = s

	return c.ServerContainer.
		AsService().
		WithHostname("server")
}
