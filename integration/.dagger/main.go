// A generated module for Integration functions
//
// This module has been generated via dagger init and serves as a reference to
// basic module structure as you get started with Dagger.
//
// Two functions have been pre-created. You can modify, delete, or add to them,
// as needed. They demonstrate usage of arguments and return types using simple
// echo and grep commands. The functions can be called from the dagger CLI or
// from one of the SDKs.
//
// The first line in this comment block is a short description line and the
// rest is a long description with more detail on the module's purpose or usage,
// if appropriate. All modules should have a short description.

package main

import (
	"context"
	"dagger/integration/internal/dagger"
)

type Integration struct{}

// dbService := dag.Container().
// 	From("postgres:latest").
// 	WithEnvVariable("POSTGRES_USER", "postgres").
// 	WithEnvVariable("POSTGRES_PASSWORD", "example").
// 	WithEnvVariable("POSTGRES_DB", "cribbage").
// 	AsService()

func (i *Integration) Test(server, test *dagger.Directory) (string, error) {

	dbService := dag.Container().
		From("postgres:latest").
		WithEnvVariable("POSTGRES_USER", "postgres").
		WithEnvVariable("POSTGRES_PASSWORD", "example").
		WithEnvVariable("POSTGRES_DB", "cribbage").
		WithExposedPort(5432).
		AsService()

	migrateService := dag.Container().
		From("golang:latest").
		WithServiceBinding("db", dbService).
		WithDirectory("/src", server).
		WithExec([]string{"go", "install", "-tags", "'postgres'", "github.com/golang-migrate/migrate/v4/cmd/migrate@latest"}).
		WithExec([]string{"migrate", "-path", "/src/migrations", "-database", "postgres://postgres:example@db:5432/cribbage?sslmode=disable", "up"}).
		AsService()

	serverService := dag.Container().
		From("golang:latest").
		WithDirectory("/src", server, dagger.ContainerWithDirectoryOpts{
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
		}).
		WithExposedPort(1323).
		WithWorkdir("/src").
		WithExec([]string{"go", "run", "/src/server/main.go"}).
		AsService()

	return dag.Container().
		From("alpine:latest").
		WithServiceBinding("migrate", migrateService).
		WithServiceBinding("server", serverService).
		WithDirectory("/workdir", test).
		WithExec([]string{"apk", "add", "openjdk17-jdk", "curl", "unzip"}).
		WithExec([]string{"/bin/sh", "-c", "mkdir /ijhttp"}).
		WithExec([]string{"curl", "-f", "-L", "-o", "/ijhttp/ijhttp.zip", "https://jb.gg/ijhttp/latest"}).
		WithExec([]string{"unzip", "/ijhttp/ijhttp.zip"}).
		WithExec([]string{"/bin/sh", "-c", "chmod +x /ijhttp/ijhttp"}).
		WithExec([]string{"sh", "/ijhttp/ijhttp", "/workdir/test.http"}).
		Stdout(context.Background())
}
