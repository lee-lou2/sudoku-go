package cmd

import (
	"testing"
)

// compareSudoku 결과를 비교하는 함수
func compareSudoku(s1, s2 [][]int) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := range s1 {
		if len(s1[i]) != len(s2[i]) {
			return false
		}
		for j := range s1[i] {
			if s1[i][j] != s2[i][j] {
				return false
			}
		}
	}
	return true
}

func TestSolveSudoku(t *testing.T) {
	tests := []struct {
		name   string
		input  [][]int
		output [][]int
	}{
		{
			name: "Test Case 1 - Valid 9x9 Sudoku",
			input: [][]int{
				{0, 0, 0, 2, 6, 0, 7, 0, 1},
				{6, 8, 0, 0, 7, 0, 0, 9, 0},
				{1, 9, 0, 0, 0, 4, 5, 0, 0},
				{8, 2, 0, 1, 0, 0, 0, 4, 0},
				{0, 0, 4, 6, 0, 2, 9, 0, 0},
				{0, 5, 0, 0, 0, 3, 0, 2, 8},
				{0, 0, 9, 3, 0, 0, 0, 7, 4},
				{0, 4, 0, 0, 5, 0, 0, 3, 6},
				{7, 0, 3, 0, 1, 8, 0, 0, 0},
			},
			output: [][]int{
				{4, 3, 5, 2, 6, 9, 7, 8, 1},
				{6, 8, 2, 5, 7, 1, 4, 9, 3},
				{1, 9, 7, 8, 3, 4, 5, 6, 2},
				{8, 2, 6, 1, 9, 5, 3, 4, 7},
				{3, 7, 4, 6, 8, 2, 9, 1, 5},
				{9, 5, 1, 7, 4, 3, 6, 2, 8},
				{5, 1, 9, 3, 2, 6, 8, 7, 4},
				{2, 4, 8, 9, 5, 7, 1, 3, 6},
				{7, 6, 3, 4, 1, 8, 2, 5, 9},
			},
		},
		{
			name: "Test Case 2 - Incomplete 9x9 Sudoku",
			input: [][]int{
				{1, 0, 0, 4, 8, 9, 0, 0, 6},
				{7, 3, 0, 0, 0, 0, 0, 4, 0},
				{0, 0, 0, 0, 0, 1, 2, 9, 5},
				{0, 0, 7, 1, 2, 0, 6, 0, 0},
				{5, 0, 0, 7, 0, 3, 0, 0, 8},
				{0, 0, 6, 0, 9, 5, 7, 0, 0},
				{9, 1, 4, 6, 0, 0, 0, 0, 0},
				{0, 2, 0, 0, 0, 0, 0, 3, 7},
				{8, 0, 0, 5, 1, 2, 0, 0, 4},
			},
			output: [][]int{
				{1, 5, 2, 4, 8, 9, 3, 7, 6},
				{7, 3, 9, 2, 5, 6, 8, 4, 1},
				{4, 6, 8, 3, 7, 1, 2, 9, 5},
				{3, 8, 7, 1, 2, 4, 6, 5, 9},
				{5, 9, 1, 7, 6, 3, 4, 2, 8},
				{2, 4, 6, 8, 9, 5, 7, 1, 3},
				{9, 1, 4, 6, 3, 7, 5, 8, 2},
				{6, 2, 5, 9, 4, 8, 1, 3, 7},
				{8, 7, 3, 5, 1, 2, 9, 6, 4},
			},
		},
		{
			name: "Test Case 3 - Valid 6x6 Sudoku",
			input: [][]int{
				{0, 0, 0, 6, 2, 0},
				{2, 4, 0, 0, 0, 0},
				{4, 3, 0, 0, 0, 0},
				{0, 0, 0, 0, 3, 2},
				{0, 0, 0, 0, 6, 3},
				{0, 5, 3, 0, 0, 0},
			},
			output: [][]int{
				{3, 1, 5, 6, 2, 4},
				{2, 4, 6, 3, 1, 5},
				{4, 3, 2, 1, 5, 6},
				{5, 6, 1, 4, 3, 2},
				{1, 2, 4, 5, 6, 3},
				{6, 5, 3, 2, 4, 1},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, _ := SolveSudoku(tc.input)
			if result == nil {
				t.Errorf("failed to solve sudoku")
			}
			if !compareSudoku(result, tc.output) {
				t.Errorf("expected %v, got %v", tc.output, result)
			}
		})
	}
}
