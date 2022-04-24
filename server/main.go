package main

import (
	"bytes"
	"distributed-parallel-game-of-life/gol"
	"encoding/binary"
	"encoding/gob"
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
var connections = make([]net.Conn, 0, 10)

func main() {
	portListener, err := net.Listen(protocol, host+":"+port)
	if err != nil {
		fmt.Println("Could not start listening on: "+
			host+":"+port+
			", because of error:", err)
		os.Exit(1)
	}

	go listen(portListener)

	for {
		board := gol.RandomBoardPart(4, 4)
		board.Print()
		b := bytes.Buffer{}

		encoder := gob.NewEncoder(&b)
		err := encoder.Encode(board)
		if err != nil {
			println("Error encoding BoardPart object, err", err.Error())
			return
		}

		for len(connections) == 0 {
			time.Sleep(1 * time.Second)
		}
		bufferLen := make([]byte, 4)
		binary.LittleEndian.PutUint32(bufferLen, uint32(b.Len()))
		connections[0].Write(bufferLen)
		connections[0].Write(b.Bytes())

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
		connections = append(connections, conn)
		lock.Unlock()
	}
}
