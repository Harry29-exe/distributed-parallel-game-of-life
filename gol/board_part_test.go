package gol

import (
	"sort"
	"testing"
)

func TestBoardPart_Split(t *testing.T) {
	board := RandomBoardPart(4, 4)
	println("original")
	board.Print()
	boards := board.Split(4)

	sort.Slice(boards, func(i, j int) bool {
		return boards[i].partNo < boards[j].partNo
	})

	for i, bPart := range boards {
		println(i)
		bPart.Print()
	}

	println("Merged")
	merged := board.Merge(boards)
	merged.Print()

}

func TestRandomBoardPart(t *testing.T) {
	board := RandomBoardPart(16, 16)
	board.Print()
}
