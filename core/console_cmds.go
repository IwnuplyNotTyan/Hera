package generate

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

func helpText() string {
	return `Commands:
  regen                — Regenerate walls, water, fire/smoke
  add player           — Add a new player
  remove player <N>    — Remove player N
  add enemy            — Add a new enemy
  remove enemy <N>     — Remove enemy N
  heal <target> <N>    — Heal target by N HP
  damage <target> <N>  — Damage target by N HP
  effect <target> <type> [dur] — Add effect (fire/wet/smoke)
  remove effect <target> <type> — Remove effect
  set hp <target> <N>  — Set exact HP
  status               — Show all entities
  clear                — Clear console output
  help                 — This help`
}

// ExecuteConsoleCmd processes a parsed console command and returns its output string.
func (m *Model) ExecuteConsoleCmd(cmd ConsoleCmd) string {
	switch cmd.Command {
	case "help":
		return helpText()
	case "status":
		return m.cmdStatus()
	case "clear":
		m.ConsoleOutput = nil
		return ""
	case "regen", "regenerate":
		return m.cmdRegen()
	case "add":
		return m.cmdAdd(cmd.Args)
	case "remove":
		return m.cmdRemove(cmd.Args)
	case "heal":
		return m.cmdHeal(cmd.Args)
	case "damage":
		return m.cmdDamage(cmd.Args)
	case "set":
		return m.cmdSet(cmd.Args)
	case "effect":
		return m.cmdEffect(cmd.Args)
	default:
		return "Unknown command. Type 'help' for commands."
	}
}

func (m *Model) cmdStatus() string {
	var b strings.Builder
	fmt.Fprintf(&b, "Seed: %d  |  Turn: player %d", m.Seed, m.CurrentPlayer+1)
	for i, pl := range m.Players {
		fmt.Fprintf(&b, "\n  Player %d: (%d,%d) HP=%d Ult=%d", i+1, pl.X, pl.Y, pl.HP, pl.UltCharges)
		for _, e := range pl.Effects {
			fmt.Fprintf(&b, " [%s:%d]", string(e.Type), e.Duration)
		}
	}
	for i, en := range m.Enemys {
		fmt.Fprintf(&b, "\n  Enemy %d: (%d,%d) HP=%d", i+1, en.X, en.Y, en.HP)
		for _, e := range en.Effects {
			fmt.Fprintf(&b, " [%s:%d]", string(e.Type), e.Duration)
		}
	}
	return b.String()
}

func (m *Model) requireGame() string {
	if m.Screen != ScreenGame {
		return "Commands are only available during a game."
	}
	return ""
}

func (m *Model) cmdRegen() string {
	if err := m.requireGame(); err != "" {
		return err
	}
	blocked := make(map[Point]bool)
	for _, pl := range m.Players {
		blocked[Point{pl.X, pl.Y}] = true
	}
	for _, en := range m.Enemys {
		blocked[Point{en.X, en.Y}] = true
	}
	rng := rand.New(rand.NewSource(m.Seed))
	wallPositions := GenerateTiles(GridW/2, GridH/2, wallCount, blocked, rng)
	walls := make(map[Point]Wall, len(wallPositions))
	for p := range wallPositions {
		walls[p] = Wall{HP: WallHP}
		blocked[p] = true
	}
	m.Walls = walls
	m.Water = GenerateTiles(GridW/2, GridH/2, waterCount, blocked, rng)
	m.FireTiles = make(map[Point]int)
	m.SmokeTiles = make(map[Point]int)
	return "Map regenerated (walls, water). Fire and smoke cleared."
}

func (m *Model) cmdAdd(args []string) string {
	if err := m.requireGame(); err != "" {
		return err
	}
	if len(args) < 1 {
		return "Usage: add <player|enemy>"
	}
	tile := m.findFreeTile()
	if tile == nil {
		return "No free tile available."
	}
	switch args[0] {
	case "player":
		idx := len(m.Players)
		m.Players = append(m.Players, Player{
			X: tile.X, Y: tile.Y,
			HP: MaxHP, UltCharges: maxUltCharges,
			Style: m.Styles.PlayerStyles[idx%len(m.Styles.PlayerStyles)],
		})
		return fmt.Sprintf("Player %d added at (%d,%d).", idx+1, tile.X, tile.Y)
	case "enemy":
		idx := len(m.Enemys)
		m.Enemys = append(m.Enemys, Enemy{
			X: tile.X, Y: tile.Y,
			HP:    MaxHP,
			Style: m.Styles.EnemysStyles[idx%len(m.Styles.EnemysStyles)],
		})
		return fmt.Sprintf("Enemy %d added at (%d,%d).", idx+1, tile.X, tile.Y)
	default:
		return "Usage: add <player|enemy>"
	}
}

