package tests

import (
	"testing"

	"hera/utils"

	"github.com/stretchr/testify/assert"
)

func TestAbs(t *testing.T) {
	assert.Equal(t, 5, utils.Abs(-5))
	assert.Equal(t, 5, utils.Abs(5))
	assert.Equal(t, 0, utils.Abs(0))
}
