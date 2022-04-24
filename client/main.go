package main

import (
	"distributed-parallel-game-of-life/gol"
	"encoding/binary"
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

	for {
		inputLenBytes := make([]byte, 4)
		_, err = conn.Read(inputLenBytes)
		if err != nil {
			println("Error reading input inputLen, err:", err.Error())
			return
		}
		inputLen := binary.LittleEndian.Uint32(inputLenBytes)
		println("Buffer len:", inputLen)

		input := make([]byte, inputLen)
		l, err := conn.Read(input)
		if err != nil {
			println("Error reading input , err:", err.Error())
			return
		} else if uint32(l) != inputLen {
			println("Payload has different inputLen that expected."+
				"\nExpected:", inputLen, "Actual:", l)
			return
		}

		board := gol.DeserializeBoardPart(input)
		board = board.CalcNext()
		println()
		board.Print()
		println()

		output := gol.SerializeBoardPart(board)
		outputLen := make([]byte, 4)
		binary.LittleEndian.PutUint32(outputLen, uint32(len(output)))

		conn.Write(outputLen)
		conn.Write(output)
	}

}
