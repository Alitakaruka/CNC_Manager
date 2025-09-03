package Connectors

import "io"

type CNCConnector interface {
	Connect() error
	Reconnect() (bool, error)
	GetName() string
	GetConnectionString() string
	io.ReadWriteCloser
}
