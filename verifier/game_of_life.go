package main

import "fmt"

type Board struct {
	board  [][]int8
	width  uint
	height uint
}

func BoardsFromStr(str []rune) ([]Board, error) {
	boards := make([]Board, 0, 40)

	currentBoard := make([][]int8, 1)
	boardY := 0

	for i := parseNoRow(str, 0) + 1; i < len(str); i++ {
		switch str[i] {
		case '\n':
			if i+1 == len(str) {
				boards = append(boards, boardFrom2dArray(currentBoard))
				return boards, nil
			}

			if str[i+1] == '=' || str[i+1] == '#' {
				boardY++
				currentBoard = append(currentBoard, make([]int8, 0, 40))
			} else {
				boards = append(boards, boardFrom2dArray(currentBoard))
				currentBoard = make([][]int8, 1)
				boardY = 0
				i = parseNoRow(str, i+1)
			}
		case '=':
			currentBoard[boardY] = append(currentBoard[boardY], 0)
		case '#':
			currentBoard[boardY] = append(currentBoard[boardY], 1)
		default:
			return nil, fmt.Errorf("invalid file, gotten %c", str[i])
		}
	}

	return boards, nil
}

func (b Board) Equals(b2 Board) bool {
	if b.width != b2.width || b.height != b2.height {
		return false
	}

	for y := 0; y < int(b.height); y++ {
		for x := 0; x < int(b.width); x++ {
			if b.board[y][x] != b2.board[y][x] {
				return false
			}
		}
	}

	return true
}

func (b Board) CalcNext() Board {
	newBoard := make([][]int8, b.height)
	for y := 0; y < int(b.height); y++ {
		newBoard[y] = make([]int8, b.width)
	}

	for y := uint(0); y < b.height; y++ {
		for x := uint(0); x < b.width; x++ {
			neighbors := b.neighborsCount(x, y)
			if (b.board[y][x] == 1 && neighbors == 2) ||
				neighbors == 3 {

				newBoard[y][x] = 1
			}
		}
	}

	return Board{
		board:  newBoard,
		width:  b.width,
		height: b.height,
	}
}

func (b Board) neighborsCount(x, y uint) int8 {
	sum := int8(0)
	if x > 0 && y > 0 {
		sum += b.board[y-1][x-1]
	}
	if y > 0 {
		sum += b.board[y-1][x]
	}
	if y > 0 && x+1 < b.width {
		sum += b.board[y-1][x+1]
	}
	if y+1 < b.height {
		sum += b.board[y+1][x]
	}
	if y+1 < b.height && x+1 < b.width {
		sum += b.board[y+1][x+1]
	}
	if y+1 < b.height && x > 0 {
		sum += b.board[y+1][x-1]
	}
	if x+1 < b.width {
		sum += b.board[y][x+1]
	}
	if x > 0 {
		sum += b.board[y][x-1]
	}

	return sum
}

func boardFrom2dArray(data [][]int8) Board {
	return Board{
		board:  data,
		width:  uint(len(data[0])),
		height: uint(len(data)),
	}
}

func parseNoRow(str []rune, i int) int {
	for {
		if str[i] == '\n' {
			return i
		}

		i++
	}
}
