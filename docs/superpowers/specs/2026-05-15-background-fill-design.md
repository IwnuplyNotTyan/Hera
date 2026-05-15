# Background Fill Feature

## Summary

Add an optional full-terminal background fill using the tint's `Bg()` color. Controllable via CLI flag (`--background` / `-b`) and togglable in the Settings menu.

## Motivation

The game grid is small (14×10 cells) and leaves most of the terminal empty. Currently, the empty area shows the terminal's default background color. With this feature, players can fill the entire terminal with the current theme's background color for a more immersive look.

## Design

### Activation

- CLI flag: `--background` / `-b` (default: `false`)
- Settings menu: toggle "Background: on/off" item, placed after the "Center: on/off" toggle
- The existing `Model.EnableBackground` field is used (already declared but unused)

### Color Source

Always reads from `m.Theme.Bg()` (the `ThemeRegistry` interface), which works identically in both `bubbletint` and `notint` builds.

### Rendering

A new method `fillBackground(s string) string` on `*Model`:

1. Split content into lines via `strings.Split`
2. For each line, pad with spaces rendered via `lipgloss.NewStyle().Background(m.Theme.Bg())` to reach `TerminalWidth`
3. Add empty background-colored lines to reach `TerminalHeight`
4. Apply to ALL screens (menu, settings, theme select, game)

Applied AFTER centering (if both are enabled): first center the content, then fill the surrounding area with the background color.

### Files Changed

| File | Change |
|---|---|
| `cmd/hera/main.go` | Add `--background` / `-b` flag |
| `core/struct.go` | Unused `EnableBackground` — no change needed |
| `core/fight.go` | `NewModel` gains `enableBackground bool` param; settings toggle in `doMenuConfirm` |
| `core/main.go` | `View()` game screen calls `fillBackground` |
| `core/menu.go` | `viewMenu`, `viewSettings`, `viewThemeSelect` call `fillBackground`; settings item + navigation |
| `i18n/locales/en.json` | Add `settings.background` key |
| `i18n/locales/ru.json` | Add `settings.background` key |

### Non-Goals

- No separate "window" abstraction or screen manager
- No new dependencies
- No changes to theme/tint system
