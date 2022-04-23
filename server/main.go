package main

import (
	"distributed-parallel-game-of-life/gol"
	"fmt"
	"net"
	"os"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "3333"
	CONN_TYPE = "tcp"
)

var connections = make([]net.Conn, 10)

func main() {
	go startServer()

	board := gol.BoardPart{}
}

func startServer() {
	portListener, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Could not start listening on: "+
			CONN_HOST+":"+CONN_PORT+
			", because of error:", err)
		os.Exit(1)
	}

	defer portListener.Close()
	for {
		conn, err := portListener.Accept()
		if err != nil {
			fmt.Println("Could not accept connection because:", err)
		}

		connections = append(connections, conn)
	}
}
