package gol

import (
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
)

func emptyBoardPart(width, height uint32) BoardPart {
	matrix := make([][]int8, width+2)
	for i := uint32(0); i < width+2; i++ {
		matrix[i] = make([]int8, height+2)
	}

	return BoardPart{
		Width:  width,
		Height: height,
		Board:  matrix,
	}
}

func RandomBoardPart(width, height uint32) BoardPart {
	board := emptyBoardPart(width, height)
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
	PartNo uint32
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
	next := emptyBoardPart(b.Width, b.Height)

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

// Split into sqrt(n)
func (b BoardPart) Split(n uint32) []BoardPart {
	nSqrt := uint32(math.Floor(math.Sqrt(float64(n)))) // split into grid nSqrt x nSqrt
	extraSplits := n - nSqrt*nSqrt                     // how many times split part of grid into 2

	parts := make([]BoardPart, n)
	partW, partH := b.Width/nSqrt, b.Height/nSqrt
	extraSplitsCounter, partNo := extraSplits, uint32(0)

	for yPart := uint32(0); yPart < nSqrt; yPart++ {
		startY := partH * yPart

		for xPart := uint32(0); xPart < nSqrt; xPart++ {
			startX := partW * xPart

			columns := make([][]int8, partW+2)

			for i := uint32(0); i < partW+2; i++ {
				columns[i] = b.Board[startX+i][startY : startY+partH+2]
			}

			parts[partNo] = BoardPart{
				Width:  partW,
				Height: partH,
				Board:  columns,
				PartNo: partNo,
			}

			if extraSplitsCounter > 0 {
				partNo += 2
				extraSplitsCounter--
			} else {
				partNo++
			}
		}
	}

	for i := uint32(0); i < extraSplits; i += 2 {
		part := parts[i]
		p1, p2 := part.splitInto2()
		p1.PartNo = part.PartNo
		p2.PartNo = part.PartNo + 1
		parts[i] = p1
		parts[i+1] = p2
	}

	return parts
}

func (b BoardPart) splitInto2() (BoardPart, BoardPart) {
	width, height := b.Width/2, b.Height
	b1Cols, b2Cols := make([][]int8, width+2), make([][]int8, width+2)

	for x := uint32(0); x < width+2; x++ {
		b1Cols[x] = b.Board[x]
		b2Cols[x] = b.Board[x+width]
	}

	return BoardPart{
			Width:  width,
			Height: height,
			Board:  b1Cols,
			PartNo: 0,
		}, BoardPart{
			Width:  width,
			Height: height,
			Board:  b2Cols,
			PartNo: 1,
		}
}

func (b BoardPart) Merge(parts []BoardPart) BoardPart {
	sort.Slice(parts, func(i, j int) bool {
		return parts[i].PartNo < parts[j].PartNo
	})

	merged := emptyBoardPart(b.Width, b.Height)
	mergedX, mergedY := uint32(0), uint32(0)

	for _, part := range parts {
		for x := uint32(1); x < part.Width+1; x++ {
			for y := uint32(1); y < part.Height+1; y++ {
				merged.Board[mergedX+x][mergedY+y] = part.Board[x][y]
			}
		}
		mergedX += part.Width
		if mergedX == b.Width {
			mergedX = 0
			mergedY += part.Height
		}
	}

	return merged
}
