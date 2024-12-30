package main

import (
	"context"
	"dagger/CribService/internal/dagger"
	"dagger/CribService/utils"
	"strings"
)

func (c *CribService) golang() *dagger.Container {
	return dag.Container().
		From("golang:latest")
}

func (c *CribService) goWithSrc() *dagger.Container {
	return c.golang()
}

func (c *CribService) sqlc(src *dagger.Directory) *dagger.Container {
	g := c.goWithSrc()
	g = g.WithExec([]string{"go", "install", "github.com/sqlc-dev/sqlc/cmd/sqlc@latest"})
	g = g.WithDirectory("/src", src.Directory("sql"))
	return g.
		WithWorkdir("/src").
		WithExec([]string{"sqlc", "generate"}).
		WithEntrypoint([]string{"ls"})
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

func (c *CribService) migrate(src *dagger.Directory) *dagger.Container {
	migration := c.golang()
	migration = migration.WithExec([]string{
		"go",
		"install",
		"-tags",
		"'postgres'",
		"github.com/golang-migrate/migrate/v4/cmd/migrate@latest",
	})
	migration = migration.WithDirectory("/src", src.Directory("sql"))
	migration = utils.GoMod(migration, src.WithDirectory("sql", src))
	migration = migration.WithServiceBinding("db", c.Db)

	return migration.
		WithDefaultArgs([]string{
			"migrate",
			"-path",
			"migrations",
			"-database",
			"postgres://postgres:example@db:5432/cribbage?sslmode=disable",
			"up",
		})
}

func (c *CribService) gameServer(src *dagger.Directory) *dagger.Container {
	server := c.goWithSrc()
	server = server.WithExec([]string{"go", "install", "github.com/swaggo/swag/cmd/swag@latest"})
	server = server.WithDirectory("/src", src)
	server = utils.GoMod(server, src.WithDirectory("server", src))
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

func (i *CribService) http(ctx context.Context, src *dagger.Directory) *dagger.Container {
	ij := dag.Container().
		From("alpine:latest").
		WithDirectory("/workdir", src.Directory("http")).
		WithExec([]string{"apk", "add", "openjdk17-jdk", "curl", "unzip"}).
		WithExec([]string{"/bin/sh", "-c", "mkdir /ijhttp"}).
		WithExec([]string{"curl", "-f", "-L", "-o", "/ijhttp/ijhttp.zip", "https://jb.gg/ijhttp/latest"}).
		WithExec([]string{"unzip", "/ijhttp/ijhttp.zip"}).
		WithExec([]string{"/bin/sh", "-c", "chmod +x /ijhttp/ijhttp"})

	ij.WithServiceBinding("server", i.Server(ctx, src))

	entries, err := src.Directory("http").Entries(context.Background())

	if err != nil {
		return nil
	}

	f := make([]string, 0)
	for _, file := range entries {
		if strings.HasSuffix(file, ".http") {
			f = append(f, "/workdir/"+file)
		}
	}

	ij = ij.WithExec(append([]string{"sh", "/ijhttp/ijhttp"}, f...))

	return ij
}
