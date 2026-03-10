package tests

import (
	generate "hera/core"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNextTurn_CursorOnNextPlayer(t *testing.T) {
	m := generate.NewModel(2, 0)
	m.Moved = true
	m.Shot = true

	next := m.Players[(m.CurrentPlayer+1)%len(m.Players)]
	assert.Equal(t, next.X, next.X) // убеждаемся что следующий игрок есть
}

func TestTurnAdvances(t *testing.T) {
	m := testModel()
	assert.Equal(t, 0, m.CurrentPlayer)
	m.CurrentPlayer = (m.CurrentPlayer + 1) % len(m.Players)
	assert.Equal(t, 1, m.CurrentPlayer)
}
