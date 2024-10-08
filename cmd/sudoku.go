package cmd

import (
	"context"
	"sort"
	"sync"
)

// SubGrid 스도쿠 박스 크기 구조체
type SubGrid struct {
	width, height int
}

// Cell 스도쿠 좌표 구조체
type Cell struct {
	x, y int
}

// GridMap 스도쿠 박스 크기 맵
var GridMap = map[int]SubGrid{
	4:  SubGrid{2, 2},
	6:  SubGrid{3, 2},
	8:  SubGrid{4, 2},
	9:  SubGrid{3, 3},
	10: SubGrid{5, 2},
	12: SubGrid{4, 3},
	15: SubGrid{5, 3},
	16: SubGrid{4, 4},
	25: SubGrid{5, 5},
	36: SubGrid{6, 6},
	49: SubGrid{7, 7},
	64: SubGrid{8, 8},
	81: SubGrid{9, 9},
}

// contains 배열에 target 이 있는지 확인
func contains(arr []int, target int) bool {
	for _, v := range arr {
		if v == target {
			return true
		}
	}
	return false
}

// hasDuplicates 중복된 값이 있는지 확인하는 함수
func hasDuplicates(arr []int) bool {
	seen := make(map[int]bool)
	for _, v := range arr {
		if v != 0 {
			if seen[v] {
				return true
			}
			seen[v] = true
		}
	}
	return false
}

// getEmptyCell 0인 좌표와 0의 개수를 반환
func getEmptyCell(grid [][]int) *Cell {
	subGrid := GridMap[len(grid)]
	type cellInfo struct {
		cell Cell
		cnt  int
	}
	cellInfoList := make([]cellInfo, 0)

	for y, row := range grid {
		for x, cell := range row {
			if cell == 0 {
				// 가로 행, 세로 열에 0이 몇 개인지 확인
				emptyCnt := 0
				// 가로 행의 0 개수 확인
				for i := 0; i < subGrid.width*subGrid.height; i++ {
					if grid[y][i] == 0 {
						emptyCnt++
					}
				}

				// 세로 열의 0 개수 확인
				for i := 0; i < subGrid.width*subGrid.height; i++ {
					if grid[i][x] == 0 {
						emptyCnt++
					}
				}

				// 박스에 0이 몇 개인지 확인
				subGridX := (x / subGrid.width) * subGrid.width
				subGridY := (y / subGrid.height) * subGrid.height
				for i := subGridY; i < subGridY+subGrid.height; i++ {
					for j := subGridX; j < subGridX+subGrid.width; j++ {
						if grid[i][j] == 0 {
							emptyCnt++
						}
					}
				}

				// 0의 개수와 좌표를 저장
				cellInfoList = append(cellInfoList, cellInfo{cell: Cell{x, y}, cnt: emptyCnt})
			}
		}
	}

	if len(cellInfoList) == 0 {
		return nil
	}

	// emptyCnt 기준으로 정렬
	sort.Slice(cellInfoList, func(i, j int) bool {
		return cellInfoList[i].cnt < cellInfoList[j].cnt
	})

	return &cellInfoList[0].cell
}

// getCandidateNumbers 해당 좌표에 들어갈 수 있는 숫자를 반환
func getCandidateNumbers(grid [][]int, cell *Cell) []int {
	subGrid := GridMap[len(grid)]

	// 해당 좌표의 행 데이터 가져오기
	rowData := grid[cell.y]
	row := make([]int, subGrid.width*subGrid.height)
	copy(row, rowData)

	// 해당 좌표의 열 데이터 가져오기
	col := make([]int, subGrid.width*subGrid.height)
	for i := 0; i < subGrid.width*subGrid.height; i++ {
		col[i] = grid[i][cell.x]
	}

	// 해당 좌표의 서브 그리드 데이터 가져오기
	subGridX := (cell.x / subGrid.width) * subGrid.width
	subGridY := (cell.y / subGrid.height) * subGrid.height
	currSubGrid := make([]int, 0, subGrid.width*subGrid.height)
	for i := subGridY; i < subGridY+subGrid.height; i++ {
		for j := subGridX; j < subGridX+subGrid.width; j++ {
			currSubGrid = append(currSubGrid, grid[i][j])
		}
	}

	// row, col, box 에 없는 숫자 찾기
	numbers := make([]int, 0, subGrid.width*subGrid.height)
	for i := 1; i <= subGrid.width*subGrid.height; i++ {
		if !contains(row, i) && !contains(col, i) && !contains(currSubGrid, i) {
			numbers = append(numbers, i)
		}
	}

	return numbers
}

