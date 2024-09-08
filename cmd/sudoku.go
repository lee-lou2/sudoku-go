package cmd

import (
	"context"
	"sort"
	"sync"
)

type SubGrid struct {
	width, height int
}

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

// getEmptyCells 0인 좌표와 0의 개수를 반환
func getEmptyCells(grid [][]int, subGrid SubGrid) []Cell {
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

	// emptyCnt 기준으로 정렬
	sort.Slice(cellInfoList, func(i, j int) bool {
		return cellInfoList[i].cnt < cellInfoList[j].cnt
	})

	emptyCells := make([]Cell, 0)
	for _, cellInfo := range cellInfoList {
		emptyCells = append(emptyCells, cellInfo.cell)
	}

	return emptyCells
}

// getCandidateNumbers 해당 좌표에 들어갈 수 있는 숫자를 반환
func getCandidateNumbers(grid [][]int, cell Cell, subGrid SubGrid) []int {
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
			currSubGrid = append(currSubGrid, grid[i][j]) // 슬라이스에 값 추가
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
func tryPlaceNumber(ctx *context.Context, grid [][]int, emptyCells []Cell, subGrid SubGrid) bool {
	cell := emptyCells[0]
	emptyCells = emptyCells[1:]
	copiedEmptyCells := make([]Cell, len(emptyCells))
	copy(copiedEmptyCells, emptyCells)

	copiedGrid := make([][]int, len(grid))
	for i := range grid {
		copiedGrid[i] = make([]int, len(grid[i]))
		copy(copiedGrid[i], grid[i])
	}

	numbers := getCandidateNumbers(copiedGrid, cell, subGrid)

	if len(numbers) == 0 {
		return false
	}

	for _, number := range numbers {
		select {
		case <-(*ctx).Done():
			return false
		default:
		}

		copiedGrid[cell.y][cell.x] = number
		if len(copiedEmptyCells) == 0 {
			*ctx = context.WithValue(*ctx, "result", copiedGrid)
			return true
		}

		if tryPlaceNumber(ctx, copiedGrid, copiedEmptyCells, subGrid) {
			return true
		}

		copiedGrid[cell.y][cell.x] = 0
	}
	return false
}

// validateResult result 검증 (zeroCheck 인자를 추가하여 0을 체크할지 말지 제어)
func validateResult(grid [][]int, subGrid SubGrid, checkEmpty bool) bool {
	// 행에 중복된 값이 없는지 확인
	for _, rowData := range grid {
		if checkEmpty && contains(rowData, 0) {
			return false
		}
		if hasDuplicates(rowData) {
			return false
		}
	}

	// 열에 중복된 값이 없는지 확인
	for i := 0; i < subGrid.width*subGrid.height; i++ {
		if checkEmpty && contains(grid[i], 0) {
			return false
		}
		if hasDuplicates(grid[i]) {
			return false
		}
	}

	// 박스에 중복된 값이 없는지 확인
	for subGridY := 0; subGridY < subGrid.height; subGridY++ {
		for subGridX := 0; subGridX < subGrid.width; subGridX++ {
			currSubGrid := make([]int, 0, subGrid.width*subGrid.height)
			// 박스 크기에 맞게 범위 수정
			for i := subGridY * subGrid.height; i < (subGridY+1)*subGrid.height && i < len(grid); i++ {
				for j := subGridX * subGrid.width; j < (subGridX+1)*subGrid.width && j < len(grid[i]); j++ {
					currSubGrid = append(currSubGrid, grid[i][j])
				}
			}
			if checkEmpty && contains(currSubGrid, 0) {
				return false
			}
			if hasDuplicates(currSubGrid) {
				return false
			}
		}
	}

	return true
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

func SolveSudoku(grid [][]int) ([][]int, string) {
	subGrid := GridMap[len(grid)]

	if !validateResult(grid, subGrid, false) {
		return nil, "입력하신 데이터가 올바르지 않습니다."
	}

	emptyCells := getEmptyCells(grid, subGrid)
	cell := emptyCells[0]
	emptyCells = emptyCells[1:]
	numbers := getCandidateNumbers(grid, cell, subGrid)

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

			if tryPlaceNumber(&ctx, tempGrid, emptyCells, subGrid) {
				result = ctx.Value("result").([][]int)
				cancel()
			}
		}(number)
	}

	wg.Wait()

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

			if tryPlaceNumber(&ctx, tempGrid, emptyCells, subGrid) {
				result = ctx.Value("result").([][]int)
				cancel()
			}
		}(number)
	}

	if result == nil || len(result) == 0 {
		return nil, "데이터 추출을 실패하였습니다."
	}

	if !validateResult(result, subGrid, true) {
		return nil, "결과 검증에 실패하였습니다."
	}

	return result, ""
}
