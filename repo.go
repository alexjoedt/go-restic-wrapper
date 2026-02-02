package restic

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// TODO:
// implement support for S3 and REST

type Repository struct {
	path     string
	password string
}

// Open returns a Repository handle for an existing repository.
// It does not validate connectivity - use Validate() to check if the repository
// is accessible and the password is correct.
//
// For local repositories, use a filesystem path. Future versions will support
// S3 and REST backends.
func Open(path, password string) *Repository {
	return &Repository{
		path:     path,
		password: password,
	}
}

// Validate checks if the repository is accessible and the password is correct.
// It performs a minimal operation (listing snapshots with --last flag) to verify connectivity.
func (r *Repository) Validate(ctx context.Context) error {
	_, err := r.command(ctx, "", "snapshots", "--json", "--last", "--no-lock")
	if err != nil {
		return fmt.Errorf("repository validation failed: %w", err)
	}
	return nil
}

// Init initializes a new restic repository at the specified path.
// Returns an error if the repository already exists (see ErrRepoExists).
func Init(ctx context.Context, path, password string) (*Repository, error) {
	repo := &Repository{
		path:     path,
		password: password,
	}

	return repo.init(ctx)
}

func (r *Repository) init(ctx context.Context) (*Repository, error) {
	_, err := r.command(ctx, "", "init")
	if err != nil {
		return nil, err
	}

	return r, nil
}

// Backup creates a backup of the specified path.
// The path must exist and be accessible. Use options to configure tags, exclusions, etc.
func (r *Repository) Backup(ctx context.Context, path string, options ...BackupOption) (*BackupSummary, error) {
	if path == "" {
		return nil, errors.New("empty path")
	}

	// Check the source to backup
	_, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	var opts backupOptions
	for _, opt := range options {
		opt(&opts)
	}

	args := []string{"backup", "--json"}
	args = append(args, opts.args()...)
	args = append(args, ".")

	out, err := r.command(ctx, path, args...)
	if err != nil {
		return nil, err
	}

	res, err := getSummary(out)
	if err != nil {
		return nil, err
	}

	var summary BackupSummary
	if err := json.Unmarshal(res, &summary); err != nil {
		return nil, fmt.Errorf("failed to parse backup summary: %w", err)
	}

	return &summary, nil
}

// Snapshots returns snapshots from the repository.
// Fetches snapshots in read-only mode (--no-lock).
func (r *Repository) Snapshots(ctx context.Context, filters ...FilterOption) ([]Snapshot, error) {
	var opts filterOptions
	for _, filter := range filters {
		filter(&opts)
	}

	args := []string{"--no-lock", "snapshots", "--json"}
	args = append(args, opts.args()...)

	sn, err := r.command(ctx, "", args...)
	if err != nil {
		return nil, err
	}

	var snapshots []Snapshot
	if err := json.Unmarshal([]byte(sn), &snapshots); err != nil {
		return nil, fmt.Errorf("failed to parse snapshots: %w", err)
	}

	return snapshots, nil
}

// SnapshotById returns the snapshot with the given ID from the repository.
func (r *Repository) SnapshotById(ctx context.Context, id string) (*Snapshot, error) {
	args := []string{"snapshots", "--json", id}

	sn, err := r.command(ctx, "", args...)
	if err != nil {
		return nil, err
	}

	var snapshots []*Snapshot
	if err := json.Unmarshal([]byte(sn), &snapshots); err != nil {
		return nil, fmt.Errorf("failed to parse snapshot: %w", err)
	}

	if len(snapshots) < 1 {
		return nil, fmt.Errorf("no snapshot with id '%s'", id)
	}

	return snapshots[0], nil
}

var (
	idRegex = regexp.MustCompile(`(^latest(:.*)?$|^[0-9a-f]{8}(:.*)?$|^[0-9a-f]{64}(:.*)?$)`)
)

// Restore restores a specific snapshot to the target directory.
// The target directory will be created if it doesn't exist.
func (r *Repository) Restore(ctx context.Context, snapshotID string, target string, options ...RestoreOption) (*RestoreSummary, error) {
	if target == "" {
		return nil, errors.New("no target path")
	}

	if !isPathExists(target) {
		if err := os.MkdirAll(target, 0755); err != nil {
			return nil, err
		}
	}

	if snapshotID == "" {
		return nil, errors.New("empty snapshot id")
	}

	if !isSnapshotID(snapshotID) {
		return nil, ErrInvalidID
	}

	var opts restoreOptions
	for _, opt := range options {
		opt(&opts)
	}

	args := []string{"restore", snapshotID, "--target", target, "--json"}
	args = append(args, opts.args()...)

	out, err := r.command(ctx, "", args...)
	if err != nil {
		return nil, err
	}

	res, err := getSummary(out)
	if err != nil {
		return nil, err
	}

	var summary RestoreSummary
	if err := json.Unmarshal(res, &summary); err != nil {
		return nil, fmt.Errorf("failed to parse restore summary: %w", err)
	}

	return &summary, nil
}

