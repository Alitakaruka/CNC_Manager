package main

import (
	"CNCManager/Service"
	TDServer "CNCManager/WEB/Server"
	"io"
	"log"
	"os"
)

func main() {
	config := Service.GetConfig("config.yaml")
	var port, addr, sqlPath, logerPath string
	if config == nil {
		port = "8080"
		addr = "0.0.0.0"
		sqlPath = "CNCManagerDB.db"
		logerPath = "Logs.log"
	} else {
		port = config.Server.Port
		addr = config.Server.Addr
		sqlPath = config.Database.Path
		logerPath = config.Logs.Path
	}
	startLoger(logerPath)
	server := TDServer.CNCServer{}
	server.InitServer(port, addr, sqlPath)
	server.Serve()
}

func startLoger(filePath string) {
	file, ex := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if ex != nil {
		log.Fatal(ex)
	}
	muliWriter := io.MultiWriter(os.Stdout, file)
	log.SetFlags(log.Ltime | log.Ldate | log.Llongfile)
	log.SetOutput(muliWriter)
}
