package main

import (
	"dagger/CribService/internal/dagger"
	"dagger/CribService/utils"
	"fmt"
)

func (c *CribService) buildServer(src *dagger.Directory) *dagger.Directory {
	gooses := []string{"linux"}
	goarches := []string{"amd64"}

	outputs := dag.Directory()

	for _, goos := range gooses {
		for _, goarch := range goarches {
			// create directory for each OS and architecture
			path := fmt.Sprintf("server/%s/%s/", goos, goarch)

			s := c.golang()
			s = utils.GoMod(src.Directory("server"), s)
			s = s.WithDirectory("/src", src)
			// build artifact
			build := s.
				WithEnvVariable("GOOS", goos).
				WithEnvVariable("GOARCH", goarch).
				WithWorkdir("/src/server").
				WithExec([]string{"go", "build", "-o", path, "main.go"})

			// add build to outputs
			outputs = outputs.WithDirectory(path, build.Directory(path))
		}
	}

	// return build directory
	return outputs
}

func (c *CribService) buildGame(src *dagger.Directory) *dagger.Directory {
	gooses := []string{"linux", "darwin"}
	goarches := []string{"amd64", "arm64"}

	outputs := dag.Directory()

	s := c.golang()
	s = utils.GoMod(src.Directory("cli"), s)
	s = s.WithDirectory("/src", src)

	for _, goos := range gooses {
		for _, goarch := range goarches {
			// create directory for each OS and architecture
			path := fmt.Sprintf("client/%s/%s/", goos, goarch)

			// build artifact
			build := s.
				WithEnvVariable("GOOS", goos).
				WithEnvVariable("GOARCH", goarch).
				WithWorkdir("/src/cli").
				WithExec([]string{"go", "build", "-o", path, "main.go"})

			// add build to outputs
			outputs = outputs.WithDirectory(path, build.Directory(path))
		}
	}

	// return build directory
	return outputs
}

func (c *CribService) buildQueries(src *dagger.Directory) *dagger.Directory {
	g := c.sqlc(src)
	return g.Directory("/src/queries")
}
