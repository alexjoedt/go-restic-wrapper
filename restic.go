package restic

import (
	"fmt"
	"os/exec"
	"strings"
	"sync"

	"github.com/hashicorp/go-version"
)

const (
	resticBin  string = "restic"
	minVersion string = "0.16.0"
)

var (
	// resticCheckOnce ensures version check happens only once
	resticCheckOnce sync.Once
	// resticCheckErr stores the result of version check
	resticCheckErr error
)

// checkResticVersion verifies that restic binary exists and meets minimum version.
// This check is performed lazily on first use and cached for subsequent calls.
func checkResticVersion() error {
	resticCheckOnce.Do(func() {
		// Check if restic binary exists in PATH
		_, err := exec.LookPath(resticBin)
		if err != nil {
			resticCheckErr = fmt.Errorf("restic not found in PATH: %w (install from https://restic.readthedocs.io/en/latest/020_installation.html)", err)
			return
		}

		// Get restic version
		cmd := exec.Command(resticBin, "version")
		out, err := cmd.CombinedOutput()
		if err != nil {
			resticCheckErr = fmt.Errorf("failed to get restic version: %w", err)
			return
		}

		// Parse version string (format: "restic 0.16.0 compiled with go1.21.0...")
		fields := strings.Fields(string(out))
		if len(fields) < 2 {
			resticCheckErr = fmt.Errorf("unexpected restic version output: %s", string(out))
			return
		}

		v, err := version.NewVersion(fields[1])
		if err != nil {
			resticCheckErr = fmt.Errorf("failed to parse restic version %q: %w", fields[1], err)
			return
		}

		minV, err := version.NewVersion(minVersion)
		if err != nil {
			resticCheckErr = fmt.Errorf("invalid minimum version %q: %w", minVersion, err)
			return
		}

		if v.LessThan(minV) {
			resticCheckErr = fmt.Errorf("restic version %s is too old, minimum required version is %s", v.String(), minVersion)
			return
		}
	})

	return resticCheckErr
}
