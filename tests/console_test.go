package tests

import (
	"testing"

	generate "hera/core"
	"hera/i18n"

	"github.com/stretchr/testify/assert"
)

var testTheme = generate.InitDefaultTheme()
var testLoc = func() i18n.Localizer {
	loc, _ := i18n.NewTranslator("../locales", "en")
	return loc
}()

func TestParseConsoleCmd_Empty(t *testing.T) {
	cmd := generate.ParseConsoleCmd("")
	assert.Equal(t, "help", cmd.Command)
}

func TestParseConsoleCmd_Simple(t *testing.T) {
	cmd := generate.ParseConsoleCmd("help")
	assert.Equal(t, "help", cmd.Command)
	assert.Len(t, cmd.Args, 0)
}

func TestParseConsoleCmd_WithArgs(t *testing.T) {
	cmd := generate.ParseConsoleCmd("heal p1 2")
	assert.Equal(t, "heal", cmd.Command)
	assert.Equal(t, []string{"p1", "2"}, cmd.Args)
}

func TestParseConsoleCmd_CaseInsensitive(t *testing.T) {
	cmd := generate.ParseConsoleCmd("HEAL P1 2")
	assert.Equal(t, "heal", cmd.Command)
}

func TestParseConsoleCmd_ExtraSpaces(t *testing.T) {
	cmd := generate.ParseConsoleCmd("  add   player  ")
	assert.Equal(t, "add", cmd.Command)
	assert.Equal(t, []string{"player"}, cmd.Args)
}

func TestParseTarget_Player(t *testing.T) {
	isPlayer, idx, ok := generate.ParseTarget("p1")
	assert.True(t, ok)
	assert.True(t, isPlayer)
	assert.Equal(t, 0, idx)
}

func TestParseTarget_Enemy(t *testing.T) {
	isPlayer, idx, ok := generate.ParseTarget("e1")
	assert.True(t, ok)
	assert.False(t, isPlayer)
	assert.Equal(t, 0, idx)
}

func TestParseTarget_Invalid(t *testing.T) {
	_, _, ok := generate.ParseTarget("x1")
	assert.False(t, ok)
	_, _, ok = generate.ParseTarget("p")
	assert.False(t, ok)
	_, _, ok = generate.ParseTarget("")
	assert.False(t, ok)
}

func TestExecHelp(t *testing.T) {
	m := generate.NewModel(2, 2, nil, nil, testLoc, testTheme, "", false, "dark")
	result := m.ExecuteConsoleCmd(generate.ParseConsoleCmd("help"))
	assert.Contains(t, result, "regen")
	assert.Contains(t, result, "heal")
	assert.Contains(t, result, "status")
}

func TestExecClear(t *testing.T) {
	m := generate.NewModel(2, 2, nil, nil, testLoc, testTheme, "", false, "dark")
	result := m.ExecuteConsoleCmd(generate.ParseConsoleCmd("clear"))
	assert.Equal(t, "", result)
}

func TestExecStatus(t *testing.T) {
	m := generate.NewModel(2, 2, nil, nil, testLoc, testTheme, "", false, "dark")
	m.Screen = generate.ScreenGame
	result := m.ExecuteConsoleCmd(generate.ParseConsoleCmd("status"))
	assert.Contains(t, result, "Player")
	assert.Contains(t, result, "Enemy")
}

func TestExecUnknownCommand(t *testing.T) {
	m := generate.NewModel(2, 2, nil, nil, testLoc, testTheme, "", false, "dark")
	result := m.ExecuteConsoleCmd(generate.ParseConsoleCmd("foobar"))
	assert.Contains(t, result, "Unknown")
}

func TestExecHeal(t *testing.T) {
	m := generate.NewModel(2, 2, nil, nil, testLoc, testTheme, "", false, "dark")
	m.Screen = generate.ScreenGame
	m.Players[0].HP = 1
	result := m.ExecuteConsoleCmd(generate.ParseConsoleCmd("heal p1 1"))
	assert.Equal(t, 2, m.Players[0].HP)
	assert.NotContains(t, result, "error")
	assert.NotContains(t, result, "unknown")
}

func TestExecDamage(t *testing.T) {
	m := generate.NewModel(2, 2, nil, nil, testLoc, testTheme, "", false, "dark")
	m.Screen = generate.ScreenGame
	oldHP := m.Players[0].HP
	result := m.ExecuteConsoleCmd(generate.ParseConsoleCmd("damage p1 1"))
	assert.Equal(t, oldHP-1, m.Players[0].HP)
	assert.NotContains(t, result, "error")
}

