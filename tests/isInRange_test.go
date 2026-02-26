package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsInRange_BlockedByWall(t *testing.T) {
	m := testModel()
	assert.False(t, m.IsInRange(2, 5))
}

func TestIsInRange_FreeCell(t *testing.T) {
	m := testModel()
	assert.True(t, m.IsInRange(4, 4))
}

func TestIsInRange_PlayerCell(t *testing.T) {
	m := testModel()
	assert.False(t, m.IsInRange(4, 5))
}
