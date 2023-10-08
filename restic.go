package restic

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
)

const (
	resticBin  string = "restic"
	minVersion string = "0.16.0"
)

func init() {
	_, err := exec.LookPath(resticBin)
	if err != nil {
		fmt.Printf("[FATAL] %v\n", err)
		fmt.Printf("[FATAL] restic must be installed and exported in $PATH\n")
		fmt.Printf("[FATAL] https://restic.readthedocs.io/en/latest/020_installation.html\n")
		os.Exit(1)
	}

	cmd := exec.Command(resticBin, "version")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("[FATAL] %v\n", err)
		os.Exit(1)
	}

	versionRegex := regexp.MustCompile(minVersion)
	if !versionRegex.MatchString(string(out)) {
		fmt.Printf("[FATAL] restic must be minimum version %s\n", minVersion)
		os.Exit(1)
	}

}
