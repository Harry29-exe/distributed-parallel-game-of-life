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

func FuzzBoardPart_Split(f *testing.F) {
	f.Add(4, 4, 2)
	f.Add(7, 1, 7)
	f.Add(10, 1, 10)
	f.Add(5, 6, 4)
	f.Add(8, 8, 4)
	f.Add(5, 5, 1)
	f.Add(5, 5, 5)
	f.Add(4, 4, 5)

	f.Fuzz(func(t *testing.T, width, height, nSplits int) {
		board := RandomBoardPart(uint32(width), uint32(height))
		boards, err := board.Split(uint32(nSplits))
		if err != nil {
			t.Errorf("fuzz args: (%d, %d, %d) board Split() returned following error: %s", width, height, nSplits, err.Error())
		}

		merged := board.Merge(boards)

		if !board.Equal(merged) {
			t.Error("pre merge board is not equal to merged board")
		}
	})
}

func TestSerializeBoardPart(t *testing.T) {
	board := RandomBoardPart(4, 4)

	data, _ := SerializeBoardPart(board)
	newBoard, _ := DeserializeBoardPart(data)

	if !board.Equal(*newBoard) {
		t.Error("board after serialization-deserialization is not equal to original board")
	}
}

func FuzzSerializeBoardPart(f *testing.F) {
	f.Add(uint32(4), uint32(4))
	f.Add(uint32(16), uint32(16))
	f.Add(uint32(35), uint32(35))

	f.Fuzz(func(t *testing.T, width, height uint32) {
		board := RandomBoardPart(width, height)

		data, _ := SerializeBoardPart(board)
		newBoard, _ := DeserializeBoardPart(data)

		if !board.Equal(*newBoard) {
			t.Error("board after serialization-deserialization is not equal to original board")
		}
	})
}
