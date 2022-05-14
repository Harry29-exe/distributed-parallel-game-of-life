package main

import (
	"distributed-parallel-game-of-life/gol"
	"fmt"
	"net"
	"os"
	"sync"
)

const (
	host = "localhost"
	port = "3333"
)

const threadCount = 2

func main() {
	registerAndListen()
}

func registerAndListen() {
	conn, err := net.Dial("tcp", host+":"+port)
	if err != nil {
		fmt.Println("Can not connect to "+host+":"+port+" \nerror: ", err)
	}

	for {
		// receive
		board, err := gol.Remote.ReceiveBoard(conn)
		if err != nil {
			println(err.Error())
			os.Exit(1)
		}
		fmt.Println("Received board")

		// calc
		wait := sync.WaitGroup{}
		wait.Add(threadCount)
		//todo handle err
		parts, _ := board.Split(threadCount)
		for i := 0; i < threadCount; i++ {
			partNo := i
			parts[partNo].PrintWithBorder()
			print("\n")
			go func() {
				parts[partNo] = parts[partNo].CalcNext()
				wait.Done()
			}()
		}
		wait.Wait()
		print("\n\n")

		for _, part := range parts {
			part.PrintWithBorder()
			print("\n")
		}
		outputBoard := board.Merge(parts)

		fmt.Println("Calculated next board")

		err = gol.Remote.SendBoard(conn, outputBoard)
		if err != nil {
			println(err.Error())
			os.Exit(1)
		}
		fmt.Println("Send board to server")
	}

}
