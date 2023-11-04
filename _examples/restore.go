package main

import (
	"context"
	"fmt"
	"log"

	"github.com/alexjoedt/go-restic-wrapper/restic"
)

func main() {

	repo, err := restic.Connect(ctx, "/path/to/local-repo", "password")
	if err != nil {
		log.Fatal(err)
	}

	// Restore the latest snapshot to /path/to/restore
	sum, err := repo.Restore(context.Background(), "latest", "/path/to/restore")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v", sum)

}
