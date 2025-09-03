package CNCService

import (
	"io"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

type PrinterBuffer struct {
	nowValue uint
	maxValue uint
	mutex    sync.Mutex
	cond     *sync.Cond
}

func (PB *PrinterBuffer) GetValueData() uint {
	return PB.nowValue
}

func (PB *PrinterBuffer) Decrement() {
	PB.mutex.Lock()
	PB.nowValue -= 1
	PB.mutex.Unlock()
}

func (PB *PrinterBuffer) Increment() {
	PB.mutex.Lock()
	if PB.nowValue < PB.maxValue {
		PB.nowValue++
		PB.cond.Signal()
	}
	PB.mutex.Unlock()
}

func (PB *PrinterBuffer) SetBufferSize(value uint) {
	PB.mutex.Lock()
	PB.nowValue = value
	PB.mutex.Unlock()
}

func (PB *PrinterBuffer) SetMaxBufferSize(value uint) {
	PB.mutex.Lock()
	PB.maxValue = value
	PB.mutex.Unlock()
}

func (PB *PrinterBuffer) GetBufferSize() uint {
	PB.mutex.Lock()
	for PB.nowValue == 0 {
		PB.cond.Wait()
	}
	value := PB.nowValue
	PB.mutex.Unlock()
	return value
}

func (PB *PrinterBuffer) WaitForNonZero() {
	PB.mutex.Lock()
	for PB.nowValue == 0 {
		PB.cond.Wait()
	}
	PB.mutex.Unlock()
}

func NewPrinterBuffer() *PrinterBuffer {
	PB := &PrinterBuffer{nowValue: 1, maxValue: 1}
	PB.cond = sync.NewCond(&PB.mutex)
	return PB
}

type Transmitter struct {
	nowValue uint
	maxValue uint
	mutex    sync.Mutex
	cond     *sync.Cond
}

func (transmitter *Transmitter) GetValueData() uint {
	return transmitter.nowValue
}

func (transmitter *Transmitter) Decrement() {
	transmitter.mutex.Lock()
	transmitter.nowValue -= 1
	transmitter.mutex.Unlock()
}

func (transmitter *Transmitter) Increment() {
	transmitter.mutex.Lock()
	if transmitter.nowValue < transmitter.maxValue {
		transmitter.nowValue++
		transmitter.cond.Signal()
	}
	transmitter.mutex.Unlock()
}

func (transmitter *Transmitter) SetBufferSize(value uint) {
	transmitter.mutex.Lock()
	transmitter.nowValue = value
	transmitter.mutex.Unlock()
}

func (transmitter *Transmitter) SetMaxBufferSize(value uint) {
	transmitter.mutex.Lock()
	transmitter.maxValue = value
	transmitter.mutex.Unlock()
}

func (transmitter *Transmitter) GetBufferSize() uint {
	transmitter.mutex.Lock()
	for transmitter.nowValue == 0 {
		transmitter.cond.Wait()
	}
	value := transmitter.nowValue
	transmitter.mutex.Unlock()
	return value
}

func (transmitter *Transmitter) WaitForNonZero() {
	transmitter.mutex.Lock()
	for transmitter.nowValue == 0 {
		transmitter.cond.Wait()
	}
	transmitter.mutex.Unlock()
}

func NewTransmitter() *Transmitter {
	transmitter := &Transmitter{nowValue: 1, maxValue: 1}
	transmitter.cond = sync.NewCond(&transmitter.mutex)
	return transmitter
}

func (transmitter *Transmitter) SyncBuffers(Connection io.ReadWriter, proto ExchangeProtocol) {
	reader := NewTimeoutReader(Connection, time.Second*2)
	Connection.Write([]byte(proto.BuildTransmitDataInt(SYNC)))
	result := reader.Read()
	if result == "" {
		return
	}
	comands := strings.Split(result, proto.Command(EndOfData))

	for _, val := range comands {
		if strings.HasPrefix(val, proto.Command(EndOfData)) {
			str, _ := strings.CutPrefix(val, proto.Command(MMaxBufferSize))
			MaxSize, err := strconv.Atoi(str)
			if err != nil {
				log.Println(err)
			}
			transmitter.SetMaxBufferSize(uint(MaxSize))
		}
	}
	Connection.Write([]byte(proto.Command(ClearBuffer)))
	reader.Read()
	transmitter.SetBufferSize(transmitter.maxValue)
}

type TimeoutReader struct {
	reader  io.Reader
	TimeOut time.Duration
	buffer  []byte
}

func NewTimeoutReader(r io.Reader, timeout time.Duration) *TimeoutReader {
	return &TimeoutReader{reader: r, TimeOut: timeout}
}

func (PR *TimeoutReader) Read() string {
	readBuf := make([]byte, 256)
	PR.buffer = PR.buffer[:0]

	timer := time.NewTimer(PR.TimeOut)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			return string(PR.buffer)
		default:
			n, err := PR.reader.Read(readBuf)
			if err != nil {
				return string(PR.buffer)
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
