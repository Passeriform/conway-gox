package cell

import (
	"testing"
)

func PopulateCellsData(r int, c int) []Cell {
	cells := []Cell{}
	for ridx := r - 1; ridx <= r+1; ridx++ {
		for cidx := c - 1; cidx <= c+1; cidx++ {
			cell := Cell{
				row:       ridx,
				column:    cidx,
				Neighbors: []*Cell{},
			}
			cells = append(cells, cell)
		}
	}

	cells[0].Neighbors = []*Cell{
		/*_,*/ &cells[1],
		&cells[3], &cells[4],
	}
	cells[1].Neighbors = []*Cell{
		&cells[0] /*_,*/, &cells[2],
		&cells[3], &cells[4], &cells[5],
	}
	cells[2].Neighbors = []*Cell{
		&cells[1], /*_,*/
		&cells[4], &cells[5],
	}
	cells[3].Neighbors = []*Cell{
		&cells[0], &cells[1],
		/*_,*/ &cells[4],
		&cells[6], &cells[7],
	}
	cells[4].Neighbors = []*Cell{
		&cells[0], &cells[1], &cells[2],
		&cells[3] /*_,*/, &cells[5],
		&cells[6], &cells[7], &cells[8],
	}
	cells[5].Neighbors = []*Cell{
		&cells[1], &cells[2],
		&cells[4], /*_,*/
		&cells[7], &cells[8],
	}
	cells[6].Neighbors = []*Cell{
		&cells[3], &cells[4],
		/*_,*/ &cells[7],
	}
	cells[7].Neighbors = []*Cell{
		&cells[3], &cells[4], &cells[5],
		&cells[6] /*_,*/, &cells[8],
	}
	cells[8].Neighbors = []*Cell{
		&cells[4], &cells[5],
		&cells[7], /*_,*/
	}
	return cells
}

func TestWillSurvive(t *testing.T) {
	cells := PopulateCellsData(3, 2)

	testCases := []struct {
		name      string
		Neighbors []*Cell
		want      bool
	}{
		{
			"Should not survive overpopulation",
			[]*Cell{
				&cells[0], &cells[1], /*_,*/
				&cells[3], /*test,*/ /*_,*/
				/*_,*/ &cells[7], /*_,*/
			},
			false,
		},
		{
			"Should not survive underpopulation",
			[]*Cell{
				/*_,*/ /*_,*/ /*_,*/
				/*_,*/ /*test,*/ &cells[5],
				/*_,*/ /*_,*/ /*_,*/
			},
			false,
		},
		{
			"Should survive normal scenario",
			[]*Cell{
				/*_,*/ /*_,*/ &cells[2],
				/*_,*/ /*test,*/ /*_,*/
				/*_,*/ &cells[7], /*_,*/
			},
			true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCell := cells[4]
			testCell.Neighbors = testCase.Neighbors
			result := testCell.WillSurvive()
			if result != testCase.want {
				t.Errorf("got %v, want %v", result, testCase.want)
			}
		})
	}
}

func TestKill(t *testing.T) {
	cells := PopulateCellsData(5, 4)
	neighbors := cells[4].Neighbors
	cells[4].Kill()
	for _, neighbor := range neighbors {
		if neighbor == &cells[4] {
			t.Errorf("got (%v, %v) with reference to killed cell, want no reference to killed cell", neighbor.row, neighbor.column)
		}
	}
}
