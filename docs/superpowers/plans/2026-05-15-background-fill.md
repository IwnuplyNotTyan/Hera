# Background Fill Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add an optional full-terminal background fill that uses the theme's `Bg()` color.

**Architecture:** A new `applyBackground` method on `*Model` pads rendered content to `TerminalWidth × TerminalHeight` with `lipgloss.NewStyle().Background(m.Theme.Bg())`. Called at the end of every `View`-family method (game, menu, settings, theme select) AFTER centering. Toggled via CLI flag `--background` and Settings menu item.

**Tech Stack:** Go, bubbletea, lipgloss

---

### Task 1: Add translation keys for background setting

**Files:**
- Modify: `i18n/locales/en.json`
- Modify: `i18n/locales/ru.json`

- [ ] **Step 1: Add `en.json` key**

Insert `"background": "Background"` in the `settings` block (after `center`):

```json
    "center": "Center",
    "background": "Background",
    "back": "Back"
```

- [ ] **Step 2: Add `ru.json` key**

Insert `"background": "Фон"` in the `settings` block:

```json
    "center": "Центрирование",
    "background": "Фон",
    "back": "Назад"
```

- [ ] **Step 3: Commit**

```bash
git add i18n/locales/en.json i18n/locales/ru.json && git commit -m "i18n: add background setting key"
```

---

### Task 2: Add `--background` / `-b` CLI flag

**Files:**
- Modify: `cmd/hera/main.go`

- [ ] **Step 1: Add flag variable and register**

In `main()` at `cmd/hera/main.go`, add `var noBackground bool` near line 23:

```go
var lang string
var theme string
var noCenter bool
var noBackground bool
```

Register the flag near line 63:

```go
cmd.Flags().BoolVarP(&noBackground, "no-background", "b", false, "Disable background fill")
```

- [ ] **Step 2: Set `EnableBackground` on the model**

After `centerWindow := !noCenter` (near line 42), add:

```go
centerWindow := !noCenter
enableBackground := !noBackground
```

And after `model.SetAvailableThemes()` (line 48), set:

```go
model.EnableBackground = enableBackground
```

- [ ] **Step 3: Verify compilation**

```bash
go build ./...
```

Expected: clean exit with no errors.

- [ ] **Step 4: Commit**

```bash
git add cmd/hera/main.go && git commit -m "cmd: add --no-background CLI flag"
```

---

### Task 3: Add `applyBackground` method and integrate in game View()

**Files:**
- Modify: `core/main.go`

- [ ] **Step 1: Add `applyBackground` method**

Add this method to `core/main.go` (before `View()` or after it):

```go
func (m *Model) applyBackground(s string) string {
	if !m.EnableBackground || m.TerminalWidth <= 0 || m.TerminalHeight <= 0 {
		return s
	}

	bg := m.Theme.Bg()
	bgStyle := lipgloss.NewStyle().Background(bg)

	lines := strings.Split(s, "\n")
	for i, line := range lines {
		w := lipgloss.Width(line)
		if w < m.TerminalWidth {
			lines[i] = line + bgStyle.Render(strings.Repeat(" ", m.TerminalWidth-w))
		}
	}
	for len(lines) < m.TerminalHeight {
		lines = append(lines, bgStyle.Render(strings.Repeat(" ", m.TerminalWidth)))
	}

	return strings.Join(lines[:m.TerminalHeight], "\n")
}
```

- [ ] **Step 2: Integrate in game `View()` — game-over case**

Replace the early return at line 236:

```go
// Before (line 236):
	return gameOver

// After:
	return m.applyBackground(gameOver)
```

- [ ] **Step 3: Integrate at the end of game `View()`**

Before `return content` at line 466, add:

```go
	content = m.applyBackground(content)
	return content
```

So the block becomes:

```go
	content = m.applyBackground(content)
	return content
```

Remove the old `return content` at line 466.

- [ ] **Step 4: Verify compilation**

```bash
go build ./...
```

Expected: clean exit.

- [ ] **Step 5: Commit**

```bash
git add core/main.go && git commit -m "core: add applyBackground method for full-terminal fill"
```

---

### Task 4: Add `enableBackground` param to `NewModel` and settings toggle

