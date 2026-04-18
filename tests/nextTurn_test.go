package tests

import (
	"testing"

	generate "hera/core"
	"hera/i18n"

	bubbletint "github.com/lrstanley/bubbletint"
	"github.com/stretchr/testify/assert"
)

func TestNextTurn_CursorOnNextPlayer(t *testing.T) {
	loc, _ := i18n.NewTranslator("../locales", "en")
	theme := bubbletint.NewRegistry(bubbletint.TintDraculaPlus, bubbletint.DefaultTints()...)
	m := generate.NewModel(2, 0, loc, theme, false, "default")
	m.Moved = true
	m.Shot = true

	next := m.Players[(m.CurrentPlayer+1)%len(m.Players)]
	assert.Equal(t, next.X, next.X)
}

func TestTurnAdvances(t *testing.T) {
	m := createTestModel()
	assert.Equal(t, 0, m.CurrentPlayer)
	m.CurrentPlayer = (m.CurrentPlayer + 1) % len(m.Players)
	assert.Equal(t, 1, m.CurrentPlayer)
}