// tryPlaceNumber 스도쿠 퍼즐 풀기
func tryPlaceNumber(ctx *context.Context, grid [][]int) bool {
	cell := getEmptyCell(grid)
	if cell == nil {
		*ctx = context.WithValue(*ctx, "result", grid)
		return true
	}

	copiedGrid := make([][]int, len(grid))
	for i := range grid {
		copiedGrid[i] = make([]int, len(grid[i]))
		copy(copiedGrid[i], grid[i])
	}

	numbers := getCandidateNumbers(copiedGrid, cell)

	for _, number := range numbers {
		select {
		case <-(*ctx).Done():
			return false
		default:
		}
		copiedGrid[cell.y][cell.x] = number
		if tryPlaceNumber(ctx, copiedGrid) {
			return true
		}
		copiedGrid[cell.y][cell.x] = 0
	}
	return false
}

// validateResult result 검증 (zeroCheck 인자를 추가하여 0을 체크할지 말지 제어)
func validateResult(grid [][]int, checkEmpty bool) bool {
	// 그리드 크기와 서브그리드 크기
	gridSize := len(grid)
	subGridSize := GridMap[gridSize]

	// 숫자가 유효 범위 내에 있는지 확인
	isValidNumber := func(n int) bool {
		return n >= 1 && n <= gridSize
	}

	// 행에 중복된 값이 없는지 확인
	for _, rowData := range grid {
		if checkEmpty && contains(rowData, 0) {
			return false
		}
		if hasDuplicates(rowData) {
			return false
		}
		// 유효 범위 검사
		for _, num := range rowData {
			if !isValidNumber(num) && num != 0 {
				return false
			}
		}
	}

	// 열에 중복된 값이 없는지 확인
	for col := 0; col < gridSize; col++ {
		columnData := make([]int, gridSize)
		for row := 0; row < gridSize; row++ {
			columnData[row] = grid[row][col]
		}
		if checkEmpty && contains(columnData, 0) {
			return false
		}
		if hasDuplicates(columnData) {
			return false
		}
		// 유효 범위 검사
		for _, num := range columnData {
			if !isValidNumber(num) && num != 0 {
				return false
			}
		}
	}

	// 서브그리드에 중복된 값이 없는지 확인
	for subGridY := 0; subGridY < subGridSize.height; subGridY++ {
		for subGridX := 0; subGridX < subGridSize.width; subGridX++ {
			currSubGrid := make([]int, 0, subGridSize.width*subGridSize.height)
			// 서브그리드 크기에 맞게 범위 설정
			for i := subGridY * subGridSize.height; i < (subGridY+1)*subGridSize.height && i < gridSize; i++ {
				for j := subGridX * subGridSize.width; j < (subGridX+1)*subGridSize.width && j < gridSize; j++ {
					currSubGrid = append(currSubGrid, grid[i][j])
				}
			}
			if checkEmpty && contains(currSubGrid, 0) {
				return false
			}
			if hasDuplicates(currSubGrid) {
				return false
			}
			// 유효 범위 검사
			for _, num := range currSubGrid {
				if !isValidNumber(num) && num != 0 {
					return false
				}
			}
		}
	}

	return true
}

func SolveSudoku(grid [][]int) ([][]int, string) {
	if !validateResult(grid, false) {
		return nil, "입력하신 데이터가 올바르지 않습니다."
	}

	cell := getEmptyCell(grid)
	numbers := getCandidateNumbers(grid, cell)

	var result [][]int

	wg := sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for _, number := range numbers {
		wg.Add(1)
		go func(number int) {
			defer wg.Done()
			tempGrid := make([][]int, len(grid))
			for i := range grid {
				tempGrid[i] = make([]int, len(grid[i]))
				copy(tempGrid[i], grid[i])
			}
			tempGrid[cell.y][cell.x] = number

			if tryPlaceNumber(&ctx, tempGrid) {
				result = ctx.Value("result").([][]int)
				cancel()
			}
		}(number)
	}

	wg.Wait()

	if result == nil || len(result) == 0 {
		return nil, "데이터 추출을 실패하였습니다."
	}

	if !validateResult(result, true) {
		return nil, "결과 검증에 실패하였습니다."
	}

	return result, ""
}
