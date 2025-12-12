package CNCService

import (
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Transmitter struct {
	CurrentFreeBytes int
	MaxBytes         int
	mutex            sync.Mutex
	cond             *sync.Cond
	Commands         []int
}

func (T *Transmitter) Trainsmit(bytes int) {
	T.mutex.Lock()
	T.CurrentFreeBytes -= bytes
	T.Commands = append(T.Commands, bytes)
	T.mutex.Unlock()
}

func (T *Transmitter) ACK() {
	if len(T.Commands) == 0 {
		fmt.Print("ACK LEN = 0!!!!!!")
		return
	}
	// fmt.Printf("T.CurrentFreeBytes: %v\n", T.CurrentFreeBytes)
	T.mutex.Lock()
	// fmt.Println("ACK")
	val := T.Commands[0]
	T.CurrentFreeBytes += val
	T.Commands = T.Commands[1:]
	T.cond.Signal()
	T.mutex.Unlock()
}

func (T *Transmitter) SetMaxBytes(MaxBytes int) {
	T.mutex.Lock()
	T.MaxBytes = MaxBytes
	T.CurrentFreeBytes = MaxBytes
	T.mutex.Unlock()
}

func (T *Transmitter) SyncBuffers(Connection io.ReadWriter) {
	T.mutex.Lock()
	defer T.mutex.Unlock()
	reader := NewTimeoutReader(Connection, time.Second*2)
	Connection.Write([]byte(SYNC + EndOfData))
	result := reader.Read()
	if result == "" {
		return
	}
	comands := strings.Split(result, EndOfData)

	for _, val := range comands {
		if strings.HasPrefix(val, MyBufferLen) {
			str, _ := strings.CutPrefix(val, MyBufferLen)
			MaxSize, err := strconv.Atoi(str)
			if err != nil {
				log.Println(err)
			}
			T.MaxBytes = MaxSize
		}
	}
	Connection.Write([]byte(ClearBuffer))
	reader.Read()
	T.CurrentFreeBytes = (T.MaxBytes)
}

func (transmitter *Transmitter) Wait(bytes int) {
	transmitter.mutex.Lock()
	for transmitter.CurrentFreeBytes < bytes {
		transmitter.cond.Wait()
	}
	transmitter.mutex.Unlock()
}

func NewTransmitter() *Transmitter {
	transmitter := &Transmitter{MaxBytes: 128, CurrentFreeBytes: 128}
	transmitter.cond = sync.NewCond(&transmitter.mutex)
	transmitter.Commands = make([]int, 0)
	return transmitter
}

type TimeoutReader struct {
	reader  io.Reader
	TimeOut time.Duration
	buffer  []byte
}

func NewTimeoutReader(r io.Reader, timeout time.Duration) *TimeoutReader {
	return &TimeoutReader{reader: r, TimeOut: timeout}
}

func (PR *TimeoutReader) ReadBytes() []byte {
	readBuf := make([]byte, 256)
	PR.buffer = PR.buffer[:0]

	timer := time.NewTimer(PR.TimeOut)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			return PR.buffer
		default:
			n, err := PR.reader.Read(readBuf)
			if err != nil {
				return PR.buffer
			}

			if n > 0 {
				PR.buffer = append(PR.buffer, readBuf[:n]...)

				if !timer.Stop() {
					<-timer.C // clear chan
				}
				timer.Reset(PR.TimeOut)
			}
		}
	}
}

func (PR *TimeoutReader) Read() string {
	readBuf := make([]byte, 256)
	PR.buffer = PR.buffer[:0]
	for {
		timer := time.NewTimer(PR.TimeOut)
		n, err := PR.readWithTimeout(readBuf, timer)
		timer.Stop()
		if err != nil {
			return string(PR.buffer)
		}

		if n > 0 {
			PR.buffer = append(PR.buffer, readBuf[:n]...)
			continue
		}
		return string(PR.buffer)
	}
}

func (PR *TimeoutReader) readWithTimeout(buf []byte, timer *time.Timer) (int, error) {
	done := make(chan struct {
		n   int
		err error
	})

	go func() {
		n, err := PR.reader.Read(buf)
		done <- struct {
			n   int
			err error
		}{n, err}
	}()

	select {
	case result := <-done:
		return result.n, result.err
	case <-timer.C:
		return 0, nil // timeout
	}
}

// type CNCBuffer
