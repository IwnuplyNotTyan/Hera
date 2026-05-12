# Mouse Input Refactoring Design

## Problem

Current mouse logic has structural issues:

1. **Duplicate mouse handling**: Both `main.go:Update()` and `menu.go:updateMenu()` process `tea.MouseMsg` for the same menu/settings/theme screens, with the second handler overriding the first.
2. **Package-level globals**: `layoutElements`, `gridOffsetX`, `gridOffsetY` in `mouse.go` are mutable global state.
3. **Theme click does nothing**: `case ElementThemeItem:` has empty body.
4. **Model uses value receivers everywhere** — wasteful for a large struct, and prevents View() from writing transient render state.

## Solution

### 1. All Model methods → pointer receivers

Every method on `Model` changes from `(m Model)` to `(m *Model)` (including `Init`, `Update`, `View`).

In `main.go` (cmd), pass `&model` to `tea.NewProgram`.

`Update` returns `m` (which is `*Model`, satisfying `tea.Model`).

### 2. Transient render state in Model

Add fields to `struct.go` `Model`:
```go
layoutElements []Element
gridOffsetX, gridOffsetY int
```

These are scratch buffers, set each frame in `View()`, consumed by `hitTest()` in `Update()`.

### 3. mouse.go methods → Model methods

- `resetLayout()` → `(m *Model) resetLayout()`
- `trackElement()` → `(m *Model) trackElement(Element)`
- `hitTest()` → `(m *Model) hitTest(screenX, screenY) *Element`
- `cellWidth()`, `cellHeight()` → static package (no state)

Remove package-level globals.

### 4. Remove mouse handling from updateMenu

Delete `case tea.MouseMsg:` block in `menu.go:updateMenu()` (lines 370–382). All mouse events handled centrally in `main.go:Update()`.

### 5. Fix ElementThemeItem click

Clicking a theme item navigates to that theme (preview), matching keyboard up/down behavior. Right-click (or Enter) confirms selection and returns to settings.

In the main `Update()` mouse handler:
```go
case ElementThemeItem:
    m.MenuSelected = elem.Index
    m.navigateToThemeByIndex(elem.Index)
```

Where `navigateToThemeByIndex(index)` resolves the filtered display index to an absolute theme name and sets `m.ThemeName` + refreshes `m.Styles` (same as `navigateTheme` but by index).

Since `trackElement` for themes stores the filtered-list index (not the absolute index), the handler must map through the filtered list to find the actual theme name.

### 6. gridOffsetX/Y always set

- When `CenterWindow` is true: set to margin values
- When `--no-center` (CenterWindow is false): set to 0,0

## Files Changed

| File | Change |
|------|--------|
| `core/struct.go` | Add `layoutElements`, `gridOffsetX`, `gridOffsetY` fields |
| `core/mouse.go` | Rewrite as Model methods, remove globals |
| `core/main.go` | Pointer receivers, centralized mouse handling |
| `core/menu.go` | Pointer receivers, remove mouse from updateMenu, fix theme click |
| `core/fight.go` | Pointer receivers |
| `core/info.go` | Pointer receivers |
| `core/pallet.go` | No changes needed |
| `core/effect.go` | No changes needed |
| `core/tiles.go` | No changes needed |
| `core/help.go` | No changes needed |
| `cmd/hera/main.go` | Pass `&model` to `tea.NewProgram` |

## `--no-center` Behavior

- `gridOffsetX = 0`, `gridOffsetY = 0`
- Content starts at terminal position (0,0)
- Elements tracked at content-relative coordinates
- `hitTest` returns screen coords unchanged → maps 1:1 to content coords
- Works correctly without centering branch

## Migration Order

1. Add fields to Model (struct.go)
2. Rewrite mouse.go methods
3. Change all Model methods to pointer receivers
4. Remove mouse from updateMenu
5. Fix theme click
6. Update main.go to pass `&model`
7. Build, test, vet
