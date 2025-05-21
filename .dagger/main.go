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

func (c *CribService) BuildQueries(ctx context.Context, src *dagger.Directory) *dagger.Directory {
	return c.buildQueries(src)
}

func (c *CribService) BuildServer(src *dagger.Directory) *dagger.Directory {
	return c.buildServer(src)
}

func (c *CribService) BuildGame(src *dagger.Directory) *dagger.Directory {
	return c.buildGame(src)
}

func (c *CribService) DbUp(ctx context.Context, src *dagger.Directory, withPort bool) (*dagger.Service, error) {
	p, err := c.postgresService(withPort).Start(ctx)
	// _, err := c.migrationService(src, p).Start(ctx)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (i *CribService) Server(ctx context.Context, src *dagger.Directory, migrate bool) (*dagger.Service, error) {
	return i.serverService(ctx, src, migrate)
}

func (i *CribService) ServerTest(ctx context.Context, src *dagger.Directory) (*dagger.Container, error) {
	return i.gameServer(src), nil
}

func (i *CribService) TestHttp(ctx context.Context, src *dagger.Directory,
	// +optional
	isCI bool,
) (string, error) {
	httpDir := src.Directory("http")
	ij := i.http(httpDir)

	e := "local"

	if isCI {
		e = "ci"

		s, err := i.serverService(ctx, src, true)
		if err != nil {
			return "", err
		}

		i.ServerService = s
		ij = ij.WithServiceBinding("server", s)
	}

	entries, err := httpDir.Entries(context.Background())
	if err != nil {
		return "", nil
	}

	f := make([]string, 0)
	for _, file := range entries {
		if strings.HasSuffix(file, ".http") {
			f = append(f, "/workdir/"+file)
		}
	}

	ij = ij.WithExec(append([]string{"sh", "/ijhttp/ijhttp", "-v", "/workdir/http-client.env.json", "-e", e}, f...))

	return ij.Stdout(ctx)
}