func TestExecSetHP(t *testing.T) {
	m := generate.NewModel(2, 2, nil, nil, testLoc, testTheme, "", false, "dark")
	m.Screen = generate.ScreenGame
	result := m.ExecuteConsoleCmd(generate.ParseConsoleCmd("set hp p1 2"))
	assert.Equal(t, 2, m.Players[0].HP)
	assert.NotContains(t, result, "error")
}

func TestExecEffectAdd(t *testing.T) {
	m := generate.NewModel(2, 2, nil, nil, testLoc, testTheme, "", false, "dark")
	m.Screen = generate.ScreenGame
	result := m.ExecuteConsoleCmd(generate.ParseConsoleCmd("effect p1 fire 3"))
	assert.True(t, generate.HasEffect(m.Players[0].Effects, generate.EffectFire))
	assert.NotContains(t, result, "error")
}

func TestExecEffectRemove(t *testing.T) {
	m := generate.NewModel(2, 2, nil, nil, testLoc, testTheme, "", false, "dark")
	m.Screen = generate.ScreenGame
	m.ExecuteConsoleCmd(generate.ParseConsoleCmd("effect p1 fire 3"))
	m.ExecuteConsoleCmd(generate.ParseConsoleCmd("remove effect p1 fire"))
	assert.False(t, generate.HasEffect(m.Players[0].Effects, generate.EffectFire))
}

func TestExecAddPlayer(t *testing.T) {
	m := generate.NewModel(2, 2, nil, nil, testLoc, testTheme, "", false, "dark")
	m.Screen = generate.ScreenGame
	oldCount := len(m.Players)
	result := m.ExecuteConsoleCmd(generate.ParseConsoleCmd("add player"))
	assert.Equal(t, oldCount+1, len(m.Players))
	assert.NotContains(t, result, "error")
}

func TestExecRemovePlayer(t *testing.T) {
	m := generate.NewModel(3, 2, nil, nil, testLoc, testTheme, "", false, "dark")
	m.Screen = generate.ScreenGame
	result := m.ExecuteConsoleCmd(generate.ParseConsoleCmd("remove player 1"))
	assert.Equal(t, 2, len(m.Players))
	assert.NotContains(t, result, "error")
}

func TestExecAddEnemy(t *testing.T) {
	m := generate.NewModel(2, 2, nil, nil, testLoc, testTheme, "", false, "dark")
	m.Screen = generate.ScreenGame
	oldCount := len(m.Enemys)
	result := m.ExecuteConsoleCmd(generate.ParseConsoleCmd("add enemy"))
	assert.Equal(t, oldCount+1, len(m.Enemys))
	assert.NotContains(t, result, "error")
}

func TestExecRemoveEnemy(t *testing.T) {
	m := generate.NewModel(2, 3, nil, nil, testLoc, testTheme, "", false, "dark")
	m.Screen = generate.ScreenGame
	result := m.ExecuteConsoleCmd(generate.ParseConsoleCmd("remove enemy 1"))
	assert.Equal(t, 2, len(m.Enemys))
	assert.NotContains(t, result, "error")
}

func TestExecInvalidTarget(t *testing.T) {
	m := generate.NewModel(2, 2, nil, nil, testLoc, testTheme, "", false, "dark")
	m.Screen = generate.ScreenGame
	result := m.ExecuteConsoleCmd(generate.ParseConsoleCmd("heal x1 1"))
	assert.Contains(t, result, "Invalid")

	result = m.ExecuteConsoleCmd(generate.ParseConsoleCmd("heal p99 1"))
	assert.Contains(t, result, "not found")
}

func TestExecRegen(t *testing.T) {
	m := generate.NewModel(2, 2, nil, nil, testLoc, testTheme, "", false, "dark")
	m.Screen = generate.ScreenGame
	result := m.ExecuteConsoleCmd(generate.ParseConsoleCmd("regen"))
	assert.NotEqual(t, 0, len(m.Walls))
	assert.NotContains(t, result, "error")
}

func TestExecHeal_NotInGame(t *testing.T) {
	m := generate.NewModel(2, 2, nil, nil, testLoc, testTheme, "", false, "dark")
	result := m.ExecuteConsoleCmd(generate.ParseConsoleCmd("heal p1 1"))
	assert.Contains(t, result, "only available during a game")
}
