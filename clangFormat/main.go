package main

import (
	"context"
)

type ClangFormat struct {}

func (m *ClangFormat) MyFunction(ctx context.Context, stringArg string) (*Container, error) {
	return dag.Container().From("alpine:latest").WithExec([]string{"echo", stringArg}).Sync(ctx)
}
