package server

import (
	"net"
	"sync/atomic"
)

type ConnWrapper struct {
	Conn   net.Conn
	connId uint64
}

type ConnCounter struct {
	nextId *uint64
}

func (c ConnCounter) GetAndIncrease() uint64 {
	return atomic.AddUint64(c.nextId, 1) - 1
}
