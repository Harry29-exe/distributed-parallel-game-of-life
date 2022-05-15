package gol

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
)

var Remote remote = remote{}

type remote struct{}

func (r remote) SendBoard(conn net.Conn, board BoardPart) error {
	serialized, err := SerializeBoardPart(board)
	if err != nil {
		return err
	}
	serializedLen := make([]byte, 4)
	binary.LittleEndian.PutUint32(serializedLen, uint32(len(serialized)))
	_, err = conn.Write(serializedLen)
	if err != nil {
		return err
	}
	_, err = conn.Write(serialized)
	if err != nil {
		return err
	}

	return nil
}

func (r remote) ReceiveBoard(conn net.Conn) (*BoardPart, error) {
	boardLen := make([]byte, 4)
	i, err := conn.Read(boardLen)
	if err != nil {
		return nil, err
	} else if i != 4 {
		return nil, errors.New("incoming data's length array should be " +
			"exactly 4 bytes long")
	}

	length := binary.LittleEndian.Uint32(boardLen)
	println(length)
	boardData := make([]byte, length)
	i, err = conn.Read(boardData)
	if err != nil {
		return nil, err
	} else if uint32(i) != length {
		return nil,
			fmt.Errorf("data should be exactly the length of previously send length (%d bytes) but is %d bytes long",
				length, i)
	}

	board, err := DeserializeBoardPart(boardData)
	if err != nil {
		return nil, err
	}

	return board, nil
}
