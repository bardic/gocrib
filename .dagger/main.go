package main

import (
	"context"
	"dagger/CribService/internal/dagger"
	"strings"
)

type CribService struct {
	ServerContainer *dagger.Container
	Db              *dagger.Service
	ServerService   *dagger.Service
}

func (c *CribService) Gen(ctx context.Context, src *dagger.Directory) *dagger.Directory {
	g := c.sqlc(src)
	return g.Directory("/src/queries")
}
func (c *CribService) DbUp(ctx context.Context, src *dagger.Directory, withPort bool) (*dagger.Service, error) {
	p := c.postgresService(withPort)
	_, err := c.migrationService(src, p).Start(ctx)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (i *CribService) Server(ctx context.Context, src *dagger.Directory, migrate bool) (*dagger.Service, error) {
	return i.gameServerService(ctx, src, migrate)
}

func (i *CribService) TestHttp(ctx context.Context, src *dagger.Directory) (string, error) {
	s, err := i.gameServerService(ctx, src, true)

	if err != nil {
		return "", err
	}

	ij := i.http(src)

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

	ij = ij.WithExec(append([]string{"sh", "/ijhttp/ijhttp", "-e", "production"}, f...))

	return ij.Stdout(ctx)
}
