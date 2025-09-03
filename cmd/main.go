package main

import (
	TDServer "PrinterManager/WEB/Server"
)

func main() {

	server := TDServer.PrinterServer{}
	server.InitServer()
	server.Serve()
}
