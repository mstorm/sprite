package inkscape

import (
	"fmt"
	"os/exec"
)

const (
	Inkscape   = "inkscape"
	DefaultDPI = 96
)

func convert(targetFile, outputFile string, ratio float32) ([]byte, error) {
	args := []string{fmt.Sprintf("--export-dpi=%f", DefaultDPI*ratio), targetFile, "-o", outputFile}
	return exec.Command(Inkscape, args...).CombinedOutput()
}
