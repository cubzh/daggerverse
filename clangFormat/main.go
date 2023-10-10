package main

import (
	"context"
)

type ClangFormat struct{}

type CheckFormatOpts struct {
	dir *Directory
}

func (m *ClangFormat) CheckFormat(ctx context.Context, opts CheckFormatOpts) (*Container, error) {

	// get Docker container from hub
	c := dag.Container().From("gaetan/clang-tools")

	// mount directory into container and make it the working directory
	c = c.WithMountedDirectory("/workdir", opts.dir)
	c = c.WithWorkdir("/workdir")

	// --dry-run: do not apply changes
	command := "set -e ; set -o pipefail ; find . -maxdepth 1 -regex '^.*\\.\\(cpp\\|hpp\\|c\\|h\\)$' -print0 | xargs -0 clang-format --dry-run --Werror -style=file"
	c = c.WithExec([]string{"ash", "-c", command})

	return c.Sync(ctx)
}
