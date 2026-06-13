package generate

import (
	"time"

	"hera/i18n"
	"hera/utils"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func NewModel(playerCount, enemysCount int, playerEffects, enemyEffects []Effect, loc i18n.Localizer, theme ThemeRegistry, centerWindow string, enableBackground bool, themeName string) Model {
	styles := NewStyles(theme)
	if playerCount < 2 {
		playerCount = 2
	}
	if playerCount > 4 {
		playerCount = 4
	}
	if enemysCount < 0 {
		enemysCount = 0
	}

	blocked := make(map[Point]bool)

	walls := GenerateTiles(GridW/2, GridH/2, wallCount, nil)
	water := GenerateTiles(GridW/2, GridH/2, waterCount, walls)

	starts := []Point{
		{X: 1, Y: 1},
		{X: GridW - 2, Y: GridH - 2},
		{X: GridW - 2, Y: 1},
		{X: 1, Y: GridH - 2},
	}

	players := make([]Player, playerCount)
	for i := range players {
		effs := make([]Effect, len(playerEffects))
		copy(effs, playerEffects)
		players[i] = Player{
			X:          starts[i].X,
			Y:          starts[i].Y,
			HP:         MaxHP,
			UltCharges: maxUltCharges,
			Effects:    effs,
			Style:      styles.PlayerStyles[i%len(styles.EnemysStyles)],
		}
	}

	for _, p := range players {
		blocked[Point{p.X, p.Y}] = true
	}
	for p := range walls {
		blocked[p] = true
	}
	for p := range water {
		blocked[p] = true
	}

	enemyStarts := GenerateTiles(GridW/2, GridH/2, enemysCount, blocked)
	enemyPositions := make([]Point, 0, enemysCount)
	for p := range enemyStarts {
		enemyPositions = append(enemyPositions, p)
	}

	enemys := make([]Enemy, enemysCount)
	for i := range enemys {
		effs := make([]Effect, len(enemyEffects))
		copy(effs, enemyEffects)
		enemys[i] = Enemy{
			X:       enemyPositions[i].X,
			Y:       enemyPositions[i].Y,
			HP:      MaxHP,
			Effects: effs,
			Style:   styles.EnemysStyles[i],
		}
	}

	return Model{
		Theme:              theme,
		Config:             &Config{ThemeName: themeName, CenterWindow: centerWindow, Background: enableBackground},
		Styles:             styles,
		Screen:             ScreenMenu,
		EasterEgg:          loc.RandomEasterEgg(),
		Players:            players,
		Enemys:             enemys,
		CurrentPlayer:      0,
		CurrentEnemy:       0,
		CursorX:            players[0].X,
		CursorY:            players[0].Y,
		Walls:              walls,
		Water:              water,
		FireTiles:          make(map[Point]int),
		SmokeTiles:         make(map[Point]int),
		keys:               newKeyMap(loc),
		menuKeys:           newMenuKeyMap(loc),
		help:               help.New(),
		Localizer:          loc,
		startPlayers:       playerCount,
		startEnemies:       enemysCount,
		startPlayerEffects: playerEffects,
		startEnemyEffects:  enemyEffects,
	}
}

func (m *Model) SetAvailableThemes() {
	if m.Theme != nil {
		m.AvailableThemes = m.Theme.TintIDs()
	}
}

func (m *Model) closestPlayer(ex, ey int) (int, int) {
	if len(m.Players) == 0 {
		return ex, ey
	}
	bestX, bestY := m.Players[0].X, m.Players[0].Y
	bestDist := utils.Abs(ex-bestX) + utils.Abs(ey-bestY)
	for _, pl := range m.Players[1:] {
		d := utils.Abs(ex-pl.X) + utils.Abs(ey-pl.Y)
		if d < bestDist {
			bestDist = d
			bestX, bestY = pl.X, pl.Y
		}
	}
	return bestX, bestY
}

func (m *Model) enemyOccupied(x, y, skipIdx int) bool {
	for i, e := range m.Enemys {
		if i != skipIdx && e.X == x && e.Y == y {
			return true
		}
	}
	for _, p := range m.Players {
		if p.X == x && p.Y == y {
			return true
		}
	}
	return false
}

func (m *Model) doEnemyTurn(idx int) *Model {
	if len(m.Players) == 0 || idx >= len(m.Enemys) {
		return m
	}

	for step := 0; step < moveRange; step++ {
		en := m.Enemys[idx]
		tx, ty := m.closestPlayer(en.X, en.Y)
		dist := utils.Abs(en.X-tx) + utils.Abs(en.Y-ty)

		if dist <= shootRange && !m.HasWallBetweenPoints(en.X, en.Y, tx, ty) {
			for j, pl := range m.Players {
				if pl.X == tx && pl.Y == ty {
					m.Players[j].HP--
					m.BoxTrigger = TriggerDamage
					m.TriggerTimer = 6
					if m.Players[j].HP <= 0 {
						m.Players = append(m.Players[:j], m.Players[j+1:]...)
						if m.CurrentPlayer >= len(m.Players) {
							m.CurrentPlayer = 0
						}
					}
					break
				}
			}
			return m
		}

		moves := []Point{}
		if tx > en.X {
			moves = append(moves, Point{en.X + 1, en.Y})
		}
		if tx < en.X {
			moves = append(moves, Point{en.X - 1, en.Y})
		}
		if ty > en.Y {
			moves = append(moves, Point{en.X, en.Y + 1})
		}
		if ty < en.Y {
			moves = append(moves, Point{en.X, en.Y - 1})
		}

		moved := false
		for _, mv := range moves {
			if mv.X < 0 || mv.X >= GridW || mv.Y < 0 || mv.Y >= GridH {
				continue
			}
			if m.Walls[mv] || m.enemyOccupied(mv.X, mv.Y, idx) {
				continue
			}
			m.Enemys[idx].X = mv.X
			m.Enemys[idx].Y = mv.Y

			p := Point{mv.X, mv.Y}
			if m.FireTiles[p] > 0 && !HasEffect(m.Enemys[idx].Effects, EffectWet) {
				m.Enemys[idx].Effects = ResolveEffects(
					m.Enemys[idx].Effects,
					Effect{Type: EffectFire, Duration: 2},
				)
			}
			if m.Water[p] {
				m.Enemys[idx].Effects = ResolveEffects(
					m.Enemys[idx].Effects,
					Effect{Type: EffectWet, Duration: 2},
				)
			}

			moved = true
			break
		}
		if !moved {
			break
		}
	}

	if HasEffect(m.Enemys[idx].Effects, EffectFire) && m.Enemys[idx].HP > 1 {
		m.Enemys[idx].HP--
	}

	m.Enemys[idx].Effects = TickEffects(m.Enemys[idx].Effects)
	return m
}

func (m *Model) Move(newX, newY int) *Model {
	return m
}

func (m *Model) currentRange() int {
	r := moveRange
	if m.ShootMode {
		return shootRange
	}
	if len(m.Players) > 0 && m.CurrentPlayer < len(m.Players) {
		if HasEffect(m.Players[m.CurrentPlayer].Effects, EffectWet) {
			r -= 2
		}
	}
	if r < 1 {
		r = 1
	}
	return r
}

func (m *Model) IsInRange(col, row int) bool {
	if len(m.Players) == 0 || m.CurrentPlayer >= len(m.Players) {
		return false
	}
	current := m.Players[m.CurrentPlayer]
	dx := utils.Abs(col - current.X)
	dy := utils.Abs(row - current.Y)
	r := m.currentRange()
	if dx+dy > r || dx+dy == 0 {
		return false
	}
	return !m.HasWallBetweenPoints(current.X, current.Y, col, row)
}

func (m *Model) Reachable(sx, sy, r int) map[Point]bool {
	type state struct {
		x, y, steps int
	}
	visited := map[Point]bool{}
	result := map[Point]bool{}
	queue := []state{{sx, sy, 0}}
	visited[Point{sx, sy}] = true

	occupied := map[Point]bool{}
	for i, pl := range m.Players {
		if i != m.CurrentPlayer {
			occupied[Point{pl.X, pl.Y}] = true
		}
	}
	for _, en := range m.Enemys {
		occupied[Point{en.X, en.Y}] = true
	}

	dirs := []Point{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		for _, d := range dirs {
			nx, ny := cur.x+d.X, cur.y+d.Y
			np := Point{nx, ny}
			if nx < 0 || nx >= GridW || ny < 0 || ny >= GridH {
				continue
			}
			if visited[np] {
				continue
			}
			if m.Walls[np] {
				continue
			}
			visited[np] = true
			if cur.steps+1 <= r {
				result[np] = true
				if !occupied[np] {
					queue = append(queue, state{nx, ny, cur.steps + 1})
				}
			}
		}
	}
	return result
}

func (m *Model) HasWallBetweenPoints(x0, y0, x1, y1 int) bool {
	startX, startY := x0, y0
	dx := utils.Abs(x1 - x0)
	dy := utils.Abs(y1 - y0)
	sx := 1
	if x0 > x1 {
		sx = -1
	}
	sy := 1
	if y0 > y1 {
		sy = -1
	}
	err := dx - dy

	for {
		isStart := x0 == startX && y0 == startY
		isEnd := x0 == x1 && y0 == y1
		if !isStart && !isEnd {
			if m.Walls[Point{x0, y0}] {
				return true
			}
		}
		if isEnd {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x0 += sx
		}
		if e2 < dx {
			err += dx
			y0 += sy
		}
	}
	return false
}

func (m *Model) ultCross(cx, cy int) []Point {
	offsets := []Point{
		{0, 0}, {1, 0}, {-1, 0}, {0, 1}, {0, -1},
	}
	var pts []Point
	for _, o := range offsets {
		p := Point{cx + o.X, cy + o.Y}
		if p.X < 0 || p.X >= GridW || p.Y < 0 || p.Y >= GridH {
			continue
		}
		if m.Walls[p] {
			continue
		}
		pts = append(pts, p)
	}
	return pts
}

func (m *Model) ultInAxisRange(cx, cy int) bool {
	current := m.Players[m.CurrentPlayer]
	return cx == current.X || cy == current.Y
}

func (m *Model) doConfirm() *Model {
	p := Point{m.CursorX, m.CursorY}
	current := m.Players[m.CurrentPlayer]
	wallBlocked := m.HasWallBetweenPoints(current.X, current.Y, m.CursorX, m.CursorY)

	if m.UltMode && !m.Shot {
		m = m.doUlt()
		cur := m.Players[m.CurrentPlayer]
		m.CursorX = cur.X
		m.CursorY = cur.Y
	} else if m.ShootMode && !m.Shot {
		if HasEffect(m.Players[m.CurrentPlayer].Effects, EffectSmoke) {
			m.Shot = true
			m.ShootMode = false
			return m
		}
		if m.IsInRange(m.CursorX, m.CursorY) && !m.Walls[p] && !m.HasWallBetweenPoints(current.X, current.Y, m.CursorX, m.CursorY) {
			var hit bool
			for i, pl := range m.Players {
				if i != m.CurrentPlayer && pl.X == m.CursorX && pl.Y == m.CursorY {
					m.Players[i].HP--
					m.BoxTrigger = TriggerDamage
					m.TriggerTimer = 6
					if m.Players[i].HP <= 0 {
						m.CurrentScore -= 5
						m.Players = append(m.Players[:i], m.Players[i+1:]...)
						if m.CurrentPlayer >= len(m.Players) {
							m.CurrentPlayer = 0
						}
					}
					hit = true
					break
				}
			}
			if !hit {
				for i, en := range m.Enemys {
					if en.X == m.CursorX && en.Y == m.CursorY {
						m.Enemys[i].HP--
						m.BoxTrigger = TriggerDamage
						m.TriggerTimer = 6
						if m.Enemys[i].HP <= 0 {
							m.CurrentScore += 10
							m.Enemys = append(m.Enemys[:i], m.Enemys[i+1:]...)
						}
						hit = true
						break
					}
				}
			}
			if hit {
				m.CurrentScore++
			}
			m.Shot = true
			m.ShootMode = false
			cur := m.Players[m.CurrentPlayer]
			m.CursorX = cur.X
			m.CursorY = cur.Y
		}
	} else if !m.ShootMode && !m.UltMode && !m.Moved {
		if m.IsInRange(m.CursorX, m.CursorY) && !m.Walls[p] && !wallBlocked && !m.OccupiedByOther(m.CursorX, m.CursorY) {
			m.Players[m.CurrentPlayer].X = m.CursorX
			m.Players[m.CurrentPlayer].Y = m.CursorY

			if m.Water[p] {
				m.Players[m.CurrentPlayer].Effects = ResolveEffects(
					m.Players[m.CurrentPlayer].Effects,
					Effect{Type: EffectWet, Duration: 2},
				)
			}
			if m.FireTiles[p] > 0 {
				if !HasEffect(m.Players[m.CurrentPlayer].Effects, EffectWet) {
					m.Players[m.CurrentPlayer].Effects = ResolveEffects(
						m.Players[m.CurrentPlayer].Effects,
						Effect{Type: EffectFire, Duration: 2},
					)
				}
			}
			if m.SmokeTiles[p] > 0 {
				m.Players[m.CurrentPlayer].Effects = ResolveEffects(
					m.Players[m.CurrentPlayer].Effects,
					Effect{Type: EffectSmoke, Duration: 2},
				)
			}

			m.Moved = true
			m.CursorX = m.Players[m.CurrentPlayer].X
			m.CursorY = m.Players[m.CurrentPlayer].Y
		}
	}

	return m
}

func (m *Model) doUlt() *Model {
	current := m.Players[m.CurrentPlayer]
	if current.UltCharges <= 0 {
		return m
	}
	if !m.ultInAxisRange(m.CursorX, m.CursorY) {
		return m
	}

	m.Players[m.CurrentPlayer].UltCharges--
	m.UltMode = false
	m.UltAxis = ""
	m.Shot = true
	m.CurrentScore += 2

	affected := m.ultCross(m.CursorX, m.CursorY)

	affectedSet := make(map[Point]bool, len(affected))
	for _, p := range affected {
		affectedSet[p] = true
	}

	for _, p := range affected {
		if m.Water[p] || m.SmokeTiles[p] > 0 {
			m.SmokeTiles[p] = 2
		} else {
			m.FireTiles[p] = 2
		}
	}

	for i, pl := range m.Players {
		p := Point{pl.X, pl.Y}
		if !affectedSet[p] {
			continue
		}
		if m.SmokeTiles[p] > 0 {
			m.Players[i].Effects = ResolveEffects(
				m.Players[i].Effects,
				Effect{Type: EffectSmoke, Duration: 2},
			)
		} else if m.FireTiles[p] > 0 {
			if HasEffect(pl.Effects, EffectWet) {
				m.Players[i].Effects = RemoveEffect(m.Players[i].Effects, EffectWet)
			} else {
				m.Players[i].Effects = ResolveEffects(
					m.Players[i].Effects,
					Effect{Type: EffectFire, Duration: 2},
				)
			}
		}
	}

	for i, en := range m.Enemys {
		p := Point{en.X, en.Y}
		if !affectedSet[p] {
			continue
		}
		if m.SmokeTiles[p] > 0 {
			m.Enemys[i].Effects = ResolveEffects(
				m.Enemys[i].Effects,
				Effect{Type: EffectSmoke, Duration: 2},
			)
		} else if m.FireTiles[p] > 0 {
			if HasEffect(en.Effects, EffectWet) {
				m.Enemys[i].Effects = RemoveEffect(m.Enemys[i].Effects, EffectWet)
			} else {
				m.Enemys[i].Effects = ResolveEffects(
					m.Enemys[i].Effects,
					Effect{Type: EffectFire, Duration: 2},
				)
			}
		}
	}

	return m
}

func (m *Model) tickFireTiles() *Model {
	for p, t := range m.FireTiles {
		t--
		if t <= 0 {
			delete(m.FireTiles, p)
		} else {
			m.FireTiles[p] = t
		}
	}
	for p, t := range m.SmokeTiles {
		t--
		if t <= 0 {
			delete(m.SmokeTiles, p)
		} else {
			m.SmokeTiles[p] = t
		}
	}
	return m
}

func (m *Model) nextTurn() *Model {
	m.Moved = false
	m.Shot = false
	m.ShootMode = false
	m.UltMode = false
	m.UltAxis = ""

	if HasEffect(m.Players[m.CurrentPlayer].Effects, EffectFire) && m.Players[m.CurrentPlayer].HP > 1 {
		m.Players[m.CurrentPlayer].HP--
		m.BoxTrigger = TriggerDamage
		m.TriggerTimer = 6
	}

	m.Players[m.CurrentPlayer].Effects = TickEffects(
		m.Players[m.CurrentPlayer].Effects,
	)

	p := Point{m.Players[m.CurrentPlayer].X, m.Players[m.CurrentPlayer].Y}
	if m.Water[p] {
		m.Players[m.CurrentPlayer].Effects = ResolveEffects(
			m.Players[m.CurrentPlayer].Effects,
			Effect{Type: EffectWet, Duration: 2},
		)
	}

	if m.CurrentPlayer == len(m.Players)-1 {
		m = m.tickFireTiles()
		for i := range m.Players {
			m.Players[i].Effects = TickEffects(m.Players[i].Effects)
		}
	}

	m.CurrentPlayer = (m.CurrentPlayer + 1) % len(m.Players)
	m.CurrentScore++
	next := m.Players[m.CurrentPlayer]
	m.CursorX = next.X
	m.CursorY = next.Y
	return m
}

func enemyTurnCmd(idx int) tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return enemyTurnMsg{enemyIdx: idx}
	})
}

