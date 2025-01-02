package utils

import "dagger/CribService/internal/dagger"

func GoMod(src *dagger.Directory, container *dagger.Container) *dagger.Container {
	return container.
		WithFiles("/src", []*dagger.File{src.File("go.mod"), src.File("go.sum")}).
		WithWorkdir("/src").
		WithExec([]string{"go", "mod", "download"})
}

func ContainerExclude(c *dagger.Container, dir *dagger.Directory) *dagger.Container {
	return c.WithDirectory("/src", dir, dagger.ContainerWithDirectoryOpts{
		Exclude: []string{
			"./.dagger/internal",
			"./.git",
			"./.dagger/dagger.gen.go",
		},
	})
}
