package resvg

import (
	"fmt"
	"os/exec"
)

const (
	Resvg = "resvg"
)

func Convert(src, output string, ratio float64) ([]byte, error) {
	args := []string{src, output, fmt.Sprintf("--zoom=%f", ratio)}
	return exec.Command(Resvg, args...).CombinedOutput()
}
