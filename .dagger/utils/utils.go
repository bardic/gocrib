package utils

import "dagger/CribService/internal/dagger"

func GoMod(container *dagger.Container, dir *dagger.Directory) *dagger.Container {
	return container.
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
