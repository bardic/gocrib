package main

import (
	"context"
	"dagger/CribService/internal/dagger"
	"os"
)

type CribService struct {
	Db *dagger.Service
}

func (c *CribService) Gen(ctx context.Context, src *dagger.Directory) *dagger.Directory {
	g := c.sqlc(src)
	return g.Directory("/src/queries")
}

func (i *CribService) Server(ctx context.Context, src *dagger.Directory) *dagger.Service {
	os.Setenv("GOCRIB_HOST", "")
	return i.startGameServer(src)
}

func (i *CribService) TestHttp(ctx context.Context, src *dagger.Directory) (string, error) {
	ij := i.http(ctx, src)
	return ij.Stdout(ctx)
}
