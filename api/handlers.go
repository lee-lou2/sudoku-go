package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sudoku-go/cmd"
)

type SudokuRequest struct {
	Grid [][]int `json:"grid"`
}

type SudokuResponse struct {
	Solution [][]int `json:"solution"`
}

func solveSudokuHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer func() { _ = r.Body.Close() }()

	var sudokuReq SudokuRequest
	err = json.Unmarshal(body, &sudokuReq)
	if err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}
	solution, msg := cmd.SolveSudoku(sudokuReq.Grid)
	if solution == nil {
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	response := SudokuResponse{
		Solution: solution,
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
