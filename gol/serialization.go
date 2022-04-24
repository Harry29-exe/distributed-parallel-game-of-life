package gol

import (
	"bytes"
	"encoding/gob"
	"os"
)

func SerializeBoardPart(part BoardPart) []byte {
	payloadBuff := bytes.Buffer{}

	encoder := gob.NewEncoder(&payloadBuff)
	err := encoder.Encode(part)
	if err != nil {
		println("Error encoding BoardPart object, err", err.Error())
		os.Exit(1)
	}

	return payloadBuff.Bytes()
}

func DeserializeBoardPart(data []byte) BoardPart {
	buffer := bytes.Buffer{}
	buffer.Write(data)
	decoder := gob.NewDecoder(&buffer)

	board := BoardPart{}
	err := decoder.Decode(&board)
	if err != nil {
		println("Could not decode board payload, err", err.Error())
		os.Exit(1)
	}

	return board
}
