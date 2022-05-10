package server

import (
	"encoding/binary"
	"fmt"
	"math"
	"net"
	"sync"
)

type DistributorServer interface {
	Start() error
	ConnectionCount() uint
	Distribute(tasks []Task) error
	Stop() error
	IsActive() bool
}

type Task struct {
	TaskData []byte
	Receiver func(data []byte, err error) error
}

type taskError struct {
	err             error
	isConnectionErr bool
	taskId          int
}

func newTaskErr(taskId int, err error) taskError {
	return taskError{
		err:             err,
		isConnectionErr: false,
		taskId:          taskId,
	}
}

func newTaskConnErr(taskId int, err error) taskError {
	return taskError{
		err:             err,
		isConnectionErr: true,
		taskId:          taskId,
	}
}

func CreateServer(protocol, host, port string) DistributorServer {
	return &distributorServer{
		connections: make([]net.Conn, 0, 10),
		lock:        &sync.Mutex{},
		stopChan:    make(chan bool),
		protocol:    protocol,
		host:        host,
		port:        port,
	}
}

type distributorServer struct {
	connections          []net.Conn
	lock                 sync.Locker
	stopChan             chan bool
	protocol, host, port string
}

func (s *distributorServer) Start() error {
	portListener, err := net.Listen(s.protocol, s.host+":"+s.port)
	if err != nil {
		return err
	}

	go s.listen(portListener, s.stopChan)

	return nil
}

func (s *distributorServer) ConnectionCount() uint {
	return uint(len(s.connections))
}

func (s *distributorServer) Distribute(tasks []Task) error {
	s.lock.Lock()
	errChan := make(chan taskError)
	waitGroup := &sync.WaitGroup{}
	waitChan := make(chan bool)
	connCounter := 0
	connsLen := len(s.connections)

	for i := 0; i < len(tasks); i++ {
		waitGroup.Add(1)
		s.sendTask(tasks[i], i, s.connections[connCounter], waitGroup, errChan)
		connCounter = (connCounter + 1) % connsLen
	}

	go func() {
		waitGroup.Wait()
		close(waitChan)
	}()

	select {
	case err := <-errChan:
		if err.isConnectionErr {

		}
	case <-waitChan:
		return nil
	}
}

func (s *distributorServer) Stop() error {

}

func (s *distributorServer) IsActive() bool {

}

func (s distributorServer) listen(portListener net.Listener, stop chan bool) {
	defer func(portListener net.Listener) {
		err := portListener.Close()
		if err != nil {
			fmt.Printf("When closing server port error has occured: %s", err.Error())
		}
	}(portListener)

	for {
		select {
		case <-stop:
			fmt.Printf("Closing connection listener\n")
			return
		default:
			conn, err := portListener.Accept()
			if err != nil {
				//todo this will be displayed when Close() will be called
				fmt.Println("Could not accept connection because:", err)
			}

			s.lock.Lock()
			s.connections = append(s.connections, conn)
			s.lock.Unlock()
		}
	}
}

func (s *distributorServer) sendTask(task Task, taskId int, connId, conn net.Conn, group *sync.WaitGroup, errChan chan taskError) {
	if len(task.TaskData) > math.MaxUint32 {
		errChan <- newTaskErr(taskId, fmt.Errorf("task data is of length %d but max protocl msg is %d",
			len(task.TaskData), math.MaxUint32))

		return
	}

	msgLen := dataLen(task.TaskData)
	_, err := conn.Write(msgLen)
	if err != nil {
		errChan <- newTaskConnErr(taskId, err)
		return
	}

	_, err = conn.Write(task.TaskData)
	if err != nil {
		errChan <- newTaskConnErr(taskId, err)
		return
	}

	n, err := conn.Read(msgLen)
	if err != nil {
		errChan <- newTaskConnErr(taskId, err)
		return
	} else if n != 4 {
		errChan <- newTaskConnErr(taskId, fmt.Errorf("distributor server expected 4 byte but get %d", n))
		return
	}
	incomingDataLen := binary.LittleEndian.Uint32(msgLen)

	msg := make([]byte, incomingDataLen)
	n, err = conn.Read(msg)
	if err != nil {
		errChan <- newTaskConnErr(taskId, err)
		return
	} else if n > math.MaxUint32 {
		errChan <- newTaskErr(taskId, fmt.Errorf("received msg of leght %d, that is protocol error", n))
		return
	} else if uint32(n) != incomingDataLen {
		errChan <- newTaskConnErr(taskId, fmt.Errorf("msg should be %d bytes long but received %d bytes",
			incomingDataLen, n))
		return
	}

	unexpectedErr := task.Receiver(msg, nil)
	if unexpectedErr != nil {
		errChan <- newTaskErr(taskId, unexpectedErr)
		return
	}

	group.Done()

	return
}

func dataLen(data []byte) (length []byte) {
	length = make([]byte, 4)
	binary.LittleEndian.PutUint32(length, uint32(len(data)))

	return
}
