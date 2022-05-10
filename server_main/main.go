package main

import (
	"distributed-parallel-game-of-life/gol"
	"fmt"
	"net"
	"os"
	"server/server"
	"sync"
	"time"
)

var lock = sync.Mutex{}
var conns = make([]net.Conn, 0, 10)
var iteration = 0

func main() {
	readInputArgs()
	s := server.server{}

	//start server
	portListener, err := net.Listen(protocol, host+":"+port)
	if err != nil {
		fmt.Println("Could not start listening on: "+
			host+":"+port+
			", because of error:", err)
		os.Exit(1)
	}

	go listen(portListener)

	//create file
	err = os.WriteFile(outputFilePath, make([]byte, 0), 0755)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	file, err := os.Create(outputFilePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	board := gol.RandomBoardPart(boardW, boardH)
	err = board.FPrint(iteration, file)

	// infinite loop
	waitGroup := sync.WaitGroup{}

	if err != nil {
		panic(err.Error())
	}
	iteration++

	for iteration < programIterations {
		for len(conns) == 0 {
			time.Sleep(1 * time.Second)
		}

		connsLen := len(conns)
		bParts := board.Split(uint32(connsLen))
		for i := 0; i < connsLen; i++ {
			waitGroup.Add(1)
			go delegateBoardPart(bParts, i, &waitGroup)
		}
		waitGroup.Wait()

		board = board.Merge(bParts)
		err := board.FPrint(iteration, file)
		if err != nil {
			println(err.Error())
			os.Exit(1)
		}

		iteration++

		time.Sleep(1 * time.Second)
	}

}

func listen(portListener net.Listener) {
	defer portListener.Close()
	for {
		conn, err := portListener.Accept()
		if err != nil {
			fmt.Println("Could not accept connection because:", err)
		}

		lock.Lock()
		conns = append(conns, conn)
		lock.Unlock()
	}
}

func delegateBoardPart(parts []gol.BoardPart, i int, gr *sync.WaitGroup) {
	err := gol.Remote.SendBoard(conns[i], parts[i])
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}

	board, err := gol.Remote.ReceiveBoard(conns[i])
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}

	parts[i] = *board
	gr.Done()
}
