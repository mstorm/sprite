package inkscape

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvert(t *testing.T) {
	inputFile := "testdata/dot.svg"
	outputFile := "testdata/dot.png"

	_, err := convert(inputFile, outputFile, 10)

	assert.NoError(t, err)
	assert.FileExists(t, outputFile)

	assert.NoError(t, os.Remove(outputFile))
}
