package main

import (
	TDServer "CNCManager/WEB/Server"
)

func main() {

	server := TDServer.PrinterServer{}
	server.InitServer()
	server.Serve()
}
