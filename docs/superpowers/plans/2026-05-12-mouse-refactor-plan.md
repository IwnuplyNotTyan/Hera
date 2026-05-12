# Mouse Input Refactoring Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Refactor mouse handling: centralized in one place, no global state, pointer receivers on Model, theme click works, `--no-center` consistent.

**Architecture:** Move `layoutElements`/`gridOffsetX`/`gridOffsetY` from package-level globals into Model fields. Change all Model methods from value to pointer receivers so `View()` can write transient render state. Remove duplicate mouse handling from `updateMenu()`. Fix `ElementThemeItem` click handler.

**Tech Stack:** Go, bubbletea, lipgloss

---

### Task 1: Add transient render state to Model + rewrite mouse.go

**Files:**
- Modify: `core/struct.go` — add `layoutElements`, `gridOffsetX`, `gridOffsetY` to Model
- Rewrite: `core/mouse.go` — global functions → Model methods, remove globals

- [ ] **Step 1: Add fields to Model**

In `core/struct.go`, add to the `Model` struct after `Localizer`:
```go
    layoutElements []Element
    gridOffsetX, gridOffsetY int
```

- [ ] **Step 2: Rewrite core/mouse.go**

Replace the entire file with:

```go
package generate

func (m *Model) resetLayout() {
    m.layoutElements = nil
    m.gridOffsetX = 0
    m.gridOffsetY = 0
}

func (m *Model) trackElement(elem Element) {
    m.layoutElements = append(m.layoutElements, elem)
}

func (m *Model) hitTest(screenX, screenY int) *Element {
    absX := screenX - m.gridOffsetX
    absY := screenY - m.gridOffsetY
    for i := range m.layoutElements {
        elem := &m.layoutElements[i]
        if absX >= elem.X && absX < elem.X+elem.Width &&
            absY >= elem.Y && absY < elem.Y+elem.Height {
            return elem
        }
    }
    return nil
}

func cellWidth() int { return 3 }
func cellHeight() int { return 1 }
```

- [ ] **Step 3: Build to check (expected failure)**

Run: `go build ./...`
Expected: compile errors — `main.go`, `menu.go`, etc. call `trackElement()`, `resetLayout()`, `hitTest()` as package functions, but they're now methods. This is fine; Task 2 fixes all callers.

- [ ] **Step 4: Commit (struct + mouse.go only, even though build is broken)**

```bash
git add core/struct.go core/mouse.go
git commit -m "refactor: move mouse render state into Model fields"
```

---

### Task 2: Change all Model methods to pointer receivers (fixes Task 1 build)

**Files:**
- Modify: `core/main.go`
- Modify: `core/menu.go`
- Modify: `core/fight.go`
- Modify: `core/info.go`

**Rules:**
- Read-only methods: change receiver to `(m *Model)`, keep return type.
- Mutating methods: change receiver to `(m *Model)`, return `*Model` (was `Model`), or `(tea.Model, tea.Cmd)` (was `(Model, tea.Cmd)`).
- `View()`: change receiver to `(m *Model)`, return type stays `string`.
- Inside methods: replace `trackElement(...)` with `m.trackElement(...)`, `resetLayout()` with `m.resetLayout()`, `hitTest(...)` with `m.hitTest(...)`.
- Replace `gridOffsetX` with `m.gridOffsetX`, `gridOffsetY` with `m.gridOffsetY`.

- [ ] **Step 1: Change core/main.go — all methods to pointer receivers**

```go
func (m *Model) Init() tea.Cmd {
    return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // Replace hitTest(msg.X, msg.Y) → m.hitTest(msg.X, msg.Y)
    // Replace resetLayout() → m.resetLayout()  (removed from here, goes to View)
    // m = m.doConfirm() stays same (both *Model)
    // return m, nil — *Model implements tea.Model
}

func (m *Model) View() string {
    m.resetLayout()
    // Replace trackElement(...) → m.trackElement(...)
    // Replace gridOffsetX → m.gridOffsetX, gridOffsetY → m.gridOffsetY
    return content
}
```

- [ ] **Step 2: Change core/menu.go — all methods to pointer receivers**

