package generate

import (
	"math/rand"
	"strings"

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

type Point struct {
	X, Y int
}

type Player struct {
    X, Y  int
    Style lipgloss.Style
}

type keyMap struct {
	Up    key.Binding
	Down  key.Binding
	Left  key.Binding
	Right key.Binding
	Confirm key.Binding
	Help  key.Binding
	Quit  key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},
		{k.Help, k.Quit},
	}
}

func NewModel(playerCount int) Model {
    if playerCount < 2 { playerCount = 2 }
    if playerCount > 4 { playerCount = 4 }

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
            Style: playerStyles[i],
        }
    }

    return Model{
        Players:       players,
        CurrentPlayer: 0,
        CursorX:       players[0].X,
        CursorY:       players[0].Y,
        Walls:         walls,
        Water:         water,
        keys:          keys,
        help:          help.New(),
    }
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "move left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "move right"),
	),
	Confirm: key.NewBinding(
        	key.WithKeys("enter"),
        	key.WithHelp("enter", "move player"),
    	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

type Model struct {
    Players       []Player
    CurrentPlayer int
    CursorX, CursorY int
    Walls         map[Point]bool
    Water         map[Point]bool
    keys          keyMap
    help          help.Model
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

func (m Model) isInRange(col, row int) bool {
    current := m.Players[m.CurrentPlayer]
    dx := utils.Abs(col - current.X)
    dy := utils.Abs(row - current.Y)
    return dx+dy <= moveRange && dx+dy > 0
}

func (m Model) inRange(x, y int) bool {
    current := m.Players[m.CurrentPlayer]
    dx := utils.Abs(x - current.X)
    dy := utils.Abs(y - current.Y)
    result := dx+dy <= moveRange
    return result
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
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

        case key.Matches(msg, m.keys.Confirm):
            p := Point{m.CursorX, m.CursorY}
            if !m.Walls[p] && !m.Water[p] && !m.OccupiedByOther(m.CursorX, m.CursorY) {
                m.Players[m.CurrentPlayer].X = m.CursorX
                m.Players[m.CurrentPlayer].Y = m.CursorY
                m.CurrentPlayer = (m.CurrentPlayer + 1) % len(m.Players)
                next := m.Players[m.CurrentPlayer]
                m.CursorX = next.X
                m.CursorY = next.Y
            }
        case key.Matches(msg, m.keys.Help):
            m.help.ShowAll = !m.help.ShowAll
        case key.Matches(msg, m.keys.Quit):
            return m, tea.Quit
        }
    }
    return m, nil
}

func (m Model) OccupiedByOther(x, y int) bool {
    for i, p := range m.Players {
        if i != m.CurrentPlayer && p.X == x && p.Y == y {
            return true
        }
    }
    return false
}

func (m Model) View() string {
    var rows []string
    for row := 0; row < GridH; row++ {
        var cells []string
        for col := 0; col < GridW; col++ {
            p := Point{col, row}

            playerIdx := -1
            for i, pl := range m.Players {
                if pl.X == col && pl.Y == row {
                    playerIdx = i
                    break
                }
            }
	switch {
	case col == m.CursorX && row == m.CursorY:
	    if playerIdx >= 0 {
	        cells = append(cells, cursorStyle.Render(
	            m.Players[playerIdx].Style.Render(" ■ "),
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
	case m.Walls[p]:
	    cells = append(cells, wallStyle.Render(" ■ "))
	case m.Water[p]:
	    cells = append(cells, waterStyle.Render(" ≈ "))
	case m.isInRange(col, row) && !m.Walls[p] && !m.Water[p]:
	    cells = append(cells, rangeStyle.Render(" · "))
	default:
	    cells = append(cells, cellStyle.Render(" · "))
		}
        }
        rows = append(rows, strings.Join(cells, ""))
    }

    grid := strings.Join(rows, "\n")
    box := boxStyle.Render(lipgloss.JoinVertical(lipgloss.Left, grid))
    helpView := helpStyle.Render(m.help.View(m.keys))

    return lipgloss.JoinVertical(lipgloss.Left, box, helpView)
}