func (m *Model) cmdRemove(args []string) string {
	if err := m.requireGame(); err != "" {
		return err
	}
	if len(args) < 1 {
		return "Usage: remove <player|enemy|effect> ..."
	}
	switch strings.ToLower(args[0]) {
	case "player":
		if len(args) < 2 {
			return "Usage: remove player <N>"
		}
		n, err := strconv.Atoi(args[1])
		if err != nil || n < 1 {
			return "Invalid number."
		}
		if n-1 >= len(m.Players) {
			return fmt.Sprintf("Player %d not found.", n)
		}
		m.Players = append(m.Players[:n-1], m.Players[n:]...)
		if m.CurrentPlayer >= len(m.Players) {
			m.CurrentPlayer = 0
		}
		if len(m.Players) == 0 {
			m.endGame(3)
			return "Last player removed — game over."
		}
		return fmt.Sprintf("Player %d removed.", n)
	case "enemy":
		if len(args) < 2 {
			return "Usage: remove enemy <N>"
		}
		n, err := strconv.Atoi(args[1])
		if err != nil || n < 1 {
			return "Invalid number."
		}
		if n-1 >= len(m.Enemys) {
			return fmt.Sprintf("Enemy %d not found.", n)
		}
		m.Enemys = append(m.Enemys[:n-1], m.Enemys[n:]...)
		if len(m.Enemys) == 0 {
			m.endGame(15)
			return "Last enemy removed — you win!"
		}
		return fmt.Sprintf("Enemy %d removed.", n)
	case "effect":
		return m.cmdRemoveEffect(args[1:])
	default:
		return "Usage: remove <player|enemy|effect> ..."
	}
}

func (m *Model) cmdRemoveEffect(args []string) string {
	if len(args) < 2 {
		return "Usage: remove effect <target> <type>"
	}
	isPlayer, idx, ok := ParseTarget(args[0])
	if !ok {
		return "Invalid target."
	}
	var effectType EffectType
	switch strings.ToLower(args[1]) {
	case "fire":
		effectType = EffectFire
	case "wet":
		effectType = EffectWet
	case "smoke":
		effectType = EffectSmoke
	default:
		return "Invalid effect type. Use: fire, wet, smoke"
	}
	if isPlayer {
		if idx >= len(m.Players) {
			return fmt.Sprintf("Player %d not found.", idx+1)
		}
		m.Players[idx].Effects = RemoveEffect(m.Players[idx].Effects, effectType)
		return fmt.Sprintf("Removed %s from Player %d.", args[1], idx+1)
	}
	if idx >= len(m.Enemys) {
		return fmt.Sprintf("Enemy %d not found.", idx+1)
	}
	m.Enemys[idx].Effects = RemoveEffect(m.Enemys[idx].Effects, effectType)
	return fmt.Sprintf("Removed %s from Enemy %d.", args[1], idx+1)
}

func (m *Model) cmdHeal(args []string) string {
	if err := m.requireGame(); err != "" {
		return err
	}
	if len(args) < 2 {
		return "Usage: heal <target> <N>"
	}
	isPlayer, idx, ok := ParseTarget(args[0])
	if !ok {
		return "Invalid target. Use p1, p2, ... or e1, e2, ..."
	}
	n, err := strconv.Atoi(args[1])
	if err != nil || n < 1 {
		return "Invalid number."
	}
	if isPlayer {
		if idx >= len(m.Players) {
			return fmt.Sprintf("Player %d not found.", idx+1)
		}
		m.Players[idx].HP += n
		if m.Players[idx].HP > MaxHP {
			m.Players[idx].HP = MaxHP
		}
		return fmt.Sprintf("Player %d healed to %d HP.", idx+1, m.Players[idx].HP)
	}
	if idx >= len(m.Enemys) {
		return fmt.Sprintf("Enemy %d not found.", idx+1)
	}
	m.Enemys[idx].HP += n
	if m.Enemys[idx].HP > MaxHP {
		m.Enemys[idx].HP = MaxHP
	}
	return fmt.Sprintf("Enemy %d healed to %d HP.", idx+1, m.Enemys[idx].HP)
}

