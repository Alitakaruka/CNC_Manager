package Connectors

import "io"

type CNCConnector interface {
	Connect() error
	Reconnect() error
	io.ReadWriteCloser
}
