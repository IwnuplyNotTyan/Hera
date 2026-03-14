package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsInRange_ShootModeSmaller(t *testing.T) {
	m := testModel()
	m.ShootMode = true
	assert.True(t, m.IsInRange(4, 4))
	assert.True(t, m.IsInRange(4, 3))
	assert.False(t, m.IsInRange(4, 2))
}
