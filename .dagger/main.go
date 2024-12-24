package main

import (
	"context"
	"dagger/CribService/internal/dagger"
	"strings"
)

type CribService struct {
	Db     *dagger.Service
	Server *dagger.Service
}

func (i *CribService) TestSqlc(ctx context.Context, src *dagger.Directory) *dagger.Container {
	return dag.Container().
		From("golang:latest").
		WithDirectory("/src", src.Directory("sql")).
		WithWorkdir("/src").
		WithExec([]string{"go", "install", "github.com/sqlc-dev/sqlc/cmd/sqlc@latest"}).
		WithEntrypoint([]string{"sqlc", "generate"})
}

func (i *CribService) Sqlc(ctx context.Context, src *dagger.Directory) *dagger.Directory {
	s := dag.Container().
		From("golang:latest").
		WithDirectory("/src", src.Directory("sql")).
		WithWorkdir("/src")

	s = gomod(s, src)

	return s.
		WithExec([]string{"go", "install", "github.com/sqlc-dev/sqlc/cmd/sqlc@latest"}).
		WithExec([]string{"sqlc", "generate"}).
		WithEntrypoint([]string{"sqlc", "generate"}).
		Directory("/src/queries")
}

func (i *CribService) GameServer(ctx context.Context, src *dagger.Directory) (*dagger.Service, error) {
	db, err := i.startPostgresDB(ctx)
	if err != nil {
		return nil, err
	}
	i.Db = db

	server, err := i.startSwagger(ctx, src)
	if err != nil {
		return nil, err
	}

	migration, err := i.startMigrationService(ctx, src)
	if err != nil {
		return nil, err
	}
	defer migration.Stop(ctx)

	return server, nil
}

func (i *CribService) Http(ctx context.Context, src *dagger.Directory) (string, error) {
	db, err := i.startPostgresDB(ctx)
	if err != nil {
		return "", err
	}
	i.Db = db

	server, err := i.startSwagger(ctx, src)
	if err != nil {
		return "", err
	}
	i.Server = server

	migration, err := i.startMigrationService(ctx, src)
	if err != nil {
		return "", err
	}
	defer migration.Stop(ctx)

	ij := i.ijhttp(src.Directory("http"))
	return ij.Stdout(ctx)
}

func (i *CribService) Postgres(ctx context.Context, src *dagger.Directory) (*dagger.Service, error) {
	db := i.postgresDB()
	i.Db = db

	migration, err := i.startMigrationService(ctx, src)
	if err != nil {
		return nil, err
	}
	defer migration.Stop(ctx)

	return db, nil
}

func (i *CribService) postgresDB() *dagger.Service {
	return dag.Container().
		From("postgres:latest").
		WithEnvVariable("POSTGRES_USER", "postgres").
		WithEnvVariable("POSTGRES_PASSWORD", "example").
		WithEnvVariable("POSTGRES_DB", "cribbage").
		AsService(dagger.ContainerAsServiceOpts{UseEntrypoint: true}).
		WithHostname("db")
}

func (i *CribService) migrationService(src *dagger.Directory) *dagger.Service {
	migration := dag.Container().
		From("golang:latest").
		WithServiceBinding("db", i.Db)

	mSrc := src.Directory("sql")
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
			"/src/sql/migrations",
			"-database",
			"postgres://postgres:example@db:5432/cribbage?sslmode=disable",
			"up",
		}).
		AsService()
}

func (i *CribService) swagger(src *dagger.Directory) *dagger.Service {

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

func (i *CribService) ijhttp(src *dagger.Directory) *dagger.Container {
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
			f = append(f, "/workdir/"+file)
		}
	}

	ij = ij.WithExec(append([]string{"sh", "/ijhttp/ijhttp"}, f...))

	return ij
}

func (i *CribService) startPostgresDB(ctx context.Context) (*dagger.Service, error) {
	db, err := i.postgresDB().Start(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (i *CribService) startSwagger(ctx context.Context, src *dagger.Directory) (*dagger.Service, error) {
	server, err := i.swagger(src).Start(ctx)
	if err != nil {
		return nil, err
	}
	return server, nil
}

func (i *CribService) startMigrationService(ctx context.Context, src *dagger.Directory) (*dagger.Service, error) {
	migration, err := i.migrationService(src).Start(ctx)
	if err != nil {
		return nil, err
	}
	return migration, nil
}

func exclude(c *dagger.Container, dir *dagger.Directory) *dagger.Container {
	return c.WithDirectory("/src", dir, dagger.ContainerWithDirectoryOpts{
		Exclude: []string{
			"./.dagger/internal",
			"./.git",
			"./.dagger/dagger.gen.go",
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
