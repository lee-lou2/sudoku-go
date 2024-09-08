package api

import (
	"log"
	"net/http"
)

// RunServer 서버 실행
func RunServer() {
	// 서버 생성
	mux := http.NewServeMux()

	// 정적 파일 서버 등록
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/", http.StripPrefix("/", fs))

	// 핸들러 등록
	mux.Handle("/v1/sudoku/solve", http.HandlerFunc(solveSudokuHandler))

	// 서버 실행
	log.Println("Server is running on :3000")
	_ = http.ListenAndServe(":3000", mux)
}
