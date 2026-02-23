package tests

import (
	"testing"

	"hera/utils"

	"github.com/stretchr/testify/assert"
)

func TestClamp(t *testing.T) {
	assert.Equal(t, 0, utils.Clamp(-5, 0, 9))
	assert.Equal(t, 9, utils.Clamp(15, 0, 9))
	assert.Equal(t, 5, utils.Clamp(5, 0, 9))
	assert.Equal(t, 0, utils.Clamp(0, 0, 9))
	assert.Equal(t, 9, utils.Clamp(9, 0, 9))
}
