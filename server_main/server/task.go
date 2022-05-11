package server

type Task struct {
	TaskData []byte
	Receiver func(data []byte, err error) error
}

type taskErrorFactory struct {
	taskId       int
	connectionId uint64
}

func (f taskErrorFactory) newTaskErr(err error) taskError {
	return taskError{
		err:             err,
		isConnectionErr: false,
		taskId:          f.taskId,
		connectionId:    f.connectionId,
	}
}

func (f taskErrorFactory) newTaskConnError(err error) taskError {
	return taskError{
		err:             err,
		isConnectionErr: true,
		taskId:          f.taskId,
		connectionId:    f.connectionId,
	}
}

type taskError struct {
	err             error
	isConnectionErr bool
	taskId          int
	connectionId    uint64
}
