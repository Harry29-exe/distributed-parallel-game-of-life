package main

import (
	"distributed-parallel-game-of-life/gol"
	"encoding/binary"
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

	for {
		board := gol.RandomBoardPart(4, 4)
		println()
		board.Print()
		println()

		for len(conns) == 0 {
			time.Sleep(1 * time.Second)
		}

		connsLen := len(conns)
		bParts := board.Split(uint32(connsLen))
		for i := 0; i < connsLen; i++ {
			go delegateBoardPart(conns[i], &bParts[i], &waitGroup)
			waitGroup.Add(1)
		}

		waitGroup.Wait()
		board = board.Merge(bParts)
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

func delegateBoardPart(conn net.Conn, part *gol.BoardPart, gr *sync.WaitGroup) {
	requestPayload := gol.SerializeBoardPart(*part)
	payloadLen := make([]byte, 4)
	binary.LittleEndian.PutUint32(payloadLen, uint32(len(requestPayload)))

	conn.Write(payloadLen)
	conn.Write(requestPayload)

	conn.Read(payloadLen)
	responsePayload := make([]byte, binary.LittleEndian.Uint32(payloadLen))
	conn.Read(responsePayload)

	b := gol.DeserializeBoardPart(responsePayload)
	part = &b

	gr.Done()
}
