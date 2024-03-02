package cell_map

import (
	"github.com/passeriform/conway-gox/internal/cell"
	"github.com/passeriform/conway-gox/internal/generation_processor"
	"github.com/passeriform/conway-gox/internal/utility"
)

type Map struct {
	cells []*cell.Cell
}

type Bounds struct {
	Top    int
	Right  int
	Bottom int
	Left   int
}

func (m *Map) recomputeNeighbors() {
	for current := range m.cells {
		m.cells[current].Neighbors = []*cell.Cell{}
		for next := range m.cells {
			if next == current {
				continue
			}

			if m.cells[current].IsNeighbor(m.cells[next]) {
				m.cells[current].Neighbors = append(m.cells[current].Neighbors, m.cells[next])
			}
		}
	}
}

func FromCells(c []cell.Cell) Map {
	cells := make([]*cell.Cell, len(c))
	for cellIndex := range c {
		cells[cellIndex] = &c[cellIndex]
	}

	cm := Map{cells}
	cm.recomputeNeighbors()
	return cm
}

func (m Map) GetCells() []*cell.Cell {
	return m.cells
}

func (m Map) GetBounds() Bounds {
	top, right, bottom, left := 0, 0, 0, 0

	for _, element := range m.cells {
		row, column := element.GetPosition()
		top = min(row, top)
		right = max(column, right)
		bottom = max(row, bottom)
		left = min(column, left)
	}

	return Bounds{
		Top:    top,
		Right:  right,
		Bottom: bottom,
		Left:   left,
	}
}

func (m *Map) Next() {
	// TODO: Move to generation_processor
	nextCells, markedForKill := utility.Partition(m.cells, func(element *cell.Cell) bool {
		return element.WillSurvive()
	})

	processor := generation_processor.FromCells(m.cells)
	processor.Expand(1)
	processor.Reduce(3, 3)

	nextCells = append(nextCells, processor.ToCells()...)

	for _, cell := range markedForKill {
		cell.Kill()
	}

	m.cells = nextCells
	m.recomputeNeighbors()
}

func (m Map) Rasterize(padding int) [][]bool {
	raster := make([][]bool, 2*padding)
	for idx := range raster {
		raster[idx] = make([]bool, 2*padding)
	}

	for _, cell := range m.cells {
		row, column := cell.GetPosition()
		raster[padding+row][padding+column] = true
	}

	return raster
}

func (m Map) EncodeJson(padding int) [][2]int {
	jsonData := [][2]int{}

	for _, cell := range m.GetCells() {
		row, column := cell.GetPosition()
		jsonData = append(jsonData, [2]int{padding + row, padding + column})
	}

	return jsonData
}

func DecodeJson(jb [][2]int, padding int) Map {
	cells := make([]cell.Cell, len(jb))

	for idx, cb := range jb {
		cells[idx] = cell.New(padding+cb[0], padding+cb[1])
	}

	return FromCells(cells)
}
