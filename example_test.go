/*
BSD 2-Clause License

Copyright (c) 2014, Alexander Neumann <alexander@bumpern.de>
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

* Redistributions of source code must retain the above copyright notice, this
  list of conditions and the following disclaimer.

* Redistributions in binary form must reproduce the above copyright notice,
  this list of conditions and the following disclaimer in the documentation
  and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER






















































































func ExampleOpen() {
	ctx := context.Background()

	// Open an existing repository
	repo := restic.Open("/path/to/repo", "my-secret-password")

	// Optionally validate connectivity
	if err := repo.Validate(ctx); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Repository opened successfully")
}

func ExampleRepository_Backup() {
	ctx := context.Background()
	repo := restic.Open("/path/to/repo", "my-secret-password")

	// Backup with tags and exclusions
	summary, err := repo.Backup(ctx, "/path/to/backup",
		restic.WithTags("daily", "important"),
		restic.WithExclude("*.tmp", "*.log"),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Backed up %d files\n", summary.FilesNew)
}

func ExampleRepository_Snapshots() {
	ctx := context.Background()
	repo := restic.Open("/path/to/repo", "my-secret-password")

	// List all snapshots
	snapshots, err := repo.Snapshots(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found %d snapshots\n", len(snapshots))

	// Filter snapshots by tag and host
	filtered, err := repo.Snapshots(ctx,
		restic.FilterByTag("daily"),
		restic.FilterByHost("myserver"),
		restic.FilterLatest(5),
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found %d filtered snapshots\n", len(filtered))
}

func ExampleRepository_Restore() {
	ctx := context.Background()
	repo := restic.Open("/path/to/repo", "my-secret-password")

	// Restore a snapshot by ID
	summary, err := repo.Restore(ctx, "abc12345", "/path/to/restore",
		restic.RestoreExclude("*.tmp"),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Restored %d files\n", summary.FilesRestored)
}

func ExampleRepository_Forget() {
	ctx := context.Background()
	repo := restic.Open("/path/to/repo", "my-secret-password")

	// Forget old snapshots, keep last 7
	summaries, err := repo.Forget(ctx,
		restic.ForgetKeepLast(7),
		restic.ForgetWithPrune(),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Processed %d snapshot groups\n", len(summaries))
}CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

package restic_test

import (
	"context"
	"fmt"
	"log"

	"github.com/alexjoedt/go-restic-wrapper"
)

func ExampleInit() {
	ctx := context.Background()

	// Initialize a new repository
	repo, err := restic.Init(ctx, "/path/to/repo", "my-secret-password")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Repository initialized")
	_ = repo
}
