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
)

var (
	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#874BFD")).
			Padding(1, 2).
			Width(32)
	helpStyle = lipgloss.NewStyle().
			Padding(1, 2)

	playerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF75B7")).
			Bold(true)
	cellStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#555555"))
	wallStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#874BFD"))
	waterStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#146fba"))
)

type Point struct {
	X, Y int
}

type keyMap struct {
	Up    key.Binding
	Down  key.Binding
	Left  key.Binding
	Right key.Binding
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

func NewModel() Model {
	startX, startY := GridW/2, GridH/2

	walls := GenerateTiles(startX, startY, wallCount, nil)
	water := GenerateTiles(startX, startY, waterCount, walls)

	return Model{
		X:     startX,
		Y:     startY,
		Walls: walls,
		Water: water,
		keys:  keys,
		help:  help.New(),
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
	Y, X  int
	Walls map[Point]bool
	Water map[Point]bool
	keys  keyMap
	help  help.Model
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
	newX = utils.Clamp(newX, 0, GridW-1)
	newY = utils.Clamp(newY, 0, GridH-1)
	if !m.Walls[Point{newX, newY}] && !m.Water[Point{newX, newY}] {
		m.X = newX
		m.Y = newY
	}
	return m
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			m = m.Move(m.X, m.Y-1)
		case "down", "j":
			m = m.Move(m.X, m.Y+1)
		case "left", "h":
			m = m.Move(m.X-1, m.Y)
		case "right", "l":
			m = m.Move(m.X+1, m.Y)
		case "?":
			m.help.ShowAll = !m.help.ShowAll
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) View() string {
	var rows []string
	for row := 0; row < GridH; row++ {
		var cells []string
		for col := 0; col < GridW; col++ {
			switch {
			case col == m.X && row == m.Y:
				cells = append(cells, playerStyle.Render(""))
			case m.Walls[Point{col, row}]:
				cells = append(cells, wallStyle.Render("▪"))
			case m.Water[Point{col, row}]:
				cells = append(cells, waterStyle.Render("≈"))
			default:
				cells = append(cells, cellStyle.Render("·"))
			}
		}
		rows = append(rows, strings.Join(cells, " "))
	}
	grid := strings.Join(rows, "\n")

	box := boxStyle.Render(grid)
	helpView := helpStyle.Render(m.help.View(m.keys))

	return lipgloss.JoinVertical(lipgloss.Left, box, helpView)
}
