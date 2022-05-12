package server

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"net"
	"sync"
	"time"
)

type DistributorServer interface {
	Start() error
	ConnectionCount() uint64
	Distribute(tasks []Task) error
	Stop() error
	IsActive() bool
}

func CreateServer(protocol, host, port string) DistributorServer {
	return &distributorServer{
		connections: map[uint64]ConnWrapper{},
		connCount:   0,
		connCounter: &ConnCounter{nextId: 0},
		lock:        &sync.Mutex{},
		stopChan:    make(chan bool),
		isActive:    false,
		protocol:    protocol,
		host:        host,
		port:        port,
	}
}

type distributorServer struct {
	connections map[uint64]ConnWrapper
	connCount   uint64
	connCounter *ConnCounter

	lock                 sync.Locker
	portListener         net.Listener
	stopChan             chan bool
	isActive             bool
	protocol, host, port string
}

func (s *distributorServer) Start() error {
	s.lock.Lock()
	if s.isActive {
		return errors.New("server already started")
	}
	s.lock.Unlock()

	portListener, err := net.Listen(s.protocol, s.host+":"+s.port)
	if err != nil {
		return err
	}

	go s.listen(portListener, s.stopChan)

	s.lock.Lock()
	s.portListener = portListener
	s.isActive = true
	s.lock.Unlock()

	fmt.Println("Server started")

	return nil
}

func (s *distributorServer) ConnectionCount() uint64 {
	return uint64(len(s.connections))
}

func (s *distributorServer) Distribute(tasks []Task) error {
	s.lock.Lock()
	if s.isActive == false {
		return errors.New("server is not active")
	} else if len(s.connections) == 0 {
		s.lock.Unlock()
		s.waitForConnection()
		s.lock.Lock()
	}

	errChan := make(chan taskError)
	waitGroup := &sync.WaitGroup{}
	waitChan := make(chan bool)

	taskCounter := 0
	for taskCounter < len(tasks) {
		for _, connWrapper := range s.connections {
			if taskCounter > len(tasks) {
				break
			}

			waitGroup.Add(1)
			s.sendTask(tasks[taskCounter], taskCounter, connWrapper, waitGroup, errChan)
			taskCounter++
		}
	}
	s.lock.Unlock()

	go func() {
		waitGroup.Wait()
		close(waitChan)
	}()

	for {
		select {
		case err := <-errChan:
			if err.isConnectionErr {
				s.lock.Lock()
				delete(s.connections, err.connectionId)
				s.lock.Unlock()

				if len(s.connections) == 0 {
					s.waitForConnection()
				}
				for _, connWrapper := range s.connections {
					s.sendTask(tasks[err.taskId], err.taskId, connWrapper, waitGroup, errChan)
				}

			} else {
				return errors.New("got unexpected error" + err.err.Error())
			}
		case <-waitChan:
			return nil
		}
	}
}

func (s *distributorServer) Stop() error {
	s.lock.Lock()
	s.stopChan <- true
	for _, connWrapper := range s.connections {
		err := connWrapper.Conn.Close()
		if err != nil {
			return err
		}
	}

	s.isActive = false
	s.lock.Unlock()

	return nil
}

func (s *distributorServer) IsActive() bool {
	return s.isActive
}

func (s distributorServer) listen(portListener net.Listener, stop chan bool) {
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

			fmt.Println("New connection")
			connId := s.connCounter.GetAndIncrease()
			s.connections[connId] = ConnWrapper{
				Conn:   conn,
				connId: connId,
			}
			s.connCount++

			s.lock.Unlock()
		}
	}
}

func (s *distributorServer) waitForConnection() {
	fmt.Println("Waiting for connections...")
	for {
		s.lock.Lock()
		if len(s.connections) > 0 {
			s.lock.Unlock()
			return
		}
		s.lock.Unlock()

		time.Sleep(10 * time.Millisecond)
	}
}

func (s *distributorServer) sendTask(task Task, taskId int, connWrapper ConnWrapper, group *sync.WaitGroup, errChan chan taskError) {
	errFactory := taskErrorFactory{
		taskId:       taskId,
		connectionId: connWrapper.connId,
	}
	conn := connWrapper.Conn

	if len(task.TaskData) > math.MaxUint32 {
		errChan <- errFactory.newTaskErr(fmt.Errorf("task data is of length %d but max protocl msg is %d",
			len(task.TaskData), math.MaxUint32))

		return
	}

	msgLen := dataLen(task.TaskData)
	_, err := conn.Write(msgLen)
	if err != nil {
		errChan <- errFactory.newTaskConnError(err)
		return
	}

	_, err = conn.Write(task.TaskData)
	if err != nil {
		errChan <- errFactory.newTaskConnError(err)
		return
	}

	n, err := conn.Read(msgLen)
	if err != nil {
		errChan <- errFactory.newTaskConnError(err)
		return
	} else if n != 4 {
		errChan <- errFactory.newTaskConnError(fmt.Errorf("distributor server expected 4 byte but get %d", n))
		return
	}
	incomingDataLen := binary.LittleEndian.Uint32(msgLen)

	msg := make([]byte, incomingDataLen)
	n, err = conn.Read(msg)
	if err != nil {
		errChan <- errFactory.newTaskConnError(err)
		return
	} else if n > math.MaxUint32 {
		errChan <- errFactory.newTaskErr(fmt.Errorf("received msg of leght %d, that is protocol error", n))
		return
	} else if uint32(n) != incomingDataLen {
		errChan <- errFactory.newTaskConnError(fmt.Errorf("msg should be %d bytes long but received %d bytes",
			incomingDataLen, n))
		return
	}

	unexpectedErr := task.Receiver(msg, nil)
	if unexpectedErr != nil {
		errChan <- errFactory.newTaskErr(unexpectedErr)
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
