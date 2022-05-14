package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("please enter path to file as first program argument")
		os.Exit(1)
	}

	filepath := os.Args[1]
	fileContent, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Printf("can not read file with following filepath: %s becouse of following  error: %e",
			filepath, err)
		os.Exit(1)
	}

	fileAsStr := string(fileContent)
	boards, err := BoardsFromStr([]rune(fileAsStr))
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Printf("Loaded %d boards from file %s\n", len(boards), filepath)

	fmt.Println("Checking boards...")
	board := boards[0]
	for i := 1; i < len(boards); i++ {
		board = board.CalcNext()
		if !boards[i].Equals(board) {
			fmt.Println("Boards in given file are invalid")
			os.Exit(1)
		}
	}
	fmt.Println("Boards in given file are correct")
}
