package patterns

import (
	"github.com/passeriform/conway-gox/internal/cell"
)

// TODO: Convert to separate saveState files and use via load

func GetPrimitive(t string, r int, c int) []cell.Cell {
	if t == "Block" {
		return []cell.Cell{
			cell.New(r, c),
			cell.New(r, c+1),
			cell.New(r+1, c),
			cell.New(r+1, c+1),
		}
	} else if t == "Beehive" {
		return []cell.Cell{
			cell.New(r-1, c-1),
			cell.New(r-1, c),
			cell.New(r, c-2),
			cell.New(r, c+1),
			cell.New(r+1, c-1),
			cell.New(r+1, c),
		}
	} else if t == "Loaf" {
		return []cell.Cell{
			cell.New(r-1, c-1),
			cell.New(r-1, c),
			cell.New(r, c-2),
			cell.New(r, c+1),
			cell.New(r+1, c-1),
			cell.New(r+1, c+1),
			cell.New(r+2, c),
		}
	} else if t == "Boat" {
		return []cell.Cell{
			cell.New(r-1, c-1),
			cell.New(r-1, c),
			cell.New(r, c-1),
			cell.New(r, c+1),
			cell.New(r+1, c),
		}
	} else if t == "Tub" {
		return []cell.Cell{
			cell.New(r-1, c),
			cell.New(r, c-1),
			cell.New(r, c+1),
			cell.New(r+1, c),
		}
	} else if t == "AircraftCarrier" {
		return []cell.Cell{
			cell.New(r-1, c-1),
			cell.New(r-1, c),
			cell.New(r, c-1),
			cell.New(r, c+2),
			cell.New(r+1, c+1),
			cell.New(r+1, c+2),
		}
	} else if t == "Blinker" {
		return []cell.Cell{
			cell.New(r-1, c),
			cell.New(r-1, c),
			cell.New(r+1, c),
		}
	} else if t == "Toad" {
		return []cell.Cell{
			cell.New(r, c-1),
			cell.New(r, c),
			cell.New(r, c+1),
			cell.New(r+1, c-2),
			cell.New(r+1, c-1),
			cell.New(r+1, c),
		}
	} else if t == "Beacon" {
		return []cell.Cell{
			cell.New(r-1, c-1),
			cell.New(r-1, c),
			cell.New(r, c-1),
			cell.New(r, c),
			cell.New(r+1, c+1),
			cell.New(r+1, c+2),
			cell.New(r+2, c+1),
			cell.New(r+2, c+2),
		}
	} else if t == "Pulsar" {
		return []cell.Cell{
			cell.New(r-6, c-4),
			cell.New(r-6, c-3),
			cell.New(r-6, c-2),
			cell.New(r-6, c+2),
			cell.New(r-6, c-3),
			cell.New(r-6, c+4),
			cell.New(r-4, c-6),
			cell.New(r-4, c+1),
			cell.New(r-4, c+1),
			cell.New(r-4, c+6),
			cell.New(r-3, c-6),
			cell.New(r-3, c+1),
			cell.New(r-3, c+1),
			cell.New(r-3, c+6),
			cell.New(r-2, c-6),
			cell.New(r-2, c+1),
			cell.New(r-2, c+1),
			cell.New(r-2, c+6),
			cell.New(r-1, c-4),
			cell.New(r-1, c-3),
			cell.New(r-1, c-2),
			cell.New(r-1, c+2),
			cell.New(r-1, c-3),
			cell.New(r-1, c+4),
			cell.New(r+1, c-4),
			cell.New(r+1, c-3),
			cell.New(r+1, c-2),
			cell.New(r+1, c+2),
			cell.New(r+1, c-3),
			cell.New(r+1, c+4),
			cell.New(r+2, c-6),
			cell.New(r+2, c+1),
			cell.New(r+2, c+1),
			cell.New(r+2, c+6),
			cell.New(r+3, c-6),
			cell.New(r+3, c+1),
			cell.New(r+3, c+1),
			cell.New(r+3, c+6),
			cell.New(r+4, c-6),
			cell.New(r+4, c+1),
			cell.New(r+4, c+1),
			cell.New(r+4, c+6),
			cell.New(r+6, c-4),
			cell.New(r+6, c-3),
			cell.New(r+6, c-2),
			cell.New(r+6, c+2),
			cell.New(r+6, c-3),
			cell.New(r+6, c+4),
		}
	} else if t == "PentaDecathlon" {
		return []cell.Cell{
			cell.New(r-6, c-1),
			cell.New(r-6, c),
			cell.New(r-6, c+1),
			cell.New(r-5, c),
			cell.New(r-4, c),
			cell.New(r-3, c-1),
			cell.New(r-3, c),
			cell.New(r-3, c+1),
			cell.New(r-1, c-1),
			cell.New(r-1, c),
			cell.New(r-1, c+1),
			cell.New(r, c-1),
			cell.New(r, c),
			cell.New(r, c+1),
			cell.New(r+2, c-1),
			cell.New(r+2, c),
			cell.New(r+2, c+1),
			cell.New(r+3, c),
			cell.New(r+4, c),
			cell.New(r+5, c-1),
			cell.New(r+5, c),
			cell.New(r+5, c+1),
		}
	}

	return []cell.Cell{}
}
