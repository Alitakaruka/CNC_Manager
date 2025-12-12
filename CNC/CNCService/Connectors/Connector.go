package Connectors

import "io"

type CNCConnector interface {
	Connect() error
	Reconnect() error
	GetName() string
	GetConnectionString() string
	io.ReadWriteCloser
}
