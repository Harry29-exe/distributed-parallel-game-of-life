package main

import (
	"bytes"
	"distributed-parallel-game-of-life/gol"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"net"
)

const (
	host = "localhost"
	port = "3333"
)

func main() {
	registerAndListen()
}

func registerAndListen() {
	conn, err := net.Dial("tcp", host+":"+port)
	if err != nil {
		fmt.Println("Can not connect to "+host+":"+port+" \nerror: ", err)
	}

	bufferLen := make([]byte, 4)
	_, err = conn.Read(bufferLen)
	if err != nil {
		println("Error reading payload length, err:", err)
		return
	}
	length := binary.LittleEndian.Uint32(bufferLen)
	println("Buffer len:", length)

	payload := make([]byte, length)
	l, err := conn.Read(payload)
	if err != nil {
		println("Error reading payload , err:", err.Error())
		return
	} else if uint32(l) != length {
		println("Payload has different length that expected."+
			"\nExpected:", length, "Actual:", l)
		return
	}

	buffer := bytes.Buffer{}
	buffer.Write(payload)
	decoder := gob.NewDecoder(&buffer)

	board := gol.BoardPart{}
	err = decoder.Decode(&board)
	if err != nil {
		println("Could not decode board payload, err", err.Error())
		return
	}

	board.Print()
}
