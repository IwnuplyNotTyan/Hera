package generate

import "hera/utils"

// currentRange returns the movement or shoot range for the current player.
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

// IsInRange checks whether (col, row) is within the current player's range and not wall-blocked.
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

// Reachable returns the set of points reachable from (sx, sy) within r steps, avoiding walls and occupied cells.
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
			if m.IsWall(np) {
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

// HasWallBetweenPoints returns true if a wall lies on the Bresenham line between the two points.
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
			if m.IsWall(Point{x0, y0}) {
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
		pts = append(pts, p)
	}
	return pts
}

func (m *Model) ultInAxisRange(cx, cy int) bool {
	if len(m.Players) == 0 || m.CurrentPlayer >= len(m.Players) {
		return false
	}
	current := m.Players[m.CurrentPlayer]
	return cx == current.X || cy == current.Y
}

// IsWall reports whether a wall exists at the given point.
func (m *Model) IsWall(p Point) bool {
	_, ok := m.Walls[p]
	return ok
}

// OccupiedByOther returns true if a non-current player or enemy occupies (x, y).
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
