package gol

import (
	"testing"
)

func TestBoardPart_Split(t *testing.T) {
	widths := []uint32{4, 15, 16, 25, 32}
	heights := []uint32{4, 13, 16, 16, 13}
	ns := []uint32{2, 4, 5, 4, 8}

	for i := 0; i < len(widths); i++ {
		board := RandomBoardPart(widths[i], heights[i])
		boards, err := board.Split(ns[i])
		if err != nil {
			t.Error("board Split() returned following error: " + err.Error())
		}

		merged := board.Merge(boards)

		if !board.Equal(merged) {
			t.Error("pre merge board is not equal to merged board")
		}
	}

}

func TestRandomBoardPart(t *testing.T) {
	board := RandomBoardPart(16, 16)
	board.PrintWithBorder()
}

func TestSerializeBoardPart(t *testing.T) {
	board := RandomBoardPart(4, 4)
	println()
	board.Print()
	println()

	data, _ := SerializeBoardPart(board)
	newBoard, _ := DeserializeBoardPart(data)
	println()
	newBoard.Print()
	println()
}
