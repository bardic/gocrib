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
	"dagger/integration/internal/dagger"
)

type Integration struct{}

func (i *Integration) Test(migrations, server, test *dagger.Directory) *dagger.Container {
	dbService := dag.Container().
		From("postgres:latest").
		WithEnvVariable("POSTGRES_USER", "postgres").
		WithEnvVariable("POSTGRES_PASSWORD", "example").
		WithEnvVariable("POSTGRES_DB", "cribbage").
		WithExposedPort(5432).
		AsService()

	dag.Container().
		From("migrate/migrate").
		WithDirectory("/migrations", server).
		WithServiceBinding("db", dbService).
		WithExec([]string{"migrate", "-database", "postgres://postgres:example@db:5432/cribbage?sslmode=disable", "-path", "migrations", "up"})

	serverService := dag.Container().
		From("golang:latest").
		WithDirectory("/src", server).
		WithWorkdir("/src").
		WithExec([]string{"go", "run", "main.go"}).
		AsService()

	return dag.Container().
		From("jetbrains/intellij-http-client").
		WithDirectory("/workdir", test).
		WithServiceBinding("db", dbService).
		WithServiceBinding("server", serverService).
		WithExec([]string{"run.http"})
}
