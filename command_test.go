package restic

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"
)

// TestCommandHook verifies the command hook is called correctly
func TestCommandHook(t *testing.T) {
	var capturedArgs []string
	var hookCalled bool

	// Set up hook
	SetCommandHook(func(ctx context.Context, args []string) {
		hookCalled = true
		capturedArgs = args
	})
	defer SetCommandHook(nil) // Clean up

	// This will fail but that's OK - we just want to verify the hook was called
	repo := Open("/tmp/nonexistent-test-repo", "test-password")
	ctx := context.Background()
	_, _ = repo.Snapshots(ctx)

	if !hookCalled {
		t.Error("command hook was not called")
	}

	// Verify args contain expected values
	argsStr := strings.Join(capturedArgs, " ")
	if !strings.Contains(argsStr, "snapshots") {
		t.Errorf("expected 'snapshots' in args, got: %v", capturedArgs)
	}
	if !strings.Contains(argsStr, "--json") {
		t.Errorf("expected '--json' in args, got: %v", capturedArgs)
	}
}

// TestCommandValidation verifies input validation
func TestCommandValidation(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		password string
		wantErr  string
	}{
		{
			name:     "empty path",
			path:     "",
			password: "password",
			wantErr:  "repository path is empty",
		},
		{
			name:     "empty password",
			path:     "/some/path",
			password: "",
			wantErr:  "repository password is empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := Open(tt.path, tt.password)
			ctx := context.Background()
			_, err := repo.Snapshots(ctx)

			if err == nil {
				t.Error("expected error, got nil")
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("expected error containing %q, got: %v", tt.wantErr, err)
			}
		})
	}
}

// TestContextCancellation verifies context cancellation is handled correctly
func TestContextCancellation(t *testing.T) {
	// Skip if restic not available
	if err := checkResticVersion(); err != nil {
		t.Skip("restic not available")
	}

	repo := Open("/tmp/test-repo", "test-password")

	// Create a context that's already cancelled
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := repo.Snapshots(ctx)
	if err == nil {
		t.Fatal("expected error for cancelled context")
	}

	// Should mention cancellation
	if !strings.Contains(err.Error(), "cancel") {
		t.Errorf("expected error to mention cancellation, got: %v", err)
	}
}

// TestContextTimeout verifies timeout handling
func TestContextTimeout(t *testing.T) {
	// Skip if restic not available
	if err := checkResticVersion(); err != nil {
		t.Skip("restic not available")
	}

	repo := Open("/tmp/test-repo", "test-password")

	// Create a context with very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	time.Sleep(10 * time.Millisecond) // Ensure context expires

	_, err := repo.Snapshots(ctx)
	if err == nil {
		t.Fatal("expected error for timed out context")
	}

	// Should mention deadline or cancellation
	errStr := err.Error()
	if !strings.Contains(errStr, "deadline") && !strings.Contains(errStr, "cancel") {
		t.Errorf("expected error to mention timeout/cancellation, got: %v", err)
	}
}

