package main

import (
	"distributed-parallel-game-of-life/gol"
	"fmt"
	"net"
	"os"
	"sync"
)

func main() {
	readInputArgs()
	registerAndListen()
}

func registerAndListen() {
	conn, err := net.Dial("tcp", host+":"+port)
	if err != nil {
		fmt.Printf("Can not connect to %s:%s because of following error:\n %s", host, port, err.Error())
		os.Exit(1)
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
		wait.Add(int(threadCount))
		//todo handle err
		parts, _ := board.Split(threadCount)
		for i := 0; i < int(threadCount); i++ {
			partNo := i

			go func() {
				parts[partNo] = parts[partNo].CalcNext()
				wait.Done()
			}()
		}
		wait.Wait()

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
