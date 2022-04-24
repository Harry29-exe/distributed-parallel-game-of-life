package main

import (
	"distributed-parallel-game-of-life/gol"
	"time"
)

func main() {
	board := gol.RandomBoardPart(5, 5)
	for {
		board = board.CalcNext()
		println()
		board.Print()
		println()

		time.Sleep(time.Second)
	}
}
