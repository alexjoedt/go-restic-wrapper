package main

import (
	"context"
	"fmt"
	"log"

	"github.com/alexjoedt/go-restic/restic"
	"github.com/alexjoedt/go-restic/restic/forget"
)

func main() {

	repo, err := restic.Connect(context.Background(), "/path/to/local-repo", "password")
	if err != nil {
		log.Fatal(err)
	}

	// Forget the latest snapshots with tag "home"
	// Keep last 4 snapshots
	// prune
	sum, err := repo.Forget(context.Background(),
		forget.WithTags("home"),
		forget.WithKeepLast(4),
		forget.WithPrune(),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v", sum)

}
