package main

import (
	"dagger/CribService/internal/dagger"
)

func (c *CribService) postgresService(
	// +optional
	withPort bool,
) {
	c.Db = c.
		postgres(withPort).
		AsService(dagger.ContainerAsServiceOpts{UseEntrypoint: true}).
		WithHostname("db")
}

func (c *CribService) startGameServer(src *dagger.Directory) *dagger.Service {
	c.postgresService(false)
	c.migrate(src)

	return c.gameServer(src).
		WithServiceBinding("db", c.Db).
		WithExposedPort(1323).
		WithDefaultArgs([]string{"go", "run", "main.go"}).
		AsService().
		WithHostname("server")
}
