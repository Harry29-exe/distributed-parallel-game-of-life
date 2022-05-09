package server

import "net"

type Connection struct {
	conn net.Conn
}
