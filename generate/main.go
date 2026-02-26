package generate

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"hera/utils"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	GridW      = 14
	GridH      = 10
	wallCount  = 10
	waterCount = 10
	moveRange  = 4
	shootRange = 2
	MaxHP      = 3
)

var (
	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#874BFD")).
			Padding(1, 2).
			Width(46)
	helpStyle = lipgloss.NewStyle().
			Padding(1, 2)
	cellStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#555555"))
	wallStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#874BFD"))
	waterStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#146fba"))
	cursorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#222222"))
	rangeStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#171717"))
)

var playerStyles = []lipgloss.Style{
	lipgloss.NewStyle().Foreground(lipgloss.Color("#FF75B7")).Bold(true),
	lipgloss.NewStyle().Foreground(lipgloss.Color("#75FFBA")).Bold(true),
	lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD700")).Bold(true),
	lipgloss.NewStyle().Foreground(lipgloss.Color("#75B7FF")).Bold(true),
}

var enemysStyles = []lipgloss.Style{
	lipgloss.NewStyle().Foreground(lipgloss.Color("#ffcbe4")).Bold(true),
	lipgloss.NewStyle().Foreground(lipgloss.Color("#a2ffd0")).Bold(true),
	lipgloss.NewStyle().Foreground(lipgloss.Color("#ffdb1e")).Bold(true),
	lipgloss.NewStyle().Foreground(lipgloss.Color("#84befe")).Bold(true),
}

type Point struct {
	X, Y int
}

type Player struct {
	X, Y  int
	HP    int
	Style lipgloss.Style
}

type Enemy struct {
	X, Y  int
	HP    int
	Style lipgloss.Style
}

type enemyTurnMsg struct {
	enemyIdx int
}

type keyMap struct {
	Up      key.Binding
	Down    key.Binding
	Left    key.Binding
	Right   key.Binding
	Confirm key.Binding
	Shoot   key.Binding
	Help    key.Binding
	Quit    key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},
		{k.Confirm, k.Shoot, k.Help, k.Quit},
	}
}

func NewModel(playerCount, enemysCount int) Model {
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
			X:     starts[i].X,
			Y:     starts[i].Y,
			HP:    MaxHP,
			Style: playerStyles[i],
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
			Style: enemysStyles[i],
		}
	}

	return Model{
		Players:       players,
		Enemys:        enemys,
		CurrentPlayer: 0,
		CurrentEnemy:  0,
		CursorX:       players[0].X,
		CursorY:       players[0].Y,
		Walls:         walls,
		Water:         water,
		keys:          keys,
		help:          help.New(),
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

	en := m.Enemys[idx]
	tx, ty := m.closestPlayer(en.X, en.Y)

	for step := 0; step < moveRange; step++ {
		en = m.Enemys[idx]
		tx, ty = m.closestPlayer(en.X, en.Y)
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
			moved = true
			break
		}
		if !moved {
			break
		}
	}
	return m
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/K", "Move Up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/J", "Move Down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/H", "Move Left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/L", "Move Right"),
	),
	Confirm: key.NewBinding(
		key.WithKeys("x", "enter"),
		key.WithHelp("X/Enter", "Move Player"),
	),
	Shoot: key.NewBinding(
		key.WithKeys("z", "shift"),
		key.WithHelp("Z/Shift", "Change Mode"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "Toggle Help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "Quit"),
	),
}

type Model struct {
	Players          []Player
	Enemys           []Enemy
	CurrentPlayer    int
	CurrentEnemy     int
	CursorX, CursorY int
	Walls            map[Point]bool
	Water            map[Point]bool
	ShootMode        bool
	Moved            bool
	Shot             bool
	keys             keyMap
	help             help.Model
	EnemyTurn        bool
	EnemyIdx         int
}

