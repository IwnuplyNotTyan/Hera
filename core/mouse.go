package generate

var layoutElements []Element
var gridOffsetX, gridOffsetY int

func resetLayout() {
	layoutElements = nil
}

func trackElement(elem Element) {
	layoutElements = append(layoutElements, elem)
}

func hitTest(screenX, screenY int) *Element {
	absX := screenX - gridOffsetX
	absY := screenY - gridOffsetY
	for i := range layoutElements {
		elem := &layoutElements[i]
		if absX >= elem.X && absX < elem.X+elem.Width &&
			absY >= elem.Y && absY < elem.Y+elem.Height {
			return elem
		}
	}
	return nil
}

func cellWidth() int {
	return 3
}

func cellHeight() int {
	return 1
}

func handleGridClick(col, row int) (int, int) {
	return col, row
}
