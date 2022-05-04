package gol_test

import (
	"distributed-parallel-game-of-life/gol"
	"errors"
	"fmt"
	"testing"
)

func TestBoardPart_Split(t *testing.T) {
	board := gol.RandomBoardPart(36, 36)
	board.Println()

	for i := 4; i < 36; i++ {
		println(i)

		boardParts := board.Split(4)
		for _, part := range boardParts {
			fmt.Println("Part no.", part.PartNo, " w:", part.Width, " h:", part.Height)
		}
		merged := board.Merge(boardParts)

		err := cmpBoards(board, merged)
		if err != nil {
			t.Error(err.Error())
		}
	}
}

func TestRandomBoardPart(t *testing.T) {
	board := gol.RandomBoardPart(16, 16)
	board.PrintWithBorder()
}

func TestSerializeBoardPart(t *testing.T) {
	board := gol.RandomBoardPart(4, 4)
	println()
	board.Print()
	println()

	data, _ := gol.SerializeBoardPart(board)
	newBoard, _ := gol.DeserializeBoardPart(data)
	println()
	newBoard.Print()
	println()
}

func cmpBoards(b1, b2 gol.BoardPart) error {
	if b1.Width != b2.Width || b1.Height != b2.Height {
		return errors.New("boards b1, b2 have different sizes")
	}

	for x := uint32(0); x < b1.Width; x++ {
		for y := uint32(0); y < b1.Height; y++ {
			if b1.Board[x][y] != b2.Board[x][y] {
				return fmt.Errorf("cells at x:%d, y:%d are different, value in b1: %d, in b2: %d",
					x, y, b1.Board[x][y], b2.Board[x][y])
			}

		}
	}

	return nil
}
