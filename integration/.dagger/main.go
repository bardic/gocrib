package main

import (
	"context"
	"dagger/integration/internal/dagger"
)

type Integration struct{}

func (i *Integration) postgresDB() *dagger.Service {
	return dag.Container().
		From("postgres:latest").
		WithEnvVariable("POSTGRES_USER", "postgres").
		WithEnvVariable("POSTGRES_PASSWORD", "example").
		WithEnvVariable("POSTGRES_DB", "cribbage").
		WithExposedPort(5432).
		AsService().
		WithHostname("db")
}

func (i *Integration) migrationService(src *dagger.Directory, db *dagger.Service) *dagger.Service {
	migration := dag.Container().
		From("golang:latest").
		WithServiceBinding("db", db)

	migration = gomod(migration, src)

	migration = migration.WithExec([]string{
		"go",
		"install",
		"-tags",
		"'postgres'",
		"github.com/golang-migrate/migrate/v4/cmd/migrate@latest",
	})

	migration = exclude(migration, src)

	return migration.
		WithExec([]string{
			"migrate",
			"-path",
			"/src/queries/migrations",
			"-database",
			"postgres://postgres:example@db:5432/cribbage?sslmode=disable",
			"up",
		}).
		AsService()
}

func (i *Integration) swagger(db *dagger.Service, src *dagger.Directory) *dagger.Service {

	server := dag.Container().
		From("golang:latest")

	server = gomod(server, src)
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
		WithServiceBinding("db", db).
		WithExec([]string{"go", "run", "main.go"}).
		WithExposedPort(1323).
		AsService().
		WithHostname("server")
}

func (i *Integration) TestSwagger(ctx context.Context, src *dagger.Directory) (*dagger.Service, error) {
	db, err := i.postgresDB().
		Start(ctx)

	if err != nil {
		return nil, err
	}

	defer db.Stop(ctx)

	server, err := i.swagger(db, src).
		Start(ctx)

	if err != nil {
		return nil, err
	}

	migration, err := i.migrationService(src, db).
		Start(ctx)

	if err != nil {
		return nil, err
	}

	defer migration.Stop(ctx)

	return server, nil
}

func (i *Integration) TestPostgres(ctx context.Context, src *dagger.Directory) (*dagger.Service, error) {
	db := i.postgresDB()

	migration, err := i.migrationService(src, db).
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
	return c.WithDirectory("/src", dir, dagger.ContainerWithDirectoryOpts{
		Include: []string{
			"go.mod",
			"go.sum",
		},
	})
}
