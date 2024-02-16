package cell

import (
	"math"

	"github.com/passeriform/conway-gox/internal/utility"
)

type Cell struct {
	row       int
	column    int
	Neighbors []*Cell
}

func Create(r int, c int) Cell {
	return Cell{r, c, []*Cell{}}
}

func (c *Cell) IsNeighbor(nc *Cell) bool {
	return math.Abs(float64(c.row-nc.row)) <= 1 &&
		math.Abs(float64(c.column-nc.column)) <= 1
}

func (c *Cell) GetPosition() (int, int) {
	return c.row, c.column
}

func (c *Cell) WillSurvive() bool {
	return len(c.Neighbors) >= 2 && len(c.Neighbors) <= 3
}

func (c *Cell) Kill() {
	for _, neighbor := range c.Neighbors {
		neighbor.Neighbors = utility.Filter[*Cell](neighbor.Neighbors, func(element *Cell) bool {
			return element != c
		})
	}
}
