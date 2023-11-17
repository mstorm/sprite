package image

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	_, svg, err := Parse("testdata/airport.svg")

	assert.NoError(t, err)
	assert.Equal(t, 15, svg.Width())
	assert.Equal(t, 15, svg.Height())
}
