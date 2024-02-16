package utility

import (
	"testing"
)

func TestFilter(t *testing.T) {
	array := []int{1, 4, 2, 4, 9, 2, 7, 3, 9, 2, 5}

	testCases := []struct {
		name string
		what func(int) bool
		want []int
	}{
		{
			"Should filter out single element",
			func(x int) bool { return x != 7 },
			[]int{1, 4, 2, 4, 9, 2, 3, 9, 2, 5},
		},
		{
			"Should filter out multiple element",
			func(x int) bool { return x != 4 },
			[]int{1, 2, 9, 2, 7, 3, 9, 2, 5},
		},
		{
			"Should filter out no element if not found",
			func(x int) bool { return x != 11 },
			[]int{1, 4, 2, 4, 9, 2, 7, 3, 9, 2, 5},
		},
		{
			"Should keep even elements only",
			func(x int) bool { return x%2 == 0 },
			[]int{4, 2, 4, 2, 2},
		},
		{
			"Should remove elements less than 5",
			func(x int) bool { return x >= 5 },
			[]int{9, 7, 9, 5},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testArray := array
			result := Filter[int](testArray, testCase.what)
			for index, resultElement := range result {
				if resultElement != testCase.want[index] {
					t.Errorf("got %v, want %v", result, testCase.want)
				}
			}
		})
	}
}
