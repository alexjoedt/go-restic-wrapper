# Copilot Instructions for go-restic-wrapper

## Project Overview
This is a Go wrapper library for [restic](https://restic.net/), providing a type-safe API for backup operations. The library shells out to the `restic` CLI binary and parses JSON output into Go structs.

**Critical**: Restic binary v0.16.0+ must be installed on the system. The package validates this in `init()` and exits fatally if missing.

## Architecture

### Core Components
- **[restic.go](restic.go)**: Package initialization with restic binary validation
- **[repo.go](repo.go)**: `Repository` struct and all operations (Backup, Restore, Forget, Snapshots)
- **[types.go](types.go)**: JSON response structs (`BackupSummary`, `RestoreSummary`, `ForgetSummary`)
- **[snapshots.go](snapshots.go)**: `Snapshot` and `ID` types (BSD 2-Clause licensed from restic project)
- **Option packages** ([backup/](backup/), [filter/](filter/), [forget/](forget/), [restore/](restore/)): Functional options for each operation

### Key Design Patterns

#### 1. Command Execution Pattern
All restic operations use the private `command()` method in [repo.go](repo.go#L251):
```go
// Sets RESTIC_PASSWORD, RESTIC_REPOSITORY as env vars
// Changes working directory if needed (for backup operations)
// Captures stdout/stderr separately
func (r *Repository) command(ctx context.Context, dir string, args ...string) (string, error)
```

#### 2. Functional Options Pattern
Each operation uses functional options to build CLI arguments:
```go
// In backup/backup.go
func WithTags(tags ...string) OptionFunc
func WithExcludes(excludes ...string) OptionFunc

// Usage:
repo.Backup(ctx, path, backup.WithTags("daily"), backup.WithExcludes("*.tmp"))
```

Options are converted to CLI args via private `args()` methods that return `[]string`.

#### 3. JSON Output Parsing
All commands use `--json` flag. Helper function `getSummary()` in [repo.go](repo.go#L324) extracts the final summary line containing `"summary"` or `"tags"` from multi-line JSON output.

## Critical Implementation Details

### Repository Initialization
- `Init()` creates a new repository - fails if repo already exists (`ErrRepoAlreadyExist`)
- `Connect()` connects to existing repository - validates by calling `Snapshots()`

### Backup Operation Quirks
- Changes working directory to the backup source path and uses `.` as target
- This is why `Backup()` takes a `dir` parameter to `command()`
- Always returns `BackupSummary` with stats about files/bytes processed

### Snapshot ID Validation
Snapshot IDs must match regex in [repo.go](repo.go#L149):
- `latest` or `latest:<n>` 
- Short ID: 8 hex characters (e.g., `a1b2c3d4`)
- Full ID: 64 hex characters
- All forms support optional `:<path>` suffix

Validated by `isSnapshotID()` before restore operations.

### Error Handling Strategy
- Use `parseStdErr()` to convert restic stderr messages into typed errors
- Defined errors: `ErrRepoAlreadyExist`, `ErrInvalidID`, `ErrRepoLocked`
- Fatal errors in `init()` print to stdout and `os.Exit(1)` - this is intentional

### Locking Behavior
- `Snapshots()` uses `--no-lock` flag for read-only access
- Write operations can hit `ErrRepoLocked` - use `Unlock()` to clear stale locks
- `Unlock()` always uses `--remove-all` flag

## Testing & Development

### Local Testing
Uses [testdata/](testdata/) directory as a local restic repository for testing:
```go
const testPath = "/Users/alex/workspace/github.com/alexjoedt/go-restic-wrapper/testdata"
```

See [_examples/full/main.go](_examples/full/main.go) for initialization patterns.

### No Unit Tests
Currently no `*_test.go` files exist. Testing requires actual restic binary and repository setup.

## Adding New Operations

When adding a new restic command:

1. **Create options package** (e.g., `check/check.go`):
   ```go
   type OptionFunc func(*options)
   func Args(opts ...OptionFunc) []string
   ```

2. **Add method to Repository** in [repo.go](repo.go):
   ```go
   func (r *Repository) Check(ctx context.Context, options ...check.OptionFunc) error {
       args := []string{"check", "--json"}
       args = append(args, check.Args(options...)...)
       out, err := r.command(ctx, "", args...)
       // Parse JSON output...
   }
   ```

3. **Add response type** to [types.go](types.go) if restic returns structured data

4. **Handle stderr errors** in `parseStdErr()` if restic has specific error messages

## Known Limitations

- **No S3/REST support** - only local filesystem repositories (marked as TODO in [repo.go](repo.go#L22))
- **Forget JSON output** - restic doesn't fully support `--json` for forget, so returns `nil, nil` on success ([repo.go](repo.go#L218))
- **No streaming progress** - only final summary is returned, not real-time progress updates
- **API stability** - README warns API may change; this is a development-phase library

## Dependencies
- `github.com/hashicorp/go-version` - for restic version validation only
- Go 1.20+ (specified in [go.mod](go.mod))

## Context Propagation
All public methods accept `context.Context` as first parameter. Pass this to `exec.CommandContext` for cancellation support.
