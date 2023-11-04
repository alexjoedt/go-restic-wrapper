package main

import (
	"context"
	"fmt"
	"log"

	"github.com/alexjoedt/go-restic-wrapper"
)

func main() {

	repo, err := restic.Connect(ctx, "/path/to/local-repo", "password")
	if err != nil {
		log.Fatal(err)
	}

	// Backup data from path /path/to/backup with defaults
	sum, err := repo.Backup(context.Background(), "/path/to/backup")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v", sum)

}
