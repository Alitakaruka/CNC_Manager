package main

import (
	"CNCManager/Service"
	TDServer "CNCManager/WEB/Server"
	"io"
	"log"
	"os"
	"path/filepath"
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
	if err := startLogger(logerPath); err != nil {
		log.Fatal(err)
	}
	server := TDServer.CNCServer{}
	server.InitServer(port, addr, sqlPath)
	server.Serve()
}

func startLogger(filePath string) error {
	dir := filepath.Dir(filePath)
	if dir != "." && dir != string(os.PathSeparator) {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	muliWriter := io.MultiWriter(os.Stdout, file)
	log.SetFlags(log.Ltime | log.Ldate | log.Llongfile)
	log.SetOutput(muliWriter)
	return nil
}