func (m *Model) OccupiedByOther(x, y int) bool {
	for i, p := range m.Players {
		if i != m.CurrentPlayer && p.X == x && p.Y == y {
			return true
		}
	}
	for _, e := range m.Enemys {
		if e.X == x && e.Y == y {
			return true
		}
	}
	return false
}

func (m *Model) turnOrder() string {
	var parts []string

	for i, pl := range m.Players {
		symbol := " ■ "
		style := pl.Style
		if i == m.CurrentPlayer && !m.EnemyTurn {
			style = style.Underline(true).Bold(true)
			symbol = " ● "
		}
		parts = append(parts, style.Render(symbol))
	}

	parts = append(parts, lipgloss.NewStyle().
		Foreground(m.Theme.SelectionBg()).Render(" · "))

	for i, en := range m.Enemys {
		symbol := " ▲ "
		style := en.Style
		if m.EnemyTurn && i == m.EnemyIdx {
			style = style.Underline(true).Bold(true)
			symbol = " ♦ "
		}
		parts = append(parts, style.Render(symbol))
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, parts...)
}

// endGame saves the current score to the active profile, persists profiles
// to disk, and returns the UI to the main menu. scoreBonus is added on top
// of CurrentScore (e.g. win/loss bonus).
func (m *Model) endGame(scoreBonus int) {
	score := m.CurrentScore + scoreBonus
	if score < 0 {
		score = 0
	}
	if m.ActiveSlot >= 0 && m.ActiveSlot < len(m.Profiles) && m.Profiles[m.ActiveSlot] != nil {
		m.Profiles[m.ActiveSlot].Score += score
	}
	_ = m.SaveProfiles()
	m.ActiveSlot = -1
	if len(m.Players) == 0 {
		m.Screen = ScreenGameOver
	} else {
		m.Screen = ScreenMenu
	}
	m.MenuSelected = 0
}

