package restic

import (
	"errors"
	"testing"
)

func TestSentinelErrors(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{name: "ErrRepoExists", err: ErrRepoExists},
		{name: "ErrRepoNotFound", err: ErrRepoNotFound},
		{name: "ErrInvalidPassword", err: ErrInvalidPassword},
		{name: "ErrInvalidID", err: ErrInvalidID},
		{name: "ErrRepoLocked", err: ErrRepoLocked},
		{name: "ErrResticNotFound", err: ErrResticNotFound},
		{name: "ErrResticVersion", err: ErrResticVersion},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err == nil {
				t.Error("sentinel error is nil")
			}
			if tt.err.Error() == "" {
				t.Error("error message is empty")
			}
			if !errors.Is(tt.err, tt.err) {
				t.Errorf("errors.Is() failed for %v", tt.name)
			}
		})
	}
}

func TestErrorUnwrapping(t *testing.T) {
	tests := []struct {
		name   string
		base   error
		target error
	}{
		{name: "ErrRepoExists", base: ErrRepoExists, target: ErrRepoExists},
		{name: "ErrInvalidPassword", base: ErrInvalidPassword, target: ErrInvalidPassword},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !errors.Is(tt.base, tt.target) {
				t.Errorf("errors.Is(%v, %v) = false, want true", tt.base, tt.target)
			}
		})
	}
}

func TestErrorMessages(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{name: "ErrRepoExists", err: ErrRepoExists},
		{name: "ErrRepoNotFound", err: ErrRepoNotFound},
		{name: "ErrInvalidPassword", err: ErrInvalidPassword},
		{name: "ErrInvalidID", err: ErrInvalidID},
		{name: "ErrRepoLocked", err: ErrRepoLocked},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := tt.err.Error()
			if msg == "" {
				t.Error("error message is empty")
			}
		})
	}
}
