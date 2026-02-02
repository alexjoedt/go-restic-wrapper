package main

import (
	"context"
	"fmt"
	"os"

	"github.com/alexjoedt/go-restic-wrapper"
)

const (
	testPath = "/Users/alex/workspace/github.com/alexjoedt/go-restic-wrapper/testdata"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize a new repository (use Open() for existing repos)
	repo, err := restic.Init(ctx, testPath, "1234")
	if err != nil {
		return err
	}

	// Validate connectivity (optional)
	if err := repo.Validate(ctx); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	return nil
}
