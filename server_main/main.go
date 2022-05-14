package main

import (
	"distributed-parallel-game-of-life/gol"
	"fmt"
	"os"
	"server/server"
	"time"
)

var iteration = 0

func main() {
	readInputArgs()
	printArgs()

	//create file
	err := os.WriteFile(outputFilePath, make([]byte, 0), 0755)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	file, err := os.Create(outputFilePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	s := server.CreateServer(protocol, host, port)
	err = s.Start()
	if err != nil {
		fmt.Println("Could not start server because of following error: ", err.Error())
	}

	board := gol.RandomBoardPart(boardW, boardH)
	err = board.FPrint(iteration, file)
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}

	for iteration = 1; iteration <= programIterations; iteration++ {
		taskCount := min(max(s.ConnectionCount(), 1), uint64(boardW*boardH))
		tasks := make([]server.Task, taskCount)

		//todo add error handling when board can not be split to so many parts
		boardParts, err := board.Split(uint32(taskCount))
		if err != nil {

		}

		for i := uint64(0); i < taskCount; i++ {
			fmt.Printf("i = %d\n", i)
			boardData, err := gol.SerializeBoardPart(boardParts[i])
			if err != nil {
				fmt.Println("could not serialize board, that is unexpected error please contact it teams")
				os.Exit(1)
			}

			iCopy := i
			tasks[i] = server.Task{
				TaskData: boardData,
				Receiver: func(data []byte, err error) error {
					if err != nil {
						return err
					}
					boardPart, err := gol.DeserializeBoardPart(data)
					if err != nil {
						return err
					}
					fmt.Printf("iCopy = %d\n", iCopy)
					boardParts[iCopy] = *boardPart
					return nil
				},
			}
		}

		err = s.Distribute(tasks)
		if err != nil {
			fmt.Println("Could not distribute tasks between clients because of following error: " + err.Error())
			os.Exit(1)
		}
		board = board.Merge(boardParts)

		err = board.FPrint(iteration, file)
		if err != nil {
			println(err.Error())
			os.Exit(1)
		}

		if delay {
			time.Sleep(delayTime)
		}
	}
}

func max(a, b uint64) uint64 {
	if a > b {
		return a
	}
	return b
}

func min(a, b uint64) uint64 {
	if a < b {
		return a
	}
	return b
}
