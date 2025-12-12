package Connectors

import (
	"context"
	"fmt"
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
	fmt.Println("Ip connector!")
	ctx, cancale := context.WithTimeout(context.Background(), time.Second*5)
	defer cancale()
	dialer := net.Dialer{}
	adr := SC.Adress + ":" + SC.Port
	conn, err := dialer.DialContext(ctx, "tcp", adr)
	if err != nil {
		return err
	}
	SC.ReadWriteCloser = conn
	return nil
}

func (SC *IPConnector) GetConnectionString() string {
	return SC.Adress + ":" + SC.Port
}
func (SC *IPConnector) Reconnect() error {
	if SC.ReadWriteCloser != nil {
		SC.Close()
	}
	ex := SC.Connect()
	if ex != nil {
		fmt.Printf("ex: %v\n", ex)
		return ex
	}
	return nil
}
func (SC *IPConnector) GetName() string {
	return ""
}