```go
func (m *Model) startGame() { /* already *Model */ }
// navigateTheme already *Model, unchanged
func (m *Model) viewMenu() string {
    // trackElement(...) → m.trackElement(...)
    // gridOffsetX/Y → m.gridOffsetX/Y
}
func (m *Model) viewSettings() string { /* same pattern */ }
func (m *Model) viewThemeSelect() string { /* same pattern */ }
func (m *Model) updateMenu(msg tea.Msg) (tea.Model, tea.Cmd) {
    // receiver change, body unchanged (for now)
    // return m, nil
}
```

- [ ] **Step 3: Change core/fight.go — methods to pointer receivers**

```go
// NewModel stays as func NewModel(...) Model (factory, not a method)
func NewModel(...) Model { ... }

func (m *Model) closestPlayer(ex, ey int) (int, int) { ... }
func (m *Model) enemyOccupied(x, y, skipIdx int) bool { ... }
func (m *Model) doEnemyTurn(idx int) *Model { ... return m }
func (m *Model) Move(newX, newY int) *Model { return m }
func (m *Model) currentRange() int { ... }
func (m *Model) IsInRange(col, row int) bool { ... }
func (m *Model) Reachable(sx, sy, r int) map[Point]bool { ... }
func (m *Model) HasWallBetweenPoints(x0, y0, x1, y1 int) bool { ... }
func (m *Model) ultCross(cx, cy int) []Point { ... }
func (m *Model) ultInAxisRange(cx, cy int) bool { ... }
func (m *Model) doConfirm() *Model { ... return m }
func (m *Model) doUlt() *Model { ... return m }
func (m *Model) tickFireTiles() *Model { ... return m }
func (m *Model) nextTurn() *Model { ... return m }
// enemyTurnCmd stays standalone
func (m *Model) OccupiedByOther(x, y int) bool { ... }
func (m *Model) turnOrder() string { ... }
func (m *Model) doMenuConfirm() (tea.Model, tea.Cmd) { ... return m, nil }
```

- [ ] **Step 4: Change core/info.go — methods to pointer receivers**

```go
func (m *Model) cursorInfo() string {
    // receiver change only
}
```

- [ ] **Step 5: Build to check errors**

Run: `go build ./...`
Expected: errors in cmd/hera/main.go (value vs pointer) and tests. Note them.

- [ ] **Step 6: Commit**

```bash
git add core/main.go core/menu.go core/fight.go core/info.go
git commit -m "refactor: change all Model methods to pointer receivers"
```

---

### Task 3: Update cmd/hera/main.go for pointer Model

**Files:**
- Modify: `cmd/hera/main.go`

- [ ] **Step 1: Pass &model to tea.NewProgram**

```go
model := generate.NewModel(rand.Intn(3)+2, rand.Intn(3)+2, loc, registry, centerWindow, themeName)
model.SetAvailableThemes()
p := tea.NewProgram(
    &model,
    tea.WithAltScreen(),
    tea.WithMouseCellMotion(),
)
```

- [ ] **Step 2: Build**

Run: `go build ./...`
Expected: core builds, tests may still fail

- [ ] **Step 3: Commit**

```bash
git add cmd/hera/main.go
git commit -m "fix: pass &model to tea.NewProgram for pointer receiver compatibility"
```

---

### Task 4: Fix tests for pointer receivers

**Files:**
- Modify: `tests/fight_test.go`
- Modify: `tests/hp_test.go`
- Modify: `tests/isInRange_test.go`
- Modify: `tests/nextTurn_test.go`
- Modify: `tests/shootMode_test.go`
- Modify: `tests/ult_test.go`
- Modify: `tests/wallBetween_test.go`

- [ ] **Step 1: Fix testModel() in fight_test.go**

Change return to `*generate.Model`:
```go
func testModel() *generate.Model {
    // ... same body ...
    return &generate.Model{
        Players:       players,
        CurrentPlayer: 0,
        CursorX:       4,
        CursorY:       5,
        Walls:         walls,
        Water:         water,
        FireTiles:     map[generate.Point]int{},
        SmokeTiles:    map[generate.Point]int{},
        Enemys:        []generate.Enemy{},
        Localizer:     loc,
    }
}
```

- [ ] **Step 2: Fix createTestModel() in hp_test.go**

