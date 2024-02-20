package generation_processor

import (
	"github.com/passeriform/conway-gox/internal/cell"
)

type Coordinates struct {
	Row    int
	Column int
}

// TODO: Add species, alleles, dominant/recessive genes and evolution index
type GenerationProcessor struct {
	cellHealth map[Coordinates]int
}

func FromCells(c []*cell.Cell) GenerationProcessor {
	gp := GenerationProcessor{make(map[Coordinates]int)}
	for _, cell := range c {
		row, column := cell.GetPosition()
		gp.cellHealth[Coordinates{row, column}] = -1
	}
	return gp
}

func (gp *GenerationProcessor) Expand(influence int) {
	cellHealth := make(map[Coordinates]int)
	for cellKey := range gp.cellHealth {
		for r := cellKey.Row - influence; r <= cellKey.Row+influence; r++ {
			for c := cellKey.Column - influence; c <= cellKey.Column+influence; c++ {
				if r == cellKey.Row && c == cellKey.Column {
					continue
				}
				_, originalFound := gp.cellHealth[Coordinates{r, c}]
				if originalFound {
					continue
				}
				health, found := cellHealth[Coordinates{r, c}]
				if found {
					cellHealth[Coordinates{r, c}] = health + 1
				} else {
					cellHealth[Coordinates{r, c}] = 1
				}
			}
		}
	}
	gp.cellHealth = cellHealth
}

func (gp *GenerationProcessor) Reduce(minimum int, maximum int) {
	cellHealth := make(map[Coordinates]int)
	for key, health := range gp.cellHealth {
		if health >= minimum && health <= maximum {
			cellHealth[key] = health
		}
	}
	gp.cellHealth = cellHealth
}

func (gp GenerationProcessor) ToCells() []*cell.Cell {
	cells := []*cell.Cell{}
	for key := range gp.cellHealth {
		newCell := cell.Create(key.Row, key.Column)
		cells = append(cells, &newCell)
	}
	return cells
}