// Forget removes snapshots from the repository based on the specified options.
// At least one option must be specified.
//
// Note: If a snapshot ID is given via ForgetSnapshot(), some filtering options
// (--host, --tag, --path) will be ignored by restic.
// See: https://restic.readthedocs.io/en/stable/060_forget.html#remove-a-single-snapshot
func (r *Repository) Forget(ctx context.Context, options ...ForgetOption) ([]ForgetSummary, error) {
	if len(options) == 0 {
		return nil, errors.New("at least one option must be set")
	}

	var opts forgetOptions
	for _, opt := range options {
		opt(&opts)
	}

	args := []string{"--json", "forget"}
	args = append(args, opts.args()...)

	out, err := r.command(ctx, "", args...)
	if err != nil {
		return nil, err
	}

	data, err := getSummary(out)
	if err != nil {
		return nil, err
	}

	// Note: restic's forget command has limited JSON support.
	// If JSON parsing fails but command succeeded, the operation completed successfully.
	var summary []ForgetSummary
	if err := json.Unmarshal(data, &summary); err != nil {
		// Command succeeded but JSON parsing failed - this is expected for some restic versions
		return nil, fmt.Errorf("forget operation completed but JSON parsing failed (restic version may have limited JSON support): %w", err)
	}

	return summary, nil
}

// Unlock removes locks that other processes created on the repository.
// This is useful for cleaning up stale locks after a crash or interrupted operation.
func (r *Repository) Unlock(ctx context.Context) error {
	args := []string{"unlock", "--remove-all", "--json"}

	_, err := r.command(ctx, "", args...)
	if err != nil {
		return err
	}

	return nil
}

// command wraps the restic command and injects repo and password as environment variables to the process
func (r *Repository) command(ctx context.Context, dir string, args ...string) (string, error) {
	// Check restic binary and version before executing any command
	if err := checkResticVersion(); err != nil {
		return "", fmt.Errorf("restic validation failed: %w", err)
	}

	envArgs := []string{
		"RESTIC_PASSWORD=" + r.password,
		"RESTIC_REPOSITORY=" + r.path,
	}

	home, err := os.UserHomeDir()
	if err == nil {
		envArgs = append(envArgs, "HOME="+home)
	}

	envArgs = append(envArgs, "PATH="+os.Getenv("PATH"))

	// buffers for output
	stdErr := new(bytes.Buffer)
	stdOut := new(bytes.Buffer)

	cmd := exec.CommandContext(ctx, resticBin, args...)

	// set the execute dir
	if dir != "" {
		cmd.Dir = dir
	}

	cmd.Env = envArgs
	cmd.Stdout = stdOut
	cmd.Stderr = stdErr

	// run the command
	if err := cmd.Run(); err != nil {
		return "", parseStdErr(stdErr.String())
	}

	return stdOut.String(), nil
}

// parseStdErr parses the stderr output from the restic command and returns appropriate errors.
func parseStdErr(stdErr string) error {
	switch {
	case strings.Contains(stdErr, "config file already exists"):
		return ErrRepoExists
	case strings.Contains(stdErr, "wrong password") || strings.Contains(stdErr, "invalid password"):
		return ErrInvalidPassword
	case strings.Contains(stdErr, "Is there a repository at the following location?"):
		return ErrRepoNotFound
	case strings.Contains(stdErr, "unable to create lock in backend") ||
		strings.Contains(stdErr, "repository is already locked"):
		return ErrRepoLocked
	case strings.Contains(stdErr, "returned error, retrying after"):
		return ErrInvalidID
	default:
		return fmt.Errorf("restic command failed: %s", stdErr)
	}
}

// isPathExists checks if the path p exists
func isPathExists(p string) bool {
	_, err := os.Stat(p)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false
		}
	}
	return true
}

func isSnapshotID(id string) bool {
	return idRegex.MatchString(id)
}

func getSummary(output string) ([]byte, error) {
	scanner := bufio.NewScanner(strings.NewReader(output))
	var lastSummaryLine []byte

	for scanner.Scan() {
		line := scanner.Bytes()
		// Check for summary message_type
		var msg struct {
			MessageType string `json:"message_type"`
		}
		if err := json.Unmarshal(line, &msg); err == nil {
			if msg.MessageType == "summary" {
				return line, nil
			}
		}
		// Fallback: check for tags field (for forget command)
		if strings.Contains(string(line), `"tags":`) {
			lastSummaryLine = append([]byte{}, line...)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read output: %w", err)
	}

	if lastSummaryLine != nil {
		return lastSummaryLine, nil
	}

	return nil, errors.New("no summary found in output")
}
