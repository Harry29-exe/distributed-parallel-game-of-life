package gol

import "math"

type BoardPart interface {
	Width() uint32
	Height() uint32
	CalcNext() BoardPart
	Split(n uint32) []BoardPart
}

func emptyBoardPart(width, height uint32) boardPart {
	matrix := make([][]int8, height+2)
	for i := uint32(0); i < height+2; i++ {
		matrix[i] = make([]int8, width+2)
	}

	return boardPart{
		width:  width,
		height: height,
		board:  matrix,
	}
}

type boardPart struct {
	width  uint32
	height uint32
	// board indexed by [x][y] matrix of board fields,
	// matrix is of dimensions width+2 x height+2 because of
	// edges of map needed to calc next iteration
	board      [][]int8
	partNumber uint32
}

func (b boardPart) Width() uint32 {
	return b.width
}

func (b boardPart) Height() uint32 {
	return b.height
}

func (b boardPart) CalcNext() BoardPart {
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

func (b boardPart) getNeighbors(x, y uint32) int8 {
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
func (b boardPart) Split(n uint32) []BoardPart {
	nSqrt := uint32(math.Floor(math.Sqrt(float64(n)))) // split into grid nSqrt x nSqrt
	extraSplits := n - nSqrt*nSqrt                     // how many times split part of grid into 2

	parts := make([]BoardPart, n)
	partW, partH := b.width/nSqrt, b.height/nSqrt
	extraSplitsCounter, partNo := extraSplits, uint32(0)

	for xPart := uint32(0); xPart < nSqrt; xPart++ {
		startX := partW * xPart
		for yPart := uint32(0); yPart < nSqrt; yPart++ {
			startY := partH * yPart
			columns := make([][]int8, partW+2)

			for i := uint32(0); i < partW+2; i++ {
				columns[i] = b.board[startX+i][startY : startY+partH+2]
			}

			parts[partNo] = boardPart{
				width:      partW,
				height:     partH,
				board:      columns,
				partNumber: partNo,
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
		part := parts[i].(boardPart)
		p1, p2 := part.splitInto2()
		p1.partNumber = part.partNumber
		p2.partNumber = part.partNumber + 1
		parts[i] = p1
		parts[i+1] = p2
	}

	return parts
}

func (b boardPart) splitInto2() (boardPart, boardPart) {
	width, height := b.width/2, b.height/2
	b1Cols, b2Cols := make([][]int8, width), make([][]int8, width)

	for y := uint32(0); y < height+2; y++ {
		b1Cols[y] = b.board[y]
	}
	for y := height; y < height+height+2; y++ {
		b1Cols[y] = b.board[y]
	}

	return boardPart{
			width:      width,
			height:     height,
			board:      b1Cols,
			partNumber: 0,
		}, boardPart{
			width:      width,
			height:     height,
			board:      b2Cols,
			partNumber: 1,
		}
}