func (m *Model) doMenuConfirm() (tea.Model, tea.Cmd) {
	if m.Screen == ScreenMenu {
		switch m.MenuSelected {
		case 0:
			m.Screen = ScreenProfiles
			m.ProfileSlot = 0
		case 1:
			m.Screen = ScreenSettings
			m.MenuSelected = 0
		case 2:
			return m, tea.Quit
		}
	} else if m.Screen == ScreenSettings {
		hasLang := len(m.Localizer.AvailableLanguages()) > 1
		hasThemes := len(m.AvailableThemes) > 1
		n := 0
		if hasLang {
			if m.MenuSelected == n {
				languages := m.Localizer.AvailableLanguages()
				currentIdx := 0
				for i, l := range languages {
					if l == m.Localizer.GetLanguage() {
						currentIdx = i
						break
					}
				}
				nextIdx := (currentIdx + 1) % len(languages)
				if err := m.SetLanguage(languages[nextIdx]); err != nil {
					return m, nil
				}
				return m, nil
			}
			n++
		}
		if hasThemes {
			if m.MenuSelected == n {
				m.Screen = ScreenThemeSelect
				m.MenuSelected = 0
				return m, nil
			}
			n++
		}
		if m.MenuSelected == n {
			m.Screen = ScreenCenterSelect
			m.CenterRow, m.CenterCol = centerToGrid(m.Config.CenterWindow)
			return m, nil
		}
		n++
		if m.MenuSelected == n {
			if m.Config != nil {
				m.Config.Background = !m.Config.Background
				if err := m.SaveConfig(); err != nil {
					m.Config.Background = !m.Config.Background
				}
			}
			return m, nil
		}
		m.Screen = ScreenMenu
		m.MenuSelected = 0
	} else if m.Screen == ScreenThemeSelect {
		m.Screen = ScreenSettings
		m.MenuSelected = 0
	} else if m.Screen == ScreenProfiles {
		slot := m.ProfileSlot
		if m.Profiles[slot] != nil {
			m.ActiveSlot = slot
			m.CurrentScore = 0
			m.Screen = ScreenGame
			m.startGame()
		} else {
			m.ProfileLetters = [3]rune{'A', 'A', 'A'}
			m.ProfileCursor = 0
			m.Screen = ScreenProfileCreate
		}
	}
	return m, nil
}

func centerToGrid(mode string) (int, int) {
	switch mode {
	case "tl":
		return 0, 0
	case "tc":
		return 0, 1
	case "tr":
		return 0, 2
	case "cl":
		return 1, 0
	case "c":
		return 1, 1
	case "cr":
		return 1, 2
	case "bl":
		return 2, 0
	case "bc":
		return 2, 1
	case "br":
		return 2, 2
	default:
		return 1, 1
	}
}

func gridToCenter(row, col int) string {
	grid := [3][3]string{
		{"tl", "tc", "tr"},
		{"cl", "c", "cr"},
		{"bl", "bc", "br"},
	}
	return grid[row][col]
}
