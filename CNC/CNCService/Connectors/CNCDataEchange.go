package Connectors

import (
	"CNCManager/CNC/CNCService"
	"bufio"
	"errors"
	"io"
	"sync"
	"time"
)

const (
	MsgTypeText  = 1
	MsgTypePing  = 2
	MsgTypePong  = 3
	MsgTypeClose = 4
)
const BitShiftMsgTypeBits = 5
const BitShiftBytesLen = 3

type Exchanger struct {
	io.ReadWriteCloser
	LastFrameACK     bool
	ConnectionString string
	mut              *sync.Mutex
	TimeoutReader    CNCService.TimeoutReader
}

func NewExchanger(conn io.ReadWriteCloser, ConnectionString string) *Exchanger {
	return &Exchanger{ReadWriteCloser: conn,
		ConnectionString: ConnectionString,
		TimeoutReader:    *CNCService.NewTimeoutReader(conn, time.Second)}
}

func (ex *Exchanger) Ping() error {
	var TypeOfMassage uint8 = 0
	TypeOfMassage |= MsgTypePing << BitShiftMsgTypeBits
	ex.mut.Lock()
	ex.Write([]byte{TypeOfMassage})
	bytes := ex.TimeoutReader.ReadBytes()
	if len(bytes) == 0 {
		ex.ReadWriteCloser.Close()
		return errors.New("connection is closed")
	}
	ex.mut.Unlock()
	return nil
}

func (ex *Exchanger) Pong() error {
	var TypeOfMassage uint8 = 0
	TypeOfMassage |= MsgTypePong << BitShiftMsgTypeBits
	ex.mut.Lock()
	ex.Write([]byte{TypeOfMassage})
	ex.mut.Unlock()
	return nil
}

func (ex *Exchanger) ReadMessage() (string, int) {
	byt := ex.ReadBytes(1)[0]
	MsgType := byt >> BitShiftMsgTypeBits
	leaghtBytes := (byt & 0b00011000) >> BitShiftBytesLen

	bt := ex.ReadBytes(uint32(leaghtBytes))
	var res uint32 = 0
	for _, val := range bt {
		res = res<<8 | uint32(val)
	}
	msgData := ex.ReadBytes(res)
	switch MsgType {
	case MsgTypeText:
		return string(msgData), MsgTypeText
	case MsgTypePing:
		ex.Pong()
	case MsgTypeClose:
		ex.ReadWriteCloser.Close()
	default:
		ex.ReadWriteCloser.Close()
	}
	ex.mut.Unlock()
	return "", 1
}

func (ex *Exchanger) WriteMessage(msg []byte) error {
	bytesLen := ex.GetMsgLen(msg)
	prefix := (0 | MsgTypeText<<BitShiftMsgTypeBits) | len(bytesLen)<<BitShiftBytesLen
	StartMsg := make([]byte, len(bytesLen)+1)
	copy(StartMsg[1:], bytesLen)
	StartMsg[0] = byte(prefix)
	ex.Write(StartMsg)
	ex.Write(msg)
	for !ex.LastFrameACK {

	}
	return nil
}

func (ex *Exchanger) ReadBytes(n uint32) []byte {
	var result []byte
	bufio := bufio.NewReader(ex.ReadWriteCloser)
	var i uint32 = 0
	for ; i < n; i++ {
		data, err := bufio.ReadByte()
		if err != nil {
			return result
		}
		result = append(result, data)
	}
	return result
}

func (ex *Exchanger) GetMsgLen(msg []byte) []byte {
	MsgLen := len(msg)
	var byte1 byte = byte(MsgLen)
	var byte2 byte = byte(MsgLen >> 8)
	var byte3 byte = byte(MsgLen >> 16)
	res := make([]byte, 0)
	res = append(res, byte1)
	if byte2 != 0 {
		res = append(res, byte2)
	}
	if byte3 != 0 {
		res = append(res, byte3)
	}
	return res
}
