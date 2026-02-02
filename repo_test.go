package restic

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestIntegration_InitAndValidate(t *testing.T) {
	if err := checkResticVersion(); err != nil {
		t.Skip("restic not available:", err)
	}

	tmpDir, err := os.MkdirTemp("", "restic-test-*")
	if err != nil {
		t.Fatal("failed to create temp dir:", err)
	}
	defer os.RemoveAll(tmpDir)

	ctx := context.Background()
	password := "test-password-123"

	repo, err := Init(ctx, tmpDir, password)
	if err != nil {
		t.Fatal("Init() failed:", err)
	}

	if err := repo.Validate(ctx); err != nil {
		t.Error("Validate() failed:", err)
	}

	_, err = Init(ctx, tmpDir, password)
	if err == nil {
		t.Error("expected error when initializing existing repo")
	}
	if !errors.Is(err, ErrRepoExists) {
		t.Errorf("expected ErrRepoExists, got: %v", err)
	}
}

func TestIntegration_BackupAndSnapshots(t *testing.T) {
	if err := checkResticVersion(); err != nil {
		t.Skip("restic not available:", err)
	}

	tmpRepo, err := os.MkdirTemp("", "restic-test-repo-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpRepo)

	tmpData, err := os.MkdirTemp("", "restic-test-data-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpData)

	testFile := filepath.Join(tmpData, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	password := "test-password-123"

	repo, err := Init(ctx, tmpRepo, password)
	if err != nil {
		t.Fatal("Init() failed:", err)
	}

	summary, err := repo.Backup(ctx, tmpData, WithTags("test", "integration"), WithHost("test-host"))
	if err != nil {
		t.Fatal("Backup() failed:", err)
	}

	if summary.FilesNew == 0 {
		t.Error("expected files to be backed up")
	}
	if summary.SnapshotID == "" {
		t.Error("expected snapshot ID in summary")
	}

	snapshots, err := repo.Snapshots(ctx)
	if err != nil {
		t.Fatal("Snapshots() failed:", err)
	}

	if len(snapshots) != 1 {
		t.Errorf("expected 1 snapshot, got %d", len(snapshots))
	}

	snapshot := snapshots[0]
	if snapshot.Hostname != "test-host" {
		t.Errorf("expected hostname 'test-host', got %s", snapshot.Hostname)
	}
	if len(snapshot.Tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(snapshot.Tags))
	}

	filtered, err := repo.Snapshots(ctx, FilterByTag("test"))
	if err != nil {
		t.Fatal("Snapshots(FilterByTag) failed:", err)
	}
	if len(filtered) != 1 {
		t.Errorf("expected 1 filtered snapshot, got %d", len(filtered))
	}

	snapshotByID, err := repo.SnapshotById(ctx, snapshot.ID.String())
	if err != nil {
		t.Fatal("SnapshotById() failed:", err)
	}
	if snapshotByID.ID.String() != snapshot.ID.String() {
		t.Error("snapshot IDs don't match")
	}
}

func TestIntegration_RestoreSnapshot(t *testing.T) {
	if err := checkResticVersion(); err != nil {
		t.Skip("restic not available:", err)
	}

	tmpRepo, err := os.MkdirTemp("", "restic-test-repo-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpRepo)

	tmpData, err := os.MkdirTemp("", "restic-test-data-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpData)

	tmpRestore, err := os.MkdirTemp("", "restic-test-restore-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpRestore)

	testContent := "test content for restore"
	testFile := filepath.Join(tmpData, "restore-test.txt")
	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	password := "test-password-123"

	repo, err := Init(ctx, tmpRepo, password)
	if err != nil {
		t.Fatal(err)
	}

	summary, err := repo.Backup(ctx, tmpData)
	if err != nil {
		t.Fatal("Backup() failed:", err)
	}

	restoreSummary, err := repo.Restore(ctx, "latest", tmpRestore)
	if err != nil {
		t.Fatal("Restore() failed:", err)
	}

	if restoreSummary.FilesRestored == 0 {
		t.Error("expected files to be restored")
	}

	restoredFile := filepath.Join(tmpRestore, "restore-test.txt")
	content, err := os.ReadFile(restoredFile)
	if err != nil {
		t.Fatal("failed to read restored file:", err)
	}

	if string(content) != testContent {
		t.Errorf("restored content = %q, want %q", string(content), testContent)
	}

	tmpRestore2, err := os.MkdirTemp("", "restic-test-restore2-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpRestore2)

	_, err = repo.Restore(ctx, summary.SnapshotID, tmpRestore2)
	if err != nil {
		t.Error("Restore() with snapshot ID failed:", err)
	}
}

func TestIntegration_ForgetSnapshots(t *testing.T) {
	if err := checkResticVersion(); err != nil {
		t.Skip("restic not available:", err)
	}

	tmpRepo, err := os.MkdirTemp("", "restic-test-repo-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpRepo)

	tmpData, err := os.MkdirTemp("", "restic-test-data-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpData)

	if err := os.WriteFile(filepath.Join(tmpData, "test.txt"), []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	password := "test-password-123"

	repo, err := Init(ctx, tmpRepo, password)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 3; i++ {
		_, err := repo.Backup(ctx, tmpData, WithTags("daily"))
		if err != nil {
			t.Fatal("Backup() failed:", err)
		}
		time.Sleep(100 * time.Millisecond)
	}

	snapshots, err := repo.Snapshots(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(snapshots) != 3 {
		t.Fatalf("expected 3 snapshots, got %d", len(snapshots))
	}

	_, err = repo.Forget(ctx, ForgetKeepLast(2), ForgetByTag("daily"))
	if err != nil {
		t.Fatal("Forget() failed:", err)
	}
}

func TestIntegration_Unlock(t *testing.T) {
	if err := checkResticVersion(); err != nil {
		t.Skip("restic not available:", err)
	}

	tmpRepo, err := os.MkdirTemp("", "restic-test-repo-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpRepo)

	ctx := context.Background()
	password := "test-password-123"

	repo, err := Init(ctx, tmpRepo, password)
	if err != nil {
		t.Fatal(err)
	}

	err = repo.Unlock(ctx)
	if err != nil {
		t.Error("Unlock() failed:", err)
	}
}

func TestIntegration_InvalidPassword(t *testing.T) {
	if err := checkResticVersion(); err != nil {
		t.Skip("restic not available:", err)
	}

	tmpRepo, err := os.MkdirTemp("", "restic-test-repo-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpRepo)

	ctx := context.Background()
	correctPassword := "correct-password"
	wrongPassword := "wrong-password"

	_, err = Init(ctx, tmpRepo, correctPassword)
	if err != nil {
		t.Fatal(err)
	}

	repo := Open(tmpRepo, wrongPassword)
	_, err = repo.Snapshots(ctx)
	if err == nil {
		t.Error("expected error with wrong password")
	}
	if !errors.Is(err, ErrInvalidPassword) {
		t.Errorf("expected ErrInvalidPassword, got: %v", err)
	}
}

func TestIntegration_RepositoryNotFound(t *testing.T) {
	if err := checkResticVersion(); err != nil {
		t.Skip("restic not available:", err)
	}

	ctx := context.Background()
	repo := Open("/nonexistent/path/to/repo", "password")

	_, err := repo.Snapshots(ctx)
	if err == nil {
		t.Error("expected error for non-existent repository")
	}
	if !errors.Is(err, ErrRepoNotFound) {
		t.Errorf("expected ErrRepoNotFound, got: %v", err)
	}
}

func TestIntegration_InvalidSnapshotID(t *testing.T) {
	if err := checkResticVersion(); err != nil {
		t.Skip("restic not available:", err)
	}

	tmpRepo, err := os.MkdirTemp("", "restic-test-repo-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpRepo)

	tmpRestore, err := os.MkdirTemp("", "restic-test-restore-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpRestore)

	ctx := context.Background()
	repo, err := Init(ctx, tmpRepo, "password")
	if err != nil {
		t.Fatal(err)
	}

	_, err = repo.Restore(ctx, "invalid-id", tmpRestore)
	if err == nil {
		t.Error("expected error for invalid snapshot ID")
	}
	if !errors.Is(err, ErrInvalidID) {
		t.Errorf("expected ErrInvalidID, got: %v", err)
	}
}
