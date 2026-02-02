package restic

import "errors"

// Sentinel errors returned by the library.
var (
	// ErrRepoExists is returned when trying to initialize a repository that already exists.
	ErrRepoExists = errors.New("repository already exists")

	// ErrRepoNotFound is returned when the specified repository doesn't exist.
	ErrRepoNotFound = errors.New("repository not found")

	// ErrInvalidPassword is returned when the repository password is incorrect.
	ErrInvalidPassword = errors.New("invalid repository password")

	// ErrInvalidID is returned when a snapshot ID has an invalid format.
	ErrInvalidID = errors.New("invalid snapshot ID")

	// ErrRepoLocked is returned when the repository is locked by another process.
	ErrRepoLocked = errors.New("repository locked by another process")

	// ErrResticNotFound is returned when the restic binary is not found in PATH.
	ErrResticNotFound = errors.New("restic binary not found in PATH")

	// ErrResticVersion is returned when the restic version doesn't meet minimum requirements.
	ErrResticVersion = errors.New("restic version does not meet minimum requirements")
)