**Files:**
- Modify: `core/fight.go`

- [ ] **Step 1: Add `enableBackground` param to `NewModel`**

Change the signature at line 14:

```go
func NewModel(playerCount, enemysCount int, loc i18n.Localizer, theme ThemeRegistry, centerWindow bool, enableBackground bool, themeName string) Model {
```

In the returned `Model` literal (line 75-95), add `EnableBackground: enableBackground,`:

```go
return Model{
    Theme:            theme,
    ThemeName:        themeName,
    Styles:           styles,
    EnableBackground: enableBackground,
    CenterWindow:     centerWindow,
    // ... rest unchanged
}
```

- [ ] **Step 2: Add background toggle in `doMenuConfirm`**

After the center toggle block at line 660-663, add the background toggle:

```go
	if m.MenuSelected == n {
		m.CenterWindow = !m.CenterWindow
		return m, nil
	}
	n++
	if m.MenuSelected == n {
		m.EnableBackground = !m.EnableBackground
		return m, nil
	}
	m.Screen = ScreenMenu
```

- [ ] **Step 3: Update caller in `cmd/hera/main.go`**

In `cmd/hera/main.go` at line 47, pass the new param:

```go
model := generate.NewModel(rand.Intn(3)+2, rand.Intn(3)+2, loc, registry, centerWindow, enableBackground, themeName)
```

- [ ] **Step 4: Verify compilation**

```bash
go build ./...
```

Expected: clean exit.

- [ ] **Step 5: Commit**

```bash
git add core/fight.go cmd/hera/main.go && git commit -m "core: add enableBackground param and settings toggle"
```

---

### Task 5: Integrate background in menu, settings, theme screens

**Files:**
- Modify: `core/menu.go`

- [ ] **Step 1: Add background to `viewSettings`**

In `viewSettings()`, after the `centerStr` line (line 183), add `bgStr`:

```go
	centerStr := "on"
	if !m.CenterWindow {
		centerStr = "off"
	}
	bgStr := "on"
	if !m.EnableBackground {
		bgStr = "off"
	}
```

At line 200-203 (the `menuItems = append` block), add the background item:

```go
	menuItems = append(menuItems,
		m.Localizer.T("settings.center")+": "+centerStr,
		m.Localizer.T("settings.background")+": "+bgStr,
		m.Localizer.T("settings.back"),
	)
	figures = append(figures, " ◆ ", " ◆ ", " ● ")
```

Before `return content` at the end of `viewSettings()`, add:

```go
	content = m.applyBackground(content)
	return content
```

(Remove the old `return content`.)

- [ ] **Step 2: Add background to `viewMenu`**

Before `return content` at the end of `viewMenu()` (line 176), add:

```go
	content = m.applyBackground(content)
	return content
```

(Remove the old `return content` at line 176.)

- [ ] **Step 3: Add background to `viewThemeSelect`**

Before `return content` at the end of `viewThemeSelect()` (line 383), add:

```go
	content = m.applyBackground(content)
	return content
```

(Remove the old `return content` at line 383.)

- [ ] **Step 4: Fix keyboard navigation in `updateMenu`**

In the "down" case (line 450-461), change `maxItem := 1` to `maxItem := 2`:

```go
			} else if m.Screen == ScreenSettings {
				maxItem := 2
				if len(m.Localizer.AvailableLanguages()) > 1 {
					maxItem++
				}
```

In the "up" case (line 432-441), change `maxItem := 1` to `maxItem := 2`:

```go
			} else if m.Screen == ScreenSettings && m.MenuSelected < 0 {
				maxItem := 2
				if len(m.Localizer.AvailableLanguages()) > 1 {
					maxItem++
				}
```

- [ ] **Step 5: Verify compilation**

```bash
go build ./...
```

Expected: clean exit.

- [ ] **Step 6: Verify with vet**

```bash
go vet ./...
```

Expected: clean exit.

- [ ] **Step 7: Run tests**

```bash
go test ./...
```

Expected: all tests pass.

- [ ] **Step 8: Commit**

```bash
git add core/menu.go && git commit -m "core: add background fill to menu/settings/theme screens"
```