func GenerateTiles(playerX, playerY, count int, blocked map[Point]bool) map[Point]bool {
	tiles := make(map[Point]bool)
	for len(tiles) < count {
		x := rand.Intn(GridW)
		y := rand.Intn(GridH)
		p := Point{x, y}

		if utils.Abs(x-playerX) <= 1 && utils.Abs(y-playerY) <= 1 {
			continue
		}
		if blocked[p] {
			continue
		}
		tiles[p] = true
	}
	return tiles
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Move(newX, newY int) Model {
	return m
}

func (m Model) currentRange() int {
	if m.ShootMode {
		return shootRange
	}
	return moveRange
}

func (m Model) IsInRange(col, row int) bool {
	current := m.Players[m.CurrentPlayer]
	dx := utils.Abs(col - current.X)
	dy := utils.Abs(row - current.Y)
	r := m.currentRange()
	if dx+dy > r || dx+dy == 0 {
		return false
	}
	return !m.HasWallBetweenPoints(current.X, current.Y, col, row)
}

func (m Model) inRange(x, y int) bool {
	current := m.Players[m.CurrentPlayer]
	dx := utils.Abs(x - current.X)
	dy := utils.Abs(y - current.Y)
	return dx+dy <= m.currentRange()
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

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !m.Moved && !m.Shot {
	} else if m.Moved {
		m.ShootMode = true
	} else {
		m.ShootMode = false
	}

	switch msg := msg.(type) {
	case enemyTurnMsg:
		if len(m.Players) == 0 {
			return m, tea.Quit
		}
		if msg.enemyIdx >= len(m.Enemys) {
			m.EnemyTurn = false
			m.EnemyIdx = 0
			return m, nil
		}
		m.EnemyIdx = msg.enemyIdx
		m = m.doEnemyTurn(msg.enemyIdx)
		if len(m.Players) == 0 {
			return m, tea.Quit
		}
		return m, enemyTurnCmd(msg.enemyIdx + 1)

	case tea.KeyMsg:
		if m.EnemyTurn {
			return m, nil
		}

		switch {
		case key.Matches(msg, m.keys.Up):
			newY := utils.Clamp(m.CursorY-1, 0, GridH-1)
			if m.inRange(m.CursorX, newY) {
				m.CursorY = newY
			}
		case key.Matches(msg, m.keys.Down):
			newY := utils.Clamp(m.CursorY+1, 0, GridH-1)
			if m.inRange(m.CursorX, newY) {
				m.CursorY = newY
			}
		case key.Matches(msg, m.keys.Left):
			newX := utils.Clamp(m.CursorX-1, 0, GridW-1)
			if m.inRange(newX, m.CursorY) {
				m.CursorX = newX
			}
		case key.Matches(msg, m.keys.Right):
			newX := utils.Clamp(m.CursorX+1, 0, GridW-1)
			if m.inRange(newX, m.CursorY) {
				m.CursorX = newX
			}
		case key.Matches(msg, m.keys.Shoot):
			if !m.Shot {
				m.ShootMode = !m.ShootMode
				current := m.Players[m.CurrentPlayer]
				m.CursorX = current.X
				m.CursorY = current.Y
			}
		case key.Matches(msg, m.keys.Confirm):
			p := Point{m.CursorX, m.CursorY}
			current := m.Players[m.CurrentPlayer]
			wallBlocked := m.HasWallBetweenPoints(current.X, current.Y, m.CursorX, m.CursorY)

			if m.ShootMode && !m.Shot {
				if !m.Walls[p] && !m.HasWallBetweenPoints(current.X, current.Y, m.CursorX, m.CursorY) {
					for i, pl := range m.Players {
						if i != m.CurrentPlayer && pl.X == m.CursorX && pl.Y == m.CursorY {
							m.Players[i].HP--
							if m.Players[i].HP <= 0 {
								m.Players = append(m.Players[:i], m.Players[i+1:]...)
								if m.CurrentPlayer >= len(m.Players) {
									m.CurrentPlayer = 0
								}
							}
							break
						}
					}
					for i, en := range m.Enemys {
						if en.X == m.CursorX && en.Y == m.CursorY {
							m.Enemys[i].HP--
							if m.Enemys[i].HP <= 0 {
								m.Enemys = append(m.Enemys[:i], m.Enemys[i+1:]...)
							}
							break
						}
					}
					m.Shot = true
					m.ShootMode = false
					cur := m.Players[m.CurrentPlayer]
					m.CursorX = cur.X
					m.CursorY = cur.Y
				}
			} else if !m.ShootMode && !m.Moved {
				if !m.Walls[p] && !wallBlocked && !m.OccupiedByOther(m.CursorX, m.CursorY) {
					m.Players[m.CurrentPlayer].X = m.CursorX
					m.Players[m.CurrentPlayer].Y = m.CursorY
					m.Moved = true
					m.CursorX = m.Players[m.CurrentPlayer].X
					m.CursorY = m.Players[m.CurrentPlayer].Y
				}
			}

			if m.Moved && m.Shot {
				m = m.nextTurn()
				if m.CurrentPlayer == 0 {
					m.EnemyTurn = true
					return m, enemyTurnCmd(0)
				}
			}

		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) nextTurn() Model {
	m.Moved = false
	m.Shot = false
	m.ShootMode = false
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

func (m Model) cursorInfo() string {
	p := Point{m.CursorX, m.CursorY}
	current := m.Players[m.CurrentPlayer]
	wallBlocked := m.HasWallBetweenPoints(current.X, current.Y, m.CursorX, m.CursorY)

	for i, pl := range m.Players {
		if pl.X == m.CursorX && pl.Y == m.CursorY {
			hp := strings.Repeat("♥ ", pl.HP) + strings.Repeat("♡ ", MaxHP-pl.HP)
			if i == m.CurrentPlayer {
				return pl.Style.Render(fmt.Sprintf("● Player %s (you)", hp))
			}
			if wallBlocked {
				return lipgloss.NewStyle().Foreground(lipgloss.Color("#FF4444")).
					Render(fmt.Sprintf("■ Player %d — wall in the way", i+1))
			}
			return pl.Style.Render(fmt.Sprintf("■ Player %s", hp))
		}
	}
	for i, en := range m.Enemys {
		if en.X == m.CursorX && en.Y == m.CursorY {
			hp := strings.Repeat("♥ ", en.HP) + strings.Repeat("♡ ", MaxHP-en.HP)
			if wallBlocked {
				return lipgloss.NewStyle().Foreground(lipgloss.Color("#FF4444")).
					Render(fmt.Sprintf("▲ Enemy %d — wall in the way", i+1))
			}
			return en.Style.Render(fmt.Sprintf("▲ Enemy %d %s", i+1, hp))
		}
	}

	switch {
	case m.Walls[p]:
		return wallStyle.Render("■ Wall — impassable")
	case wallBlocked:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#FF4444")).
			Render("⊘ Wall in the way")
	case m.Water[p]:
		return waterStyle.Render("≈ Water — passable")
	case m.IsInRange(m.CursorX, m.CursorY):
		if m.ShootMode {
			return lipgloss.NewStyle().Foreground(lipgloss.Color("#FF4444")).Render("· In shoot range")
		}
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#AAAAAA")).Render("· In move range")
	default:
		return cellStyle.Render("· Empty — out of range")
	}
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

func (m Model) View() string {
	current := m.Players[m.CurrentPlayer]
	hp := strings.Repeat("♥ ", current.HP) + strings.Repeat("♡ ", MaxHP-current.HP)
	hpStyle := current.Style
	if current.HP == 1 {
		hpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).Bold(true).Blink(true)
	}
	hpStr := hpStyle.Render(fmt.Sprintf("Player %d  %s", m.CurrentPlayer+1, hp))

	var modeStr string
	if m.ShootMode {
		modeStr = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF4444")).
			Bold(true).
			Render("♡ S ")
	} else {
		modeStr = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#AAAAAA")).
			Render("♧ M ")
	}

	var rows []string
	for row := 0; row < GridH; row++ {
		var cells []string
		for col := 0; col < GridW; col++ {
			p := Point{col, row}
			playerIdx := -1
			enemyIdx := -1
			for i, pl := range m.Players {
				if pl.X == col && pl.Y == row {
					playerIdx = i
					break
				}
			}
			for i, en := range m.Enemys {
				if en.X == col && en.Y == row {
					enemyIdx = i
					break
				}
			}
			switch {
			case col == m.CursorX && row == m.CursorY:
				if playerIdx >= 0 {
					cells = append(cells, cursorStyle.Render(
						m.Players[playerIdx].Style.Render(" ■ "),
					))
				} else if enemyIdx >= 0 {
					cells = append(cells, cursorStyle.Render(
						m.Enemys[enemyIdx].Style.Render(" ▲ "),
					))
				} else {
					cells = append(cells, cursorStyle.Render(" · "))
				}
			case playerIdx >= 0:
				symbol := " ■ "
				if playerIdx == m.CurrentPlayer {
					symbol = " ● "
				}
				cells = append(cells, m.Players[playerIdx].Style.Render(symbol))
			case enemyIdx >= 0:
				symbol := " ▲ "
				if enemyIdx == m.CurrentEnemy {
					symbol = " ♦ "
				}
				cells = append(cells, m.Enemys[enemyIdx].Style.Render(symbol))
			case m.Walls[p]:
				cells = append(cells, wallStyle.Render(" ■ "))
			case m.Water[p]:
				cells = append(cells, waterStyle.Render(" ≈ "))
			case m.IsInRange(col, row):
				cells = append(cells, rangeStyle.Render(" · "))
			default:
				cells = append(cells, cellStyle.Render(" · "))
			}
		}
		rows = append(rows, strings.Join(cells, ""))
	}

	info := m.cursorInfo()

	line0 := m.turnOrder()

	line1 := lipgloss.JoinHorizontal(lipgloss.Top,
		modeStr,
		" ",
		hpStr,
	)

	line2 := lipgloss.JoinHorizontal(lipgloss.Top,
		lipgloss.NewStyle().Foreground(lipgloss.Color("#666666")).Render(
			fmt.Sprintf("(%d, %d)  ", m.CursorX, m.CursorY),
		),
		info,
	)

	status := boxStyle.Render(fmt.Sprintf("%s\n%s\n%s", line1, line2, line0))
	grid := strings.Join(rows, "\n")
	box := boxStyle.Render(lipgloss.JoinVertical(lipgloss.Left, grid))
	helpView := helpStyle.Render(m.help.View(m.keys))
	return lipgloss.JoinVertical(lipgloss.Left, box, status, helpView)
}
