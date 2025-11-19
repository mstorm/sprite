package resvg

import (
	"fmt"
	"os/exec"
)

const (
	Resvg = "resvg"
)

func Convert(src, output string, ratio float64) ([]byte, error) {
	args := []string{"-z", fmt.Sprintf("%f", ratio), src, output}
	return exec.Command(Resvg, args...).CombinedOutput()
}
