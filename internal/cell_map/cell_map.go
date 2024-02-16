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
		for next, nextCell := range m.cells {
			if next == current {
				continue
			}

			if m.cells[current].IsNeighbor(nextCell) {
				m.cells[current].Neighbors = append(m.cells[current].Neighbors, nextCell)
			}
		}
	}
}

func Create() Map {
	return Map{cells: []*cell.Cell{}}
}

func (m *Map) AddCells(c []cell.Cell) {
	for cellIndex := range c {
		m.cells = append(m.cells, &c[cellIndex])
	}

	m.recomputeNeighbors()
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

func (m *Map) Step() {
	markedForKill := []*cell.Cell{}

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

// TODO: Implement loadState and saveState
