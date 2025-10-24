package Connectors

import (
	"context"
	"io"
	"net"
	"time"
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
	ctx, cancale := context.WithTimeout(context.Background(), time.Second*2)
	defer cancale()
	dialer := net.Dialer{}
	adr := SC.Adress + ":" + SC.Port
	conn, err := dialer.DialContext(ctx, "tcp", adr)
	if err != nil {
		return err
	}
	SC.ReadWriteCloser = conn
	// URL := url.URL{Scheme: "ws", Host: adr, Path: "/ws"}
	// WS, _, err := websocket.DefaultDialer.Dial(URL.String(), nil)
	// if err != nil {
	// 	return err
	// }
	// CNCSock := CNCSocket{Conn: WS}
	// SC.ReadWriteCloser = &CNCSock
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
