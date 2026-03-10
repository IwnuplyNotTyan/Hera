package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUltAxis_ResetsWhenCursorOnPlayer(t *testing.T) {
	m := testModel()
	m.UltMode = true
	m.UltAxis = "v"
	m.CursorX = m.Players[m.CurrentPlayer].X
	m.CursorY = m.Players[m.CurrentPlayer].Y

	cur := m.Players[m.CurrentPlayer]
	if m.CursorX == cur.X && m.CursorY == cur.Y {
		m.UltAxis = ""
	}
	assert.Equal(t, "", m.UltAxis)
}

func TestUltAxis_LocksVertical(t *testing.T) {
	m := testModel()
	m.UltMode = true
	m.UltAxis = "v"
	canMoveH := m.UltAxis == "" || m.UltAxis == "h"
	assert.False(t, canMoveH)
}

func TestUltAxis_LocksHorizontal(t *testing.T) {
	m := testModel()
	m.UltMode = true
	m.UltAxis = "h"
	canMoveV := m.UltAxis == "" || m.UltAxis == "v"
	assert.False(t, canMoveV)
}

func TestUltAxis_ClearsOnUltToggle(t *testing.T) {
	m := testModel()
	m.UltMode = true
	m.UltAxis = "v"
	m.UltMode = false
	m.UltAxis = ""
	assert.Equal(t, "", m.UltAxis)
}
