package Connectors

import (
	"strconv"

	Serial "github.com/jacobsa/go-serial/serial"
)

type SerialConnector struct {
	trackCloser
	PortName string
	BaudRate int
}

func NewSerialConnector(PortName string, BaudRate int) *SerialConnector {
	return &SerialConnector{PortName: PortName, BaudRate: BaudRate}
}

func (SC *SerialConnector) Connect() error {
	options := Serial.OpenOptions{
		PortName:              SC.PortName,
		BaudRate:              uint(SC.BaudRate),
		DataBits:              8,
		StopBits:              1,
		MinimumReadSize:       1,
		ParityMode:            Serial.PARITY_NONE,
		RTSCTSFlowControl:     false,
		InterCharacterTimeout: 100,
	}

	port, ex := Serial.Open(options)
	if ex != nil {
		return ex
	}

	SC.trackCloser.InitTracker(port)
	return nil
}

func (SC *SerialConnector) GetConnectionString() string {
	return SC.PortName + ":" + strconv.Itoa(SC.BaudRate)
}

func (SC *SerialConnector) Reconnect() error {
	if SC.ReadWriteCloser != nil {
		SC.Close()
	}

	ex := SC.Connect()
	if ex != nil {
		return ex
	}
	return nil
}

func (SC *SerialConnector) GetName() string {
	return ""
}
