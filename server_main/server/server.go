package server

import (
	"net"
	"sync"
)

type Server struct {
	connections net.Conn
	lock        sync.Locker
}

func (s Server) Start(protocol, host, port string) error {
	portListener, err := net.Listen(protocol, host+":"+port)
	if err != nil {
		return err
	}

}
