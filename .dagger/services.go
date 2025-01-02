package main

import (
	"context"
	"dagger/CribService/internal/dagger"
)

func (c *CribService) migrationService(src *dagger.Directory, p *dagger.Service) *dagger.Service {
	return c.golang().
		WithExec([]string{
			"go",
			"install",
			"-tags",
			"'postgres'",
			"github.com/golang-migrate/migrate/v4/cmd/migrate@latest"}).
		WithDirectory("/src", src).
		WithServiceBinding("db", p).
		WithExec([]string{
			"migrate",
			"-path",
			"/src/sql/migrations",
			"-database",
			"postgres://postgres:example@db:5432/cribbage?sslmode=disable",
			"up",
		}).
		AsService()
}

func (c *CribService) postgresService(
	// +optional
	withPort bool,
) *dagger.Service {
	p := c.
		postgres(withPort).
		AsService(dagger.ContainerAsServiceOpts{UseEntrypoint: true}).
		WithHostname("db")

	return p
}

func (c *CribService) gameServerService(ctx context.Context, src *dagger.Directory, migrate bool) (*dagger.Service, error) {
	c.Db = c.postgresService(true)
	if migrate {
		m := c.migrationService(src, c.Db)
		_, err := m.Start(ctx)

		if err != nil {
			return nil, err
		}
	}
	return c.
		gameServer(src).
		WithServiceBinding("db", c.Db).
		WithExposedPort(1323).
		WithDefaultArgs([]string{"go", "run", "main.go"}).AsService().
		WithHostname("server"), nil

}
