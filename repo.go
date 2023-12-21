package restic

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/alexjoedt/go-restic-wrapper/backup"
	"github.com/alexjoedt/go-restic-wrapper/filter"
	"github.com/alexjoedt/go-restic-wrapper/forget"
	"github.com/alexjoedt/go-restic-wrapper/restore"
)

// TODO:
// implement support for S3 and Rest

type Repository struct {
	path     string
	password string
}

// Connect creates a new instance of a exiting restic repository.
func Connect(ctx context.Context, repoPath string, password string) (*Repository, error) {

	repo := &Repository{
		path:     repoPath,
		password: password,
	}

	_, err := repo.Snapshots(ctx)
	if err != nil {
		return nil, errors.New("failed to connect to restic repo")
	}

	return repo, nil
}

// Init initialize a new restic repository
func Init(ctx context.Context, repoPath string, password string) (*Repository, error) {
	repo := &Repository{
		path:     repoPath,
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

// Backup backing up the given path
func (r *Repository) Backup(ctx context.Context, path string, options ...backup.OptionFunc) (*BackupSummary, error) {

	// Check the path
	if path == "" {
		return nil, errors.New("empty path")
	}

	// Check the source to backup
	_, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	args := []string{"backup", "--json"}
	args = append(args, backup.Args(options...)...)
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
	err = json.Unmarshal(res, &summary)
	if err != nil {
		return nil, nil
	}

	return &summary, nil
}

// Snapshots returns snapshots from the repository.
// Fetches Snapshots in read only mode (--no-lock)
func (r *Repository) Snapshots(ctx context.Context, filters ...filter.OptionFunc) ([]Snapshot, error) {

	args := []string{"snapshots", "--json"}
	args = append(args, filter.Args(filters...)...)

	sn, err := r.command(ctx, "", args...)
	if err != nil {
		return nil, err
	}

	var snapshots []Snapshot
	err = json.Unmarshal([]byte(sn), &snapshots)
	if err != nil {
		return nil, err
	}

	return snapshots, nil
}

// SnapshotById returns the snapshot with given id from the repository
func (r *Repository) SnapshotById(ctx context.Context, id string) (*Snapshot, error) {

	args := []string{"snapshots", "--json"}
	args = append(args, id)

	sn, err := r.command(ctx, "", args...)
	if err != nil {
		return nil, err
	}

	var snapshots []*Snapshot
	err = json.Unmarshal([]byte(sn), &snapshots)
	if err != nil {
		return nil, err
	}

	if len(snapshots) < 1 {
		return nil, fmt.Errorf("no snapshot wiht id '%s'", id)
	}

	return snapshots[0], nil
}

var (
	idRegex regexp.Regexp = *regexp.MustCompile(`(^latest(:.*)?$|^[0-9a-f]{8}(:.*)?$|^[0-9a-f]{64}(:.*)?$)`)
)

// Restore restores a specific snapshot
func (r *Repository) Restore(ctx context.Context, snapshotID string, target string, options ...restore.OptionFunc) (*RestoreSummary, error) {
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
		return nil, errors.New("invalid snapshot ID")
	}

	args := []string{"restore", snapshotID, "--target", target, "--json"}

	args = append(args, restore.Args(options...)...)
	out, err := r.command(ctx, "", args...)
	if err != nil {
		return nil, err
	}

	res, err := getSummary(out)
	if err != nil {
		return nil, err
	}

	var summary RestoreSummary
	err = json.Unmarshal(res, &summary)
	if err != nil {
		return nil, nil
	}

	return &summary, nil
}

// Forget forgets a snapshot.
// If a snapshot ID is given, some option will be ignored by restic.
// E.g. --host, --tag and --path. See documentation: https://restic.readthedocs.io/en/stable/060_forget.html#remove-a-single-snapshot
func (r *Repository) Forget(ctx context.Context, options ...forget.OptionFunc) ([]ForgetSummary, error) {

	if len(options) == 0 {
		return nil, errors.New("at least one option must be set")
	}

	args := []string{
		"--json", // json output seems not supported yet, so there is no output with exit 0
		"forget",
	}

	args = append(args, forget.Args(options...)...)
	out, err := r.command(ctx, "", args...)
	if err != nil {
		return nil, err
	}

	data, err := getSummary(out)
	if err != nil {
		return nil, err
	}

	var summary []ForgetSummary
	err = json.Unmarshal(data, &summary)
	if err != nil {
		// as long --json is not supported on forget, we return nil, nil
		return nil, nil
	}

	return summary, nil
}

// Unlock remove locks other processes created on the repository
func (r *Repository) Unlock(ctx context.Context) error {

	args := []string{"unlock", "--json"}

	_, err := r.command(ctx, "", args...)
	if err != nil {
		return err
	}

	return nil
}

// command wraps the restic command and injects repo and password as environment variables to the process
func (r *Repository) command(ctx context.Context, dir string, args ...string) (string, error) {

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

var (
	ErrRepoAlreadyExist error = errors.New("restic repo already exist, use restic.Connect")
	ErrInvalidID        error = errors.New("invalid snapshot ID")
)

// parseStdErr parses the stderr output from the restic command
func parseStdErr(stdErr string) error {
	switch {
	case strings.Contains(stdErr, "failed: config file already exists"):
		return ErrRepoAlreadyExist
	}

	return errors.New(stdErr)
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
	reader := bufio.NewReader(strings.NewReader(output))
	res := make([]byte, 0)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, errors.New("failed to read output")
		}

		if strings.Contains(string(line), "summary") || strings.Contains(string(line), `"tags":`) {
			res = line
		}
	}
	return res, nil
}
