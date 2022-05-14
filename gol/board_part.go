package gol

import (
	"errors"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

func emptyBoardPart(width, height, startW, startH uint32) BoardPart {
	matrix := make([][]int8, width+2)
	for i := uint32(0); i < width+2; i++ {
		matrix[i] = make([]int8, height+2)
	}

	return BoardPart{
		Width:  width,
		Height: height,
		Board:  matrix,
		StartW: startW,
		StartH: startH,
	}
}

func RandomBoardPart(width, height uint32) BoardPart {
	board := emptyBoardPart(width, height, 0, 0)
	for x := uint32(1); x < board.Width; x++ {
		for y := uint32(1); y < board.Height; y++ {
			board.Board[x][y] = int8(rand.Intn(2))
		}
	}

	return board
}

//type BoardPartEncoder struct{}
//
//func (e BoardPartEncoder) Serialize(part BoardPart) []byte {
//
//}

type BoardPart struct {
	Width  uint32
	Height uint32
	// Board indexed by [x][y] matrix of Board fields,
	// matrix is of dimensions Width+2 x Height+2 because of
	// edges of map needed to calc next iteration
	Board  [][]int8
	StartW uint32
	StartH uint32
}

func (b BoardPart) FPrint(iteration int, file *os.File) error {
	strBuilder := strings.Builder{}
	strBuilder.WriteString(strconv.Itoa(iteration))
	strBuilder.WriteRune('\n')

	for y := uint32(1); y < b.Height+1; y++ {
		for x := uint32(1); x < b.Width+1; x++ {
			if b.Board[x][y] == 0 {
				strBuilder.WriteRune('=')
			} else {
				strBuilder.WriteRune('#')
			}
		}
		strBuilder.WriteRune('\n')
	}
	str := strBuilder.String()
	_, err := file.Write([]byte(str))
	if err != nil {
		return err
	}

	return nil
}

func (b BoardPart) Print() {
	for y := uint32(1); y < b.Height+1; y++ {
		for x := uint32(1); x < b.Width+1; x++ {
			if b.Board[x][y] == 0 {
				print("=")
			} else {
				print("#")
			}
		}
		println()
	}
}

func (b BoardPart) Println() {
	println()
	for y := uint32(1); y < b.Height+1; y++ {
		for x := uint32(1); x < b.Width+1; x++ {
			if b.Board[x][y] == 0 {
				print("=")
			} else {
				print("#")
			}
		}
		println()
	}
	println()
}

func (b BoardPart) PrintWithBorder() {
	for y := uint32(0); y < b.Height+2; y++ {
		for x := uint32(0); x < b.Width+2; x++ {
			if b.Board[x][y] == 0 {
				print("=")
			} else {
				print("#")
			}
		}
		println()
	}
}

func (b BoardPart) CalcNext() BoardPart {
	next := emptyBoardPart(b.Width, b.Height, b.StartW, b.StartH)

	for x := uint32(1); x < b.Width+1; x++ {
		for y := uint32(1); y < b.Height+1; y++ {
			neighbors := b.getNeighbors(x, y)

			if (b.Board[x][y] == 1 && neighbors == 2) ||
				neighbors == 3 {

				next.Board[x][y] = 1
			}

		}
	}

	return next
}

func (b BoardPart) getNeighbors(x, y uint32) int8 {
	sum := int8(0)
	sum += b.Board[x-1][y-1]
	sum += b.Board[x-1][y]
	sum += b.Board[x-1][y+1]

	sum += b.Board[x][y-1]
	sum += b.Board[x][y+1]

	sum += b.Board[x+1][y-1]
	sum += b.Board[x+1][y]
	sum += b.Board[x+1][y+1]

	return sum
}

func (b BoardPart) Merge(parts []BoardPart) BoardPart {
	data := make([][]int8, b.Width+2)
	for i := 0; i < int(b.Width+2); i++ {
		data[i] = make([]int8, b.Height+2)
	}

	for _, part := range parts {
		partIterator := 1
		partEnd := part.StartH + part.Height
		for x := part.StartW; x < part.StartW+part.Width; x++ {
			copy(data[x+1-b.StartW][part.StartH+1-b.StartH:partEnd+1-b.StartH], part.Board[partIterator][1:part.Height+1])
			partIterator++
		}
	}

	return BoardPart{
		Width:  b.Width,
		Height: b.Height,
		Board:  data,
		StartW: b.StartW,
		StartH: b.StartH,
	}
}

func (b BoardPart) Split(n uint32) ([]BoardPart, error) {
	parts := make([]BoardPart, n)
	parts[0] = b
	partsIteratorLen := 1
	partsIterator := 0

	for i := 1; i < int(n); i++ {
		err := b.splitInto2(parts, partsIterator, i)
		if err != nil {
			return nil, err
		}

		partsIterator++
		if partsIterator == partsIteratorLen {
			partsIterator = 0
			partsIteratorLen = i
		}
	}

	return parts, nil
}

func (b BoardPart) splitInto2(parts []BoardPart, partIndex, newPartIndex int) error {
	board := parts[partIndex]
	if board.Width > board.Height {

		b1Width := uint32(math.Ceil(float64(board.Width) / 2))
		parts[partIndex] = BoardPart{
			Width:  b1Width,
			Height: board.Height,
			Board:  b.copyData(board.StartW, board.StartH, b1Width, board.Height),
			StartW: board.StartW,
			StartH: board.StartH,
		}

		b2Width := uint32(math.Floor(float64(board.Width) / 2))
		if b2Width < 1 {
			return errors.New("can not split board")
		}
		parts[newPartIndex] = BoardPart{
			Width:  b2Width,
			Height: board.Height,
			Board:  b.copyData(board.StartW+b1Width, board.StartH, b2Width, board.Height),
			StartW: board.StartW + b1Width,
			StartH: board.StartH,
		}

	} else {
		b1Height := uint32(math.Ceil(float64(board.Height) / 2))
		parts[partIndex] = BoardPart{
			Width:  board.Width,
			Height: b1Height,
			Board:  b.copyData(board.StartW, board.StartH, board.Width, b1Height),
			StartW: board.StartW,
			StartH: board.StartH,
		}

		b2Height := uint32(math.Floor(float64(board.Height) / 2))
		if b2Height < 1 {
			return errors.New("can not split board")
		}
		b2HeightStart := board.StartH + b1Height
		parts[newPartIndex] = BoardPart{
			Width:  board.Width,
			Height: b2Height,
			Board:  b.copyData(board.StartW, b2HeightStart, board.Width, b2Height),
			StartW: board.StartW,
			StartH: b2HeightStart,
		}
	}

	return nil
}

func (b BoardPart) copyData(startX, startY uint32, width, height uint32) [][]int8 {
	data := make([][]int8, width+2)
	for x := 0; x < int(width+2); x++ {
		data[x] = make([]int8, height+2)
	}

	dataStartY := startY
	dataEndY := startY + height + 2
	dataIter := 0
	for i := startX; i < startX+width+2; i++ {
		copy(data[dataIter], b.Board[i-b.StartW][dataStartY-b.StartH:dataEndY-b.StartH])
		dataIter++
	}

	return data
}

func (b BoardPart) Equal(b2 BoardPart) bool {
	if b.Width != b2.Width || b.Height != b2.Height {
		return false
	}

	for x := 1; x < int(b.Width+1); x++ {
		for y := 1; y < int(b.Height+1); y++ {
			if b.Board[x][y] != b2.Board[x][y] {
				return false
			}
		}
	}

	return true
}
