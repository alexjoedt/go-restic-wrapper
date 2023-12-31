package restic

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/hashicorp/go-version"
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

	fields := strings.Fields(string(out))
	v, err := version.NewVersion(fields[1])
	must(err)

	minV, err := version.NewVersion(minVersion)
	must(err)

	if v.LessThan(minV) {
		fmt.Printf("[FATAL] restic must be minimum version %s\n", minVersion)
		os.Exit(1)
	}

}

func must(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
