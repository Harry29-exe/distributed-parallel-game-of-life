package gol

import (
	"sort"
	"testing"
)

func TestBoardPart_Split(t *testing.T) {
	board := RandomBoardPart(4, 4)
	println("original")
	board.PrintWithBorder()
	boards := board.Split(4)

	sort.Slice(boards, func(i, j int) bool {
		return boards[i].PartNo < boards[j].PartNo
	})

	for i, bPart := range boards {
		println(i)
		bPart.PrintWithBorder()
	}

	println("Merged")
	merged := board.Merge(boards)
	merged.PrintWithBorder()

}

func TestRandomBoardPart(t *testing.T) {
	board := RandomBoardPart(16, 16)
	board.PrintWithBorder()
}