```go
func createTestModel() *generate.Model {
    return &generate.Model{
        Theme:         theme,
        Players:       players,
        // ... rest same, prefixed with &
    }
}
```

- [ ] **Step 3: Fix inline Model constructions**

`tests/isInRange_test.go` `TestReachable_BlockedByWall`:
```go
m := &generate.Model{
    Players: players, ...
}
```

`tests/wallBetween_test.go` `TestHasWallBetween_StartNotCounted`:
```go
m := &generate.Model{
    Players: players, ...
}
```

- [ ] **Step 4: Build and test**

Run: `go build ./... && go test ./...`
Expected: all pass

- [ ] **Step 5: Commit**

```bash
git add tests/
git commit -m "test: adjust for *Model pointer receivers"
```

---

### Task 5: Centralize mouse + fix theme click + remove duplicate handler

**Files:**
- Modify: `core/menu.go`
- Modify: `core/main.go`

- [ ] **Step 1: Remove mouse handling from updateMenu**

In `core/menu.go`, delete:
```go
case tea.MouseMsg:
    if msg.Button != tea.MouseButtonLeft || msg.Action != tea.MouseActionPress {
        return m, nil
    }
    elem := hitTest(msg.X, msg.Y)
    if elem != nil {
        switch elem.Type {
        case ElementMenuItem, ElementSettingsItem:
            m.MenuSelected = elem.Index
        case ElementThemeItem:
        }
    }
    return m, nil
```

Keep `case tea.WindowSizeMsg:`.

- [ ] **Step 2: Fix theme element tracking to store theme name**

In `core/menu.go` `viewThemeSelect()`, change:
```go
m.trackElement(Element{
    Type:   ElementThemeItem,
    X:      4,
    Y:      row,
    Width:  lipgloss.Width("   ●  ") + lipgloss.Width(theme),
    Height: 1,
    ID:     theme,  // was fmt.Sprintf("theme-%d", i)
    Index:  i,
})
```

- [ ] **Step 3: Centralize ALL mouse handling in main.go**

In `core/main.go` `Update()`, replace the mouse handler:

```go
case tea.MouseMsg:
    elem := m.hitTest(msg.X, msg.Y)
    if elem != nil {
        switch elem.Type {
        case ElementGridCell:
            if m.Screen == ScreenGame && !m.EnemyTurn {
                var col, row int
                if _, err := fmt.Sscanf(elem.ID, "cell-%d-%d", &col, &row); err == nil {
                    m.CursorX = col
                    m.CursorY = row
                }
            }
        case ElementMenuItem:
            if m.Screen == ScreenMenu {
                m.MenuSelected = elem.Index
            }
        case ElementSettingsItem:
            if m.Screen == ScreenSettings {
                m.MenuSelected = elem.Index
            }
        case ElementThemeItem:
            if m.Screen == ScreenThemeSelect {
                m.ThemeName = elem.ID
                if m.Theme != nil {
                    m.Theme.SetTintID(elem.ID)
                }
                m.Styles = NewStyles(m.Theme)
            }
        }
    }
```

- [ ] **Step 4: Set gridOffsetX/Y in View for all screens**

In `core/main.go` `View()` — at the centering block, add else branch:
```go
if m.CenterWindow && m.TerminalWidth > 0 && m.TerminalHeight > 0 {
    // ... existing centering ...
    m.gridOffsetX = marginX
    m.gridOffsetY = marginY
} else {
    m.gridOffsetX = 0
    m.gridOffsetY = 0
}
```

Same pattern in `viewMenu()`, `viewSettings()`, `viewThemeSelect()`.

- [ ] **Step 5: Build and test**

Run: `go build ./... && go test ./...`
Expected: all pass

- [ ] **Step 6: Run vet**

Run: `go vet ./...`
Expected: no warnings

- [ ] **Step 7: Commit**

```bash
git add core/main.go core/menu.go
git commit -m "feat: centralize mouse handling, fix theme click, remove duplicate handler"
```

---

### Task 6: Final verification

- [ ] **Step 1: Full build**

```bash
go build ./...
```

- [ ] **Step 2: Full test**

```bash
go test ./...
```

- [ ] **Step 3: Vet**

```bash
go vet ./...
```

- [ ] **Step 4: Review diff**

```bash
git diff --stat
```

Expected: clean changes across ~8 files
