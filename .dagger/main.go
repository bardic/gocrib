package main

import (
	"context"
	"dagger/CribService/internal/dagger"
	"os"
	"strings"
)

type CribService struct {
	ServerContainer *dagger.Container
	Db              *dagger.Service
	ServerService   *dagger.Service
	// Migrator         *dagger.Container
	// MigrationService *dagger.Service
}

func (c *CribService) Gen(ctx context.Context, src *dagger.Directory) *dagger.Directory {
	g := c.sqlc(src)
	return g.Directory("/src/queries")
}
func (c *CribService) Postgres(ctx context.Context, src *dagger.Directory, withPort bool) *dagger.Service {
	os.Setenv("GOCRIB_HOST", "localhost")
	return c.postgresService(withPort)
}

func (i *CribService) Server(ctx context.Context, src *dagger.Directory, migrate bool) *dagger.Service {
	os.Setenv("GOCRIB_HOST", "")
	return i.startGameServer(src, migrate)
}

func (i *CribService) TestHttp(ctx context.Context, src *dagger.Directory) (string, error) {
	s := i.startGameServer(src, true)

	ij := i.http(ctx, src)

	ij = ij.WithServiceBinding("server", s)

	i.ServerService = s

	entries, err := src.Directory("http").Entries(context.Background())

	if err != nil {
		return "", nil
	}

	f := make([]string, 0)
	for _, file := range entries {
		if strings.HasSuffix(file, ".http") {
			f = append(f, "/workdir/"+file)
		}
	}

	ij = ij.WithExec(append([]string{"sh", "/ijhttp/ijhttp"}, f...))

	return ij.Stdout(ctx)
}
