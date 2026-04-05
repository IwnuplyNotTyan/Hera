package i18n

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTranslator(t *testing.T) {
	loc, err := NewTranslator("../locales", "en")
	require.NoError(t, err)
	assert.Equal(t, "en", loc.GetLanguage())
}

func TestNewTranslator_InvalidPath(t *testing.T) {
	loc, err := NewTranslator("./nonexistent", "en")
	assert.NoError(t, err)
	assert.Equal(t, "en", loc.GetLanguage())
}

func TestTranslator_SetLanguage(t *testing.T) {
	loc, err := NewTranslator("../locales", "en")
	require.NoError(t, err)

	err = loc.SetLanguage("ru")
	require.NoError(t, err)
	assert.Equal(t, "ru", loc.GetLanguage())
}

func TestTranslator_SetLanguage_Invalid(t *testing.T) {
	loc, err := NewTranslator("../locales", "en")
	require.NoError(t, err)

	err = loc.SetLanguage("invalid")
	assert.Error(t, err)
}

func TestTranslator_AvailableLanguages(t *testing.T) {
	loc, err := NewTranslator("../locales", "en")
	require.NoError(t, err)

	languages := loc.AvailableLanguages()
	assert.Contains(t, languages, "en")
	assert.Contains(t, languages, "ru")
}

func TestTranslator_T_English(t *testing.T) {
	loc, err := NewTranslator("../locales", "en")
	require.NoError(t, err)

	result := loc.T("game.gameOver")
	assert.Equal(t, " Game Over ", result)
}

func TestTranslator_T_Russian(t *testing.T) {
	loc, err := NewTranslator("../locales", "en")
	require.NoError(t, err)

	err = loc.SetLanguage("ru")
	require.NoError(t, err)

	result := loc.T("game.gameOver")
	assert.Equal(t, " Конец игры ", result)
}

func TestTranslator_T_WithNumericInterpolation(t *testing.T) {
	loc, err := NewTranslator("../locales", "en")
	require.NoError(t, err)

	result := loc.T("status.player", 1, "♥ ♥ ♥")
	assert.Equal(t, "Player 1 ♥ ♥ ♥", result)
}

func TestTranslator_T_WithMapInterpolation(t *testing.T) {
	loc, err := NewTranslator("../locales", "en")
	require.NoError(t, err)

	result := loc.T("cursor.coordinates", map[string]interface{}{"x": 5, "y": 10})
	assert.Equal(t, "(5, 10) ", result)
}

func TestTranslator_T_FallbackToDefault(t *testing.T) {
	loc, err := NewTranslator("../locales", "en")
	require.NoError(t, err)

	err = loc.SetLanguage("ru")
	require.NoError(t, err)

	result := loc.T("status.ultCharges", 1)
	assert.Equal(t, "⽕×1", result)
}

func TestTranslator_T_WithEffectDuration(t *testing.T) {
	loc, err := NewTranslator("../locales", "en")
	require.NoError(t, err)

	result := loc.T("effects.fire", 3)
	assert.Equal(t, "⽕ 3", result)
}

func TestTranslator_T_MissingKey(t *testing.T) {
	loc, err := NewTranslator("../locales", "en")
	require.NoError(t, err)

	result := loc.T("nonexistent.key")
	assert.Equal(t, "nonexistent.key", result)
}

func TestTranslator_T_UltCharges(t *testing.T) {
	loc, err := NewTranslator("../locales", "en")
	require.NoError(t, err)

	result := loc.T("status.ultCharges", 3)
	assert.Equal(t, "⽕×3", result)
}

func TestTranslator_T_UltChargesZero(t *testing.T) {
	loc, err := NewTranslator("../locales", "en")
	require.NoError(t, err)

	result := loc.T("status.ultChargesZero")
	assert.Equal(t, "⽕×0", result)
}

func TestTranslator_T_AllKeysPresent(t *testing.T) {
	loc, err := NewTranslator("../locales", "en")
	require.NoError(t, err)

	keys := []string{
		"game.gameOver",
		"status.player",
		"status.move",
		"status.shoot",
		"status.ult",
		"status.ultCharges",
		"status.ultChargesZero",
		"cursor.coordinates",
		"cursor.player.you",
		"cursor.player.other",
		"cursor.player.wallBlocked",
		"cursor.enemy.default",
		"cursor.enemy.wallBlocked",
		"cursor.tile.wall",
		"cursor.tile.wallInWay",
		"cursor.tile.water",
		"cursor.tile.smoke",
		"cursor.tile.fire",
		"cursor.range.ult",
		"cursor.range.outOfUltAxis",
		"cursor.range.inShootRange",
		"cursor.range.inMoveRange",
		"cursor.range.empty",
		"effects.fire",
		"effects.wet",
		"effects.smoke",
		"help.moveUp",
		"help.moveDown",
		"help.moveLeft",
		"help.moveRight",
		"help.movePlayer",
		"help.changeMode",
		"help.secondAttack",
		"help.toggleHelp",
		"help.quit",
		"keys.up",
		"keys.down",
		"keys.left",
		"keys.right",
		"keys.confirm",
		"keys.shoot",
		"keys.ult",
		"keys.help",
		"keys.quit",
	}

	for _, key := range keys {
		result := loc.T(key)
		assert.NotEqual(t, key, result, "Key %s should have a translation", key)
	}
}

func TestTranslator_LocalesDir(t *testing.T) {
	loc, err := NewTranslator("./nonexistent", "en")
	require.NoError(t, err)
	assert.Equal(t, "en", loc.GetLanguage())
	assert.Equal(t, " Game Over ", loc.T("game.gameOver"))
}
