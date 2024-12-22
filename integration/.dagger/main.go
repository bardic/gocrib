package main

import (
	"context"
	"dagger/integration/internal/dagger"
)

type Integration struct{}

func (i *Integration) PostgresDB() *dagger.Service {
	return dag.Container().
		From("postgres:latest").
		WithEnvVariable("POSTGRES_USER", "postgres").
		WithEnvVariable("POSTGRES_PASSWORD", "example").
		WithEnvVariable("POSTGRES_DB", "cribbage").
		WithExposedPort(5432).
		AsService().
		WithHostname("db")
}

func (i *Integration) MigrationService(src *dagger.Directory, db *dagger.Service) *dagger.Service {
	migration := dag.Container().
		From("golang:latest").
		WithServiceBinding("db", db)

	migration = exclude(migration, src)

	return migration.
		WithExec([]string{"go", "install", "-tags", "'postgres'", "github.com/golang-migrate/migrate/v4/cmd/migrate@latest"}).
		WithExec([]string{"migrate", "-path", "/src/queries/migrations", "-database", "postgres://postgres:example@db:5432/cribbage?sslmode=disable", "up"}).
		AsService()
}

func (i *Integration) RestServer(src *dagger.Directory) *dagger.Service {
	db := i.PostgresDB()

	server := dag.Container().
		From("golang:latest").
		WithServiceBinding("db", db).
		WithServiceBinding("migration", i.MigrationService(src, db)).
		WithExposedPort(1323)

	server = exclude(server, src)

	return server.
		WithWorkdir("/src").
		WithExec([]string{"go", "run", "/src/server/main.go"}).
		AsService().WithHostname("server")
}

func (i *Integration) Swagger(src *dagger.Directory) *dagger.Service {
	db, err := i.PostgresDB().Start(context.Background())
	if err != nil {
		return nil
	}

	migration, err := i.MigrationService(src, db).Start(context.Background())

	if err != nil {
		return nil
	}

	defer db.Stop(context.Background())
	defer migration.Stop(context.Background())

	server := dag.Container().
		From("golang:latest")

	server = exclude(server, src)

	server = server.
		WithoutEntrypoint().
		WithWorkdir("/src/server").
		WithExec([]string{"go", "install", "github.com/swaggo/swag/cmd/swag@latest"}).
		WithExec([]string{
			"/go/bin/swag",
			"init",
			"-d",
			"/src/server",
			"--parseDependency",
			"--parseInternal",
		})
		// .
		// Sync(context.Background())

	// if err != nil {
	// 	return nil
	// }

	return server.WithExec([]string{"go", "run", "main.go"}).WithExposedPort(1323).AsService()
}

func (i *Integration) Test(src, test *dagger.Directory) (string, error) {

	return dag.Container().
		From("alpine:latest").
		WithServiceBinding("server", i.RestServer(src)).
		WithDirectory("/workdir", test).
		WithExec([]string{"apk", "add", "openjdk17-jdk", "curl", "unzip"}).
		WithExec([]string{"/bin/sh", "-c", "mkdir /ijhttp"}).
		WithExec([]string{"curl", "-f", "-L", "-o", "/ijhttp/ijhttp.zip", "https://jb.gg/ijhttp/latest"}).
		WithExec([]string{"unzip", "/ijhttp/ijhttp.zip"}).
		WithExec([]string{"/bin/sh", "-c", "chmod +x /ijhttp/ijhttp"}).
		WithExec([]string{"sh", "/ijhttp/ijhttp", "/workdir/match.http"}).
		Stdout(context.Background())
}

func exclude(c *dagger.Container, dir *dagger.Directory) *dagger.Container {
	return c.WithDirectory("/src", dir, dagger.ContainerWithDirectoryOpts{
		Exclude: []string{
			"./integration/.dagger/internal",
			"./.git",
			"./.vscode",
			".gitignore",
			"README.md",
			"UNLICENSE",
			"crib.log",
			"./integration/.dagger/dagger.gen.go",
		},
	})
}
