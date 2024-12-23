package main

import (
	"context"
	"dagger/integration/internal/dagger"
	"strings"
)

type Integration struct {
	Db     *dagger.Service
	Server *dagger.Service
}

func (i *Integration) postgresDB() *dagger.Service {
	return dag.Container().
		From("postgres:latest").
		WithEnvVariable("POSTGRES_USER", "postgres").
		WithEnvVariable("POSTGRES_PASSWORD", "example").
		WithEnvVariable("POSTGRES_DB", "cribbage").
		AsService(dagger.ContainerAsServiceOpts{UseEntrypoint: true}).
		WithHostname("db")
}

func (i *Integration) migrationService(src *dagger.Directory) *dagger.Service {
	migration := dag.Container().
		From("golang:latest").
		WithServiceBinding("db", i.Db)

	mSrc := src.Directory("queries")
	migration = gomod(migration, mSrc)

	migration = migration.WithExec([]string{
		"go",
		"install",
		"-tags",
		"'postgres'",
		"github.com/golang-migrate/migrate/v4/cmd/migrate@latest",
	})

	migration = exclude(migration, src)

	return migration.
		WithDefaultArgs([]string{
			"migrate",
			"-path",
			"/src/queries/migrations",
			"-database",
			"postgres://postgres:example@db:5432/cribbage?sslmode=disable",
			"up",
		}).
		AsService()
}

func (i *Integration) swagger(src *dagger.Directory) *dagger.Service {

	server := dag.Container().
		From("golang:latest")

	sSrc := src.Directory("server")
	server = gomod(server, sSrc)
	server = server.WithExec([]string{"go", "install", "github.com/swaggo/swag/cmd/swag@latest"})
	server = exclude(server, src)

	server = server.
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

	return server.
		WithServiceBinding("db", i.Db).
		WithExposedPort(1323).
		WithDefaultArgs([]string{"go", "run", "main.go"}).
		AsService().
		WithHostname("server")
}

func (i *Integration) ijhttp(src *dagger.Directory) *dagger.Container {
	ij := dag.Container().
		From("alpine:latest").
		WithDirectory("/workdir", src).
		WithExec([]string{"apk", "add", "openjdk17-jdk", "curl", "unzip"}).
		WithExec([]string{"/bin/sh", "-c", "mkdir /ijhttp"}).
		WithExec([]string{"curl", "-f", "-L", "-o", "/ijhttp/ijhttp.zip", "https://jb.gg/ijhttp/latest"}).
		WithExec([]string{"unzip", "/ijhttp/ijhttp.zip"}).
		WithExec([]string{"/bin/sh", "-c", "chmod +x /ijhttp/ijhttp"}).
		WithServiceBinding("server", i.Server)

	entries, err := src.Entries(context.Background())

	if err != nil {
		return nil
	}

	f := make([]string, 0)
	for _, file := range entries {
		if strings.HasSuffix(file, ".http") {
			//f += "/workdir/" + file + " "
			f = append(f, "/workdir/"+file)
			//ij = ij.WithExec([]string{"sh", "/ijhttp/ijhttp", "/workdir/" + file, " >> /workdir/output.txt"})
		}
	}

	ij = ij.WithExec(append([]string{"sh", "/ijhttp/ijhttp"}, f...))

	return ij

}

func (i *Integration) TestSwagger(ctx context.Context, src *dagger.Directory) (*dagger.Service, error) {
	db, err := i.postgresDB().
		Start(ctx)

	if err != nil {
		return nil, err
	}

	i.Db = db

	server, err := i.swagger(src).
		Start(ctx)

	if err != nil {
		return nil, err
	}

	migration, err := i.migrationService(src).
		Start(ctx)

	if err != nil {
		return nil, err
	}

	defer migration.Stop(ctx)

	return server, nil
}

func (i *Integration) Test(ctx context.Context, src *dagger.Directory) (string, error) {
	db, err := i.postgresDB().
		Start(ctx)

	if err != nil {
		return "", err
	}

	i.Db = db

	server, err := i.swagger(src).
		Start(ctx)

	if err != nil {
		return "", err
	}

	i.Server = server

	migration, err := i.migrationService(src).
		Start(ctx)

	if err != nil {
		return "", err
	}

	defer migration.Stop(ctx)

	ij := i.ijhttp(src.Directory("integration/http"))

	return ij.Stdout(ctx)
}

func (i *Integration) TestPostgres(ctx context.Context, src *dagger.Directory) (*dagger.Service, error) {
	db := i.postgresDB()

	i.Db = db

	migration, err := i.migrationService(src).
		Start(ctx)

	if err != nil {
		return nil, err
	}

	defer migration.Stop(ctx)

	return db, nil
}

func exclude(c *dagger.Container, dir *dagger.Directory) *dagger.Container {
	return c.WithDirectory("/src", dir, dagger.ContainerWithDirectoryOpts{
		Exclude: []string{
			"./integration/.dagger/internal",
			"./.git",
			"./integration/.dagger/dagger.gen.go",
		},
	})
}

func gomod(c *dagger.Container, dir *dagger.Directory) *dagger.Container {
	return c.
		WithDirectory("/src", dir,
			dagger.ContainerWithDirectoryOpts{
				Include: []string{
					"go.mod",
					"go.sum",
				},
			}).
		WithWorkdir("/src").
		WithExec([]string{
			"go",
			"mod",
			"download",
		})
}
