package generate

import (
	"time"

	"hera/i18n"
	"hera/utils"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	bubbletint "github.com/lrstanley/bubbletint"
	bz "github.com/lrstanley/bubblezone"
)

func NewModel(playerCount, enemysCount int, loc i18n.Localizer, theme *bubbletint.Registry) Model {
	styles := NewStyles(theme)
	if playerCount < 2 {
		playerCount = 2
	}
	if playerCount > 4 {
		playerCount = 4
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
		players[i] = Player{
			X:          starts[i].X,
			Y:          starts[i].Y,
			HP:         MaxHP,
			UltCharges: maxUltCharges,
			Style:      styles.PlayerStyles[i],
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
		enemys[i] = Enemy{
			X:     enemyPositions[i].X,
			Y:     enemyPositions[i].Y,
			HP:    MaxHP,
			Style: styles.EnemysStyles[i],
		}
	}

	return Model{
		Theme:         theme,
		Styles:        styles,
		Players:       players,
		Enemys:        enemys,
		CurrentPlayer: 0,
		CurrentEnemy:  0,
		CursorX:       players[0].X,
		CursorY:       players[0].Y,
		Walls:         walls,
		Water:         water,
		FireTiles:     make(map[Point]int),
		SmokeTiles:    make(map[Point]int),
		keys:          newKeyMap(loc),
		help:          help.New(),
		Localizer:     loc,
		Z:             bz.New(),
	}
}

func (m Model) closestPlayer(ex, ey int) (int, int) {
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

func (m Model) enemyOccupied(x, y, skipIdx int) bool {
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

func (m Model) doEnemyTurn(idx int) Model {
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
			if m.FireTiles[p] > 0 && !hasEffect(m.Enemys[idx].Effects, EffectWet) {
				m.Enemys[idx].Effects = resolveEffects(
					m.Enemys[idx].Effects,
					Effect{Type: EffectFire, Duration: 2},
				)
			}
			if m.Water[p] {
				m.Enemys[idx].Effects = resolveEffects(
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

	if hasEffect(m.Enemys[idx].Effects, EffectFire) && m.Enemys[idx].HP > 1 {
		m.Enemys[idx].HP--
	}

	m.Enemys[idx].Effects = tickEffects(m.Enemys[idx].Effects)
	return m
}

func (m Model) Move(newX, newY int) Model {
	return m
}

func (m Model) currentRange() int {
	r := moveRange
	if m.ShootMode {
		return shootRange
	}
	if len(m.Players) > 0 && m.CurrentPlayer < len(m.Players) {
		if hasEffect(m.Players[m.CurrentPlayer].Effects, EffectWet) {
			r -= 2
		}
	}
	if r < 1 {
		r = 1
	}
	return r
}

func (m Model) IsInRange(col, row int) bool {
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

func (m Model) Reachable(sx, sy, r int) map[Point]bool {
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

func (m Model) HasWallBetweenPoints(x0, y0, x1, y1 int) bool {
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

func (m Model) ultCross(cx, cy int) []Point {
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

func (m Model) ultInAxisRange(cx, cy int) bool {
	current := m.Players[m.CurrentPlayer]
	return cx == current.X || cy == current.Y
}

func (m Model) doUlt() Model {
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
			m.Players[i].Effects = resolveEffects(
				m.Players[i].Effects,
				Effect{Type: EffectSmoke, Duration: 2},
			)
		} else if m.FireTiles[p] > 0 {
			if hasEffect(pl.Effects, EffectWet) {
				m.Players[i].Effects = removeEffect(m.Players[i].Effects, EffectWet)
			} else {
				m.Players[i].Effects = resolveEffects(
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
			m.Enemys[i].Effects = resolveEffects(
				m.Enemys[i].Effects,
				Effect{Type: EffectSmoke, Duration: 2},
			)
		} else if m.FireTiles[p] > 0 {
			if hasEffect(en.Effects, EffectWet) {
				m.Enemys[i].Effects = removeEffect(m.Enemys[i].Effects, EffectWet)
			} else {
				m.Enemys[i].Effects = resolveEffects(
					m.Enemys[i].Effects,
					Effect{Type: EffectFire, Duration: 2},
				)
			}
		}
	}

	return m
}

func (m Model) tickFireTiles() Model {
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

func (m Model) nextTurn() Model {
	m.Moved = false
	m.Shot = false
	m.ShootMode = false
	m.UltMode = false
	m.UltAxis = ""

	if hasEffect(m.Players[m.CurrentPlayer].Effects, EffectFire) && m.Players[m.CurrentPlayer].HP > 1 {
		m.Players[m.CurrentPlayer].HP--
	}

	m.Players[m.CurrentPlayer].Effects = tickEffects(
		m.Players[m.CurrentPlayer].Effects,
	)

	p := Point{m.Players[m.CurrentPlayer].X, m.Players[m.CurrentPlayer].Y}
	if m.Water[p] {
		m.Players[m.CurrentPlayer].Effects = resolveEffects(
			m.Players[m.CurrentPlayer].Effects,
			Effect{Type: EffectWet, Duration: 2},
		)
	}

	if m.CurrentPlayer == len(m.Players)-1 {
		m = m.tickFireTiles()
		for i := range m.Players {
			m.Players[i].Effects = tickEffects(m.Players[i].Effects)
		}
	}

	m.CurrentPlayer = (m.CurrentPlayer + 1) % len(m.Players)
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

func (m Model) OccupiedByOther(x, y int) bool {
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

func (m Model) turnOrder() string {
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
		Foreground(lipgloss.Color("#444444")).Render(" · "))

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
