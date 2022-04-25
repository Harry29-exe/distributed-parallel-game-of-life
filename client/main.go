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

const threadCount = 4

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

		// calc
		wait := sync.WaitGroup{}
		wait.Add(threadCount)
		parts := board.Split(threadCount)
		for i := 0; i < threadCount; i++ {
			partNo := i
			go func() {
				parts[partNo] = parts[partNo].CalcNext()
				wait.Done()
			}()
		}
		wait.Wait()
		outputBoard := board.Merge(parts)

		err = gol.Remote.SendBoard(conn, outputBoard)
		if err != nil {
			println(err.Error())
			os.Exit(1)
		}
	}

}