func (m *Model) cmdDamage(args []string) string {
	if err := m.requireGame(); err != "" {
		return err
	}
	if len(args) < 2 {
		return "Usage: damage <target> <N>"
	}
	isPlayer, idx, ok := ParseTarget(args[0])
	if !ok {
		return "Invalid target."
	}
	n, err := strconv.Atoi(args[1])
	if err != nil || n < 1 {
		return "Invalid number."
	}
	if isPlayer {
		if idx >= len(m.Players) {
			return fmt.Sprintf("Player %d not found.", idx+1)
		}
		m.Players[idx].HP -= n
		if m.Players[idx].HP <= 0 {
			m.Players[idx].HP = 0
			m.Players = append(m.Players[:idx], m.Players[idx+1:]...)
			if m.CurrentPlayer >= len(m.Players) {
				m.CurrentPlayer = 0
			}
			if len(m.Players) == 0 {
				m.endGame(3)
				return "Player died — game over."
			}
			return fmt.Sprintf("Player %d died.", idx+1)
		}
		return fmt.Sprintf("Player %d now has %d HP.", idx+1, m.Players[idx].HP)
	}
	if idx >= len(m.Enemys) {
		return fmt.Sprintf("Enemy %d not found.", idx+1)
	}
	m.Enemys[idx].HP -= n
	if m.Enemys[idx].HP <= 0 {
		m.Enemys[idx].HP = 0
		m.Enemys = append(m.Enemys[:idx], m.Enemys[idx+1:]...)
		if len(m.Enemys) == 0 {
			m.endGame(15)
			return "Enemy died — you win!"
		}
		return fmt.Sprintf("Enemy %d died.", idx+1)
	}
	return fmt.Sprintf("Enemy %d now has %d HP.", idx+1, m.Enemys[idx].HP)
}

func (m *Model) cmdSet(args []string) string {
	if err := m.requireGame(); err != "" {
		return err
	}
	if len(args) < 3 {
		return "Usage: set hp <target> <N>"
	}
	if strings.ToLower(args[0]) != "hp" {
		return "Usage: set hp <target> <N>"
	}
	isPlayer, idx, ok := ParseTarget(args[1])
	if !ok {
		return "Invalid target."
	}
	n, err := strconv.Atoi(args[2])
	if err != nil || n < 0 {
		return "Invalid number."
	}
	if n > MaxHP {
		n = MaxHP
	}
	if isPlayer {
		if idx >= len(m.Players) {
			return fmt.Sprintf("Player %d not found.", idx+1)
		}
		if n <= 0 {
			m.Players = append(m.Players[:idx], m.Players[idx+1:]...)
			if m.CurrentPlayer >= len(m.Players) {
				m.CurrentPlayer = 0
			}
			if len(m.Players) == 0 {
				m.endGame(3)
				return "Player died — game over."
			}
			return fmt.Sprintf("Player %d died.", idx+1)
		}
		m.Players[idx].HP = n
		return fmt.Sprintf("Player %d HP set to %d.", idx+1, n)
	}
	if idx >= len(m.Enemys) {
		return fmt.Sprintf("Enemy %d not found.", idx+1)
	}
	if n <= 0 {
		m.Enemys = append(m.Enemys[:idx], m.Enemys[idx+1:]...)
		if len(m.Enemys) == 0 {
			m.endGame(15)
			return "Enemy died — you win!"
		}
		return fmt.Sprintf("Enemy %d died.", idx+1)
	}
	m.Enemys[idx].HP = n
	return fmt.Sprintf("Enemy %d HP set to %d.", idx+1, n)
}

func (m *Model) cmdEffect(args []string) string {
	if err := m.requireGame(); err != "" {
		return err
	}
	if len(args) < 2 {
		return "Usage: effect <target> <type> [dur]"
	}
	isPlayer, idx, ok := ParseTarget(args[0])
	if !ok {
		return "Invalid target."
	}
	var effectType EffectType
	switch strings.ToLower(args[1]) {
	case "fire":
		effectType = EffectFire
	case "wet":
		effectType = EffectWet
	case "smoke":
		effectType = EffectSmoke
	default:
		return "Invalid effect type. Use: fire, wet, smoke"
	}
	dur := 3
	if len(args) >= 3 {
		d, err := strconv.Atoi(args[2])
		if err == nil && d > 0 {
			dur = d
		}
	}
	eff := Effect{Type: effectType, Duration: dur}
	if isPlayer {
		if idx >= len(m.Players) {
			return fmt.Sprintf("Player %d not found.", idx+1)
		}
		m.Players[idx].Effects = ResolveEffects(m.Players[idx].Effects, eff)
		return fmt.Sprintf("Added %s (dur:%d) to Player %d.", args[1], dur, idx+1)
	}
	if idx >= len(m.Enemys) {
		return fmt.Sprintf("Enemy %d not found.", idx+1)
	}
	m.Enemys[idx].Effects = ResolveEffects(m.Enemys[idx].Effects, eff)
	return fmt.Sprintf("Added %s (dur:%d) to Enemy %d.", args[1], dur, idx+1)
}

func (m *Model) findFreeTile() *Point {
	for x := 0; x < GridW; x++ {
		for y := 0; y < GridH; y++ {
			p := Point{x, y}
			if m.IsWall(p) {
				continue
			}
			occupied := false
			for _, pl := range m.Players {
				if pl.X == x && pl.Y == y {
					occupied = true
					break
				}
			}
			if occupied {
				continue
			}
			for _, en := range m.Enemys {
				if en.X == x && en.Y == y {
					occupied = true
					break
				}
			}
			if occupied {
				continue
			}
			return &Point{x, y}
		}
	}
	return nil
}
