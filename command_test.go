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
