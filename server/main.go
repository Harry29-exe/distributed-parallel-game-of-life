package main

import (
	"distributed-parallel-game-of-life/gol"
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

const (
	host     = "localhost"
	port     = "3333"
	protocol = "tcp"
)

var lock = sync.Mutex{}
var conns = make([]net.Conn, 0, 10)

func main() {
	portListener, err := net.Listen(protocol, host+":"+port)
	if err != nil {
		fmt.Println("Could not start listening on: "+
			host+":"+port+
			", because of error:", err)
		os.Exit(1)
	}

	go listen(portListener)

	waitGroup := sync.WaitGroup{}
	board := gol.RandomBoardPart(8, 8)
	println("Input")
	board.Println()

	for {
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
		for _, part := range bParts {
			part.Println()
		}

		board = board.Merge(bParts)
		println("Result")

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
