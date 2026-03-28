package Connectors

import (
	"fmt"
	"io"
	"sync"
)

type CNCConnector interface {
	Connect() error
	Reconnect() error
	io.ReadWriteCloser
	WaitClosed()
}

// trackCloser оборачивает io.ReadWriteCloser и сигнализирует через канал о закрытии
type trackCloser struct {
	io.ReadWriteCloser
	closed   chan struct{} // сигнал закрытия
	once     sync.Once     // чтобы закрытие сработало только один раз
	isClosed bool
	Rmu      sync.RWMutex
	Wmu      sync.RWMutex
}

func (t *trackCloser) InitTracker(rwc io.ReadWriteCloser) {
	t.closed = make(chan struct{})
	t.ReadWriteCloser = rwc
}

// Read проксирует Read
func (t *trackCloser) Read(p []byte) (int, error) {
	t.Rmu.Lock()
	n, err := t.ReadWriteCloser.Read(p)
	t.Rmu.Unlock()
	return n, err
}

// Write проксирует Write
func (t *trackCloser) Write(p []byte) (int, error) {

	t.Wmu.Lock()
	n, err := t.ReadWriteCloser.Write(p)
	t.Wmu.Unlock()
	return n, err
}

// Close проксирует Close и закрывает канал только один раз
func (t *trackCloser) Close() error {
	fmt.Println("Close")
	var err error
	err = t.ReadWriteCloser.Close()
	fmt.Println("Rw closer")
	t.once.Do(func() {
		close(t.closed) // уведомляем всех, кто ждет
	})
	return err
}

// WaitClosed блокирует до момента закрытия
func (t *trackCloser) WaitClosed() {
	<-t.closed
}
