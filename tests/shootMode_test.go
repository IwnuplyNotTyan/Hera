package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsInRange_ShootModeSmaller(t *testing.T) {
    m := testModel()
    m.ShootMode = true
    assert.False(t, m.IsInRange(4, 2)) // dy = 3
}
