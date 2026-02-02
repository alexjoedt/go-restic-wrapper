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






















































































}	fmt.Printf("Processed %d snapshot groups\n", len(summaries))	}		log.Fatal(err)	if err != nil {	)		restic.ForgetWithPrune(),		restic.ForgetKeepLast(7),	summaries, err := repo.Forget(ctx,	// Forget old snapshots, keep last 7	repo := restic.Open("/path/to/repo", "my-secret-password")	ctx := context.Background()func ExampleRepository_Forget() {}	fmt.Printf("Restored %d files\n", summary.FilesRestored)	}		log.Fatal(err)	if err != nil {	)		restic.RestoreExclude("*.tmp"),	summary, err := repo.Restore(ctx, "abc12345", "/path/to/restore",	// Restore a snapshot by ID	repo := restic.Open("/path/to/repo", "my-secret-password")	ctx := context.Background()func ExampleRepository_Restore() {}	fmt.Printf("Found %d filtered snapshots\n", len(filtered))	}		log.Fatal(err)	if err != nil {	)		restic.FilterLatest(5),		restic.FilterByHost("myserver"),		restic.FilterByTag("daily"),	filtered, err := repo.Snapshots(ctx,	// Filter snapshots by tag and host	fmt.Printf("Found %d snapshots\n", len(snapshots))	}		log.Fatal(err)	if err != nil {	snapshots, err := repo.Snapshots(ctx)	// List all snapshots	repo := restic.Open("/path/to/repo", "my-secret-password")	ctx := context.Background()func ExampleRepository_Snapshots() {}	fmt.Printf("Backed up %d files\n", summary.FilesNew)	}		log.Fatal(err)	if err != nil {	)		restic.WithExclude("*.tmp", "*.log"),		restic.WithTags("daily", "important"),	summary, err := repo.Backup(ctx, "/path/to/backup",	// Backup with tags and exclusions	repo := restic.Open("/path/to/repo", "my-secret-password")	ctx := context.Background()func ExampleRepository_Backup() {}	fmt.Println("Repository opened successfully")	}		log.Fatal(err)	if err := repo.Validate(ctx); err != nil {	// Optionally validate connectivity	repo := restic.Open("/path/to/repo", "my-secret-password")	// Open an existing repository	ctx := context.Background()func ExampleOpen() {CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
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
