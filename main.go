package main

import (
	"strings"
	"math/rand"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
)

const (
	gridSize = 10
	wallCount = 10
	waterCount = 10
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(1, 2).
			Width(25)
	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#874BFD")).
			Padding(1, 2).
			Width(23)
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

type point struct {
	x, y int
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

func newModel() model {
	startX, startY := gridSize/2, gridSize/2
	return model{
		x: 10 / 2,
		y: 10 / 2,
		walls: generateWalls(startX, startY),
		water: generateWater(startX, startY),
		keys:       keys,
		help:       help.New(),
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

type model struct {
	x, y int
	width, height int
	walls map[point]bool
	water map[point]bool
	keys keyMap
	help help.Model
}


func generateWater(PlayerX, PlayerY int) map[point]bool {
	water := make(map[point]bool)
	for len(water) < wallCount {
		x := rand.Intn(gridSize)
		y := rand.Intn(gridSize)
		p := point{x, y}

		if abs(x-PlayerX) <= 1 && abs(y-PlayerY) <= 1 {
			continue
		}
		water[p] = true
	}
	return water
}

func generateWalls(PlayerX, PlayerY int) map[point]bool {
	walls := make(map[point]bool)
	for len(walls) < wallCount {
		x := rand.Intn(gridSize)
		y := rand.Intn(gridSize)
		p := point{x, y}

		if abs(x-PlayerX) <= 1 && abs(y-PlayerY) <= 1 {
			continue
		}
		walls[p] = true
	}
	return walls
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (m model) Init() tea.Cmd {
	return nil
}

func clamp(val, min, max int) int {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

func (m model) Move(newX, newY int) model {
	newX = clamp(newX, 0, gridSize-1)
	newY = clamp(newY, 0, gridSize-1)
	if !m.walls[point{newX, newY}] && !m.water[point{newX, newY}] {
		m.x = newX
		m.y = newY
	}
	return m
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			m = m.Move(m.x, m.y-1)
		case "down", "j":
			m = m.Move(m.x, m.y+1)
		case "left", "h":
			m = m.Move(m.x-1, m.y)
		case "right", "l":
			m = m.Move(m.x+1, m.y)
		case "?":
			m.help.ShowAll = !m.help.ShowAll
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	var rows []string
	for row := 0; row < 10; row++ {
		var cells []string
		for col := 0; col < 10; col++ {
			switch {
			case col == m.x && row == m.y:
				cells = append(cells, playerStyle.Render(""))
			case m.walls[point{col, row}]:
				cells = append(cells, wallStyle.Render("▪"))
			case m.water[point{col, row}]:
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

func main() {
	p := tea.NewProgram(newModel())
	if _, err := p.Run(); err != nil {
		log.Error(err)
	}
}
