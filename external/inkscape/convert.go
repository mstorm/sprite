package inkscape

import (
	"fmt"
	"os/exec"
)

const (
	Inkscape   = "inkscape"
	DefaultDPI = 96
)

func Convert(src, output string, ratio float64) ([]byte, error) {
	args := []string{fmt.Sprintf("--export-dpi=%f", DefaultDPI*ratio), src, "-o", output}
	return exec.Command(Inkscape, args...).CombinedOutput()
}
