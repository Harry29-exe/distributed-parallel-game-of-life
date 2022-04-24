package gol

import (
	"math"
	"math/rand"
	"sort"
)

func emptyBoardPart(width, height uint32) BoardPart {
	matrix := make([][]int8, width+2)
	for i := uint32(0); i < width+2; i++ {
		matrix[i] = make([]int8, height+2)
	}

	return BoardPart{
		width:  width,
		height: height,
		board:  matrix,
	}
}

func RandomBoardPart(width, height uint32) BoardPart {
	board := emptyBoardPart(width, height)
	for x := uint32(0); x < board.width+2; x++ {
		for y := uint32(0); y < board.height+2; y++ {
			board.board[x][y] = int8(rand.Intn(2))
		}
	}

	return board
}

type BoardPart struct {
	width  uint32
	height uint32
	// board indexed by [x][y] matrix of board fields,
	// matrix is of dimensions width+2 x height+2 because of
	// edges of map needed to calc next iteration
	board  [][]int8
	partNo uint32
}

func (b BoardPart) Print() {
	for y := uint32(0); y < b.height+2; y++ {
		for x := uint32(0); x < b.width+2; x++ {
			if b.board[x][y] == 0 {
				print("=")
			} else {
				print("#")
			}
		}
		println()
	}
}

func (b BoardPart) CalcNext() BoardPart {
	next := emptyBoardPart(b.width, b.height)

	for x := uint32(1); x < b.width+1; x++ {
		for y := uint32(1); y < b.height+1; y++ {
			neighbors := b.getNeighbors(x, y)

			if (b.board[x][y] == 1 && neighbors == 2) ||
				neighbors == 3 {

				next.board[x][y] = 1
			}

		}
	}

	return next
}

func (b BoardPart) getNeighbors(x, y uint32) int8 {
	sum := int8(0)
	sum += b.board[x-1][y-1]
	sum += b.board[x-1][y]
	sum += b.board[x-1][y+1]

	sum += b.board[x][y-1]
	sum += b.board[x][y+1]

	sum += b.board[x+1][y-1]
	sum += b.board[x+1][y]
	sum += b.board[x+1][y+1]

	return sum
}

// Split into sqrt(n)
func (b BoardPart) Split(n uint32) []BoardPart {
	nSqrt := uint32(math.Floor(math.Sqrt(float64(n)))) // split into grid nSqrt x nSqrt
	extraSplits := n - nSqrt*nSqrt                     // how many times split part of grid into 2

	parts := make([]BoardPart, n)
	partW, partH := b.width/nSqrt, b.height/nSqrt
	extraSplitsCounter, partNo := extraSplits, uint32(0)

	for yPart := uint32(0); yPart < nSqrt; yPart++ {
		startY := partH * yPart

		for xPart := uint32(0); xPart < nSqrt; xPart++ {
			startX := partW * xPart

			columns := make([][]int8, partW+2)

			for i := uint32(0); i < partW+2; i++ {
				columns[i] = b.board[startX+i][startY : startY+partH+2]
			}

			parts[partNo] = BoardPart{
				width:  partW,
				height: partH,
				board:  columns,
				partNo: partNo,
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
		p1.partNo = part.partNo
		p2.partNo = part.partNo + 1
		parts[i] = p1
		parts[i+1] = p2
	}

	return parts
}

func (b BoardPart) splitInto2() (BoardPart, BoardPart) {
	width, height := b.width/2, b.height/2
	b1Cols, b2Cols := make([][]int8, width), make([][]int8, width)

	for y := uint32(0); y < height+2; y++ {
		b1Cols[y] = b.board[y]
	}
	for y := height; y < height+height+2; y++ {
		b1Cols[y] = b.board[y]
	}

	return BoardPart{
			width:  width,
			height: height,
			board:  b1Cols,
			partNo: 0,
		}, BoardPart{
			width:  width,
			height: height,
			board:  b2Cols,
			partNo: 1,
		}
}

func (b BoardPart) Merge(parts []BoardPart) BoardPart {
	sort.Slice(parts, func(i, j int) bool {
		return parts[i].partNo < parts[j].partNo
	})

	merged := emptyBoardPart(b.width, b.height)
	mergedX, mergedY := uint32(0), uint32(0)

	for _, part := range parts {
		for x := uint32(1); x < part.width+1; x++ {
			for y := uint32(1); y < part.height+1; y++ {
				merged.board[mergedX+x][mergedY+y] = part.board[x][y]
			}
		}
		mergedX += part.width
		if mergedX == b.width {
			mergedX = 0
			mergedY += part.height
		}
	}

	return merged
}
