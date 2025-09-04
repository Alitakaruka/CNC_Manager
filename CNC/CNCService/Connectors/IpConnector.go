package Connectors

import (
	"errors"
	"io"
	"net/url"

	"github.com/gorilla/websocket"
)

type IPConnector struct {
	io.ReadWriteCloser
	Adress string
	Port   string
}

func NewIpConnector(Adress, port string) *IPConnector {
	return &IPConnector{Adress: Adress, Port: port}
}

func (SC *IPConnector) Connect() error {
	adr := SC.Adress + ":" + SC.Port
	URL := url.URL{Scheme: "ws", Host: adr, Path: "/ws"}
	WS, _, err := websocket.DefaultDialer.Dial(URL.String(), nil)
	if err != nil {
		return err
	}
	CNCSock := CNCSocket{Conn: WS}
	SC.ReadWriteCloser = &CNCSock
	return nil
}

func (SC *IPConnector) GetConnectionString() string {
	return SC.Adress + ":" + SC.Port
}
func (SC *IPConnector) Reconnect() (bool, error) {
	if SC.ReadWriteCloser != nil {
		SC.Close()
	}

	ex := SC.Connect()
	if ex != nil {
		return false, ex
	}
	return true, nil
}
func (SC *IPConnector) GetName() string {
	return ""
}

const (
	MessageTypeText   = 1
	MessageTypeBinaty = 2
	MessageTypeClose  = 8
	MessageTypePing   = 9
	MessageTypePong   = 10
)

type CNCSocket struct {
	*websocket.Conn
	reader io.Reader
}

func (CNC_Soc *CNCSocket) Read(buffer []byte) (int, error) {
	// Если ридер пустой — читаем новое сообщение
	for CNC_Soc.reader == nil {
		_, r, err := CNC_Soc.NextReader()
		if err != nil {
			return 0, err
		}
		CNC_Soc.reader = r
	}

	// Читаем данные из текущего сообщения
	n, err := CNC_Soc.reader.Read(buffer)
	if errors.Is(err, io.EOF) {
		// сообщение закончилось — сбрасываем ридер, чтобы взять следующее
		CNC_Soc.reader = nil
		if n > 0 {
			return n, nil
		}
		return 0, io.EOF
	}
	return n, err
}

func (CNC_Soc *CNCSocket) Write(buffer []byte) (int, error) {
	w, err := CNC_Soc.Conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return 0, err
	}
	n, err := w.Write(buffer)
	if err != nil {
		return n, err
	}
	err = w.Close()
	return n, err
}