// TestLimitedBuffer verifies buffer size limits
func TestLimitedBuffer(t *testing.T) {
	tests := []struct {
		name      string
		limit     int
		data      string
		wantError bool
	}{
		{
			name:      "within limit",
			limit:     100,
			data:      "hello world",
			wantError: false,
		},
		{
			name:      "exceeds limit",
			limit:     10,
			data:      "this is a very long string that exceeds the limit",
			wantError: true,
		},
		{
			name:      "exact limit",
			limit:     5,
			data:      "12345",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lb := &limitedBuffer{
				buf:   new(bytes.Buffer),
				limit: tt.limit,
			}

			_, err := lb.Write([]byte(tt.data))

			if tt.wantError && err == nil {
				t.Error("expected error for data exceeding limit, got nil")
			}
			if !tt.wantError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

// TestBackupOptions_args verifies backup option argument generation
func TestBackupOptions_args(t *testing.T) {
	tests := []struct {
		name string
		opts []BackupOption
		want []string
	}{
		{
			name: "no options",
			opts: nil,
			want: []string{},
		},
		{
			name: "with host",
			opts: []BackupOption{WithHost("myhost")},
			want: []string{"--host", "myhost"},
		},
		{
			name: "with multiple tags",
			opts: []BackupOption{WithTags("daily", "important")},
			want: []string{"--tag", "daily", "--tag", "important"},
		},
		{
			name: "with excludes",
			opts: []BackupOption{WithExclude("*.tmp", "*.log")},
			want: []string{"--exclude", "*.tmp", "--exclude", "*.log"},
		},
		{
			name: "with includes",
			opts: []BackupOption{WithInclude("*.go", "*.md")},
			want: []string{"--include", "*.go", "--include", "*.md"},
		},
		{
			name: "combined options",
			opts: []BackupOption{
				WithHost("server1"),
				WithTags("daily", "prod"),
				WithExclude("*.tmp"),
				WithInclude("*.go"),
			},
			want: []string{
				"--host", "server1",
				"--tag", "daily", "--tag", "prod",
				"--exclude", "*.tmp",
				"--include", "*.go",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var opts backupOptions
			for _, opt := range tt.opts {
				opt(&opts)
			}
			got := opts.args()
			if !stringSliceEqual(got, tt.want) {
				t.Errorf("backupOptions.args() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestFilterOptions_args verifies filter option argument generation
func TestFilterOptions_args(t *testing.T) {
	tests := []struct {
		name string
		opts []FilterOption
		want []string
	}{
		{
			name: "no options",
			opts: nil,
			want: []string{},
		},
		{
			name: "with hosts",
			opts: []FilterOption{FilterByHost("host1", "host2")},
			want: []string{"--host", "host1", "--host", "host2"},
		},
		{
			name: "with paths",
			opts: []FilterOption{FilterByPath("/home", "/etc")},
			want: []string{"--path", "/home", "--path", "/etc"},
		},
		{
			name: "with tags",
			opts: []FilterOption{FilterByTag("daily", "prod")},
			want: []string{"--tag", "daily", "--tag", "prod"},
		},
		{
			name: "with latest",
			opts: []FilterOption{FilterLatest(5)},
			want: []string{"--latest", "5"},
		},
		{
			name: "combined options",
			opts: []FilterOption{
				FilterByHost("server1"),
				FilterByPath("/data"),
				FilterByTag("important"),
				FilterLatest(3),
			},
			want: []string{
				"--host", "server1",
				"--path", "/data",
				"--tag", "important",
				"--latest", "3",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var opts filterOptions
			for _, opt := range tt.opts {
				opt(&opts)
			}
			got := opts.args()
			if !stringSliceEqual(got, tt.want) {
				t.Errorf("filterOptions.args() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestForgetOptions_args verifies forget option argument generation
func TestForgetOptions_args(t *testing.T) {
	tests := []struct {
		name string
		opts []ForgetOption
		want []string
	}{
		{
			name: "with snapshot ID only",
			opts: []ForgetOption{ForgetSnapshot("abc12345")},
			want: []string{"abc12345"},
		},
		{
			name: "with keep last",
			opts: []ForgetOption{ForgetKeepLast(7)},
			want: []string{"--keep-last", "7"},
		},
		{
			name: "with prune",
			opts: []ForgetOption{ForgetWithPrune()},
			want: []string{"--prune"},
		},
		{
			name: "with hosts",
			opts: []ForgetOption{ForgetByHost("host1", "host2")},
			want: []string{"--host", "host1", "--host", "host2"},
		},
		{
			name: "combined with snapshot ID first",
			opts: []ForgetOption{
				ForgetSnapshot("abc12345"),
				ForgetByHost("server1"),
				ForgetKeepLast(5),
				ForgetWithPrune(),
			},
			want: []string{
				"abc12345",
				"--host", "server1",
				"--keep-last", "5",
				"--prune",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var opts forgetOptions
			for _, opt := range tt.opts {
				opt(&opts)
			}
			got := opts.args()
			if !stringSliceEqual(got, tt.want) {
				t.Errorf("forgetOptions.args() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestRestoreOptions_args verifies restore option argument generation
func TestRestoreOptions_args(t *testing.T) {
	tests := []struct {
		name string
		opts []RestoreOption
		want []string
	}{
		{
			name: "no options",
			opts: nil,
			want: []string{},
		},
		{
			name: "with hosts",
			opts: []RestoreOption{RestoreByHost("host1")},
			want: []string{"--host", "host1"},
		},
		{
			name: "with excludes and includes",
			opts: []RestoreOption{
				RestoreExclude("*.tmp", "*.log"),
				RestoreInclude("*.conf"),
			},
			want: []string{
				"--exclude", "*.tmp", "--exclude", "*.log",
				"--include", "*.conf",
			},
		},
		{
			name: "combined options",
			opts: []RestoreOption{
				RestoreByHost("server1"),
				RestoreByPath("/data"),
				RestoreByTag("prod"),
				RestoreExclude("*.tmp"),
			},
			want: []string{
				"--host", "server1",
				"--path", "/data",
				"--tag", "prod",
				"--exclude", "*.tmp",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var opts restoreOptions
			for _, opt := range tt.opts {
				opt(&opts)
			}
			got := opts.args()
			if !stringSliceEqual(got, tt.want) {
				t.Errorf("restoreOptions.args() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestParseStdErr verifies error parsing from stderr
func TestParseStdErr(t *testing.T) {
	tests := []struct {
		name   string
		stderr string
		want   error
	}{
		{
			name:   "repo already exists",
			stderr: "Fatal: config file already exists",
			want:   ErrRepoExists,
		},
		{
			name:   "wrong password",
			stderr: "Fatal: wrong password or no key found",
			want:   ErrInvalidPassword,
		},
		{
			name:   "invalid password variant",
			stderr: "invalid password",
			want:   ErrInvalidPassword,
		},
		{
			name:   "repo not found",
			stderr: "Is there a repository at the following location?",
			want:   ErrRepoNotFound,
		},
		{
			name:   "repo locked",
			stderr: "unable to create lock in backend",
			want:   ErrRepoLocked,
		},
		{
			name:   "repo locked variant",
			stderr: "repository is already locked",
			want:   ErrRepoLocked,
		},
		{
			name:   "invalid ID",
			stderr: "returned error, retrying after 1s",
			want:   ErrInvalidID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseStdErr(tt.stderr)
			if got != tt.want {
				t.Errorf("parseStdErr() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestIsSnapshotID verifies snapshot ID validation
func TestIsSnapshotID(t *testing.T) {
	tests := []struct {
		name string
		id   string
		want bool
	}{
		{
			name: "latest",
			id:   "latest",
			want: true,
		},
		{
			name: "latest with path",
			id:   "latest:/home/user",
			want: true,
		},
		{
			name: "short ID",
			id:   "a1b2c3d4",
			want: true,
		},
		{
			name: "short ID with path",
			id:   "a1b2c3d4:/data",
			want: true,
		},
		{
			name: "full ID",
			id:   "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
			want: true,
		},
		{
			name: "full ID with path",
			id:   "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef:/etc",
			want: true,
		},
		{
			name: "invalid - too short",
			id:   "abc123",
			want: false,
		},
		{
			name: "invalid - wrong length",
			id:   "a1b2c3d4e",
			want: false,
		},
		{
			name: "invalid - non-hex characters",
			id:   "g1h2i3j4",
			want: false,
		},
		{
			name: "empty",
			id:   "",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isSnapshotID(tt.id)
			if got != tt.want {
				t.Errorf("isSnapshotID(%q) = %v, want %v", tt.id, got, tt.want)
			}
		})
	}
}

// TestIsPathExists verifies path existence checking
func TestIsPathExists(t *testing.T) {
	tests := []struct {
		name string
		path string
		want bool
	}{
		{
			name: "current directory",
			path: ".",
			want: true,
		},
		{
			name: "this test file",
			path: "command_test.go",
			want: true,
		},
		{
			name: "non-existent path",
			path: "/nonexistent/path/that/does/not/exist",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isPathExists(tt.path)
			if got != tt.want {
				t.Errorf("isPathExists(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

// TestGetSummary verifies JSON summary extraction
func TestGetSummary(t *testing.T) {
	tests := []struct {
		name    string
		output  string
		wantErr bool
		check   func([]byte) bool
	}{
		{
			name:    "with summary message_type",
			output:  `{"message_type":"status","percent_done":0.5}` + "\n" + `{"message_type":"summary","files_new":10}`,
			wantErr: false,
			check: func(b []byte) bool {
				return strings.Contains(string(b), `"summary"`) && strings.Contains(string(b), `"files_new":10`)
			},
		},
		{
			name:    "with tags field (forget output)",
			output:  `{"message_type":"status"}` + "\n" + `{"tags":["daily"],"host":"myhost"}`,
			wantErr: false,
			check: func(b []byte) bool {
				return strings.Contains(string(b), `"tags":["daily"]`)
			},
		},
		{
			name:    "no summary found",
			output:  `{"message_type":"status"}`,
			wantErr: true,
		},
		{
			name:    "empty output",
			output:  "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getSummary(tt.output)
			if (err != nil) != tt.wantErr {
				t.Errorf("getSummary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.check != nil {
				if !tt.check(got) {
					t.Errorf("getSummary() result failed check: got %s", string(got))
				}
			}
		})
	}
}

// stringSliceEqual checks if two string slices are equal
func stringSliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
