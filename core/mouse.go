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

func cellWidth() int  { return 3 }
func cellHeight() int { return 1 }
