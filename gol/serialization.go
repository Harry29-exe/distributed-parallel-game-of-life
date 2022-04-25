package gol

import (
	"bytes"
	"encoding/gob"
)

func SerializeBoardPart(part BoardPart) ([]byte, error) {
	payloadBuff := bytes.Buffer{}

	encoder := gob.NewEncoder(&payloadBuff)
	err := encoder.Encode(part)
	if err != nil {
		return nil, err
	}

	return payloadBuff.Bytes(), nil
}

func DeserializeBoardPart(data []byte) (*BoardPart, error) {
	buffer := bytes.Buffer{}
	buffer.Write(data)
	decoder := gob.NewDecoder(&buffer)

	board := &BoardPart{}
	err := decoder.Decode(&board)
	if err != nil {
		return nil, err
	}

	return board, nil
}
