package main

import (
	"dagger/CribService/internal/dagger"
	"dagger/CribService/utils"
)

func (c *CribService) golang() *dagger.Container {
	return dag.Container().
		From("golang:latest")
}

func (c *CribService) sqlc(src *dagger.Directory) *dagger.Container {
	g := c.golang()
	g = g.WithExec([]string{"go", "install", "github.com/sqlc-dev/sqlc/cmd/sqlc@latest"})
	g = g.WithDirectory("/src", src.Directory("sql"))
	return g.
		WithWorkdir("/src").
		WithExec([]string{"sqlc", "generate"})
}

func (c *CribService) postgres(withPort bool) *dagger.Container {
	p := dag.Container().
		From("postgres:latest").
		WithEnvVariable("POSTGRES_USER", "postgres").
		WithEnvVariable("POSTGRES_PASSWORD", "example").
		WithEnvVariable("POSTGRES_DB", "cribbage")

	if withPort {
		p = p.WithExposedPort(5432)
	}

	return p
}

func (c *CribService) gameServer(src *dagger.Directory) *dagger.Container {
	server := c.golang()
	server = server.WithExec([]string{"go", "install", "github.com/swaggo/swag/cmd/swag@latest"})
	server = utils.GoMod(src.Directory("server"), server)
	server = server.WithDirectory("/src", src)
	return server.
		WithoutEntrypoint().
		WithWorkdir("/src/server").
		WithExec([]string{
			"/go/bin/swag",
			"init",
			"-d",
			"/src/server",
			"--parseDependency",
			"--parseInternal",
		})
}

func (i *CribService) http(src *dagger.Directory) *dagger.Container {
	ij := dag.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "openjdk17-jdk", "curl", "unzip"}).
		WithExec([]string{"/bin/sh", "-c", "mkdir /ijhttp"}).
		WithExec([]string{"curl", "-f", "-L", "-o", "/ijhttp/ijhttp.zip", "https://jb.gg/ijhttp/latest"}).
		WithExec([]string{"unzip", "/ijhttp/ijhttp.zip"}).
		WithExec([]string{"/bin/sh", "-c", "chmod +x /ijhttp/ijhttp"}).
		WithDirectory("/workdir", src.Directory("http"))

	return ij
}
