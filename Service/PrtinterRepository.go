package Service

import (
	"CNCManager/CNC"
	"database/sql"
	"log"
)

type PrinterRepository struct {
	Db *sql.DB
}

func (PR *PrinterRepository) InitRepository(sqlPath string) error {
	db, ex := sql.Open("sqlite", "file:"+sqlPath)
	if ex != nil {
		PR.Db = nil
		return ex
	}
	if ex = db.Ping(); ex != nil {
		PR.Db = nil
		return ex
	}
	log.Println("Connection to database successful")

	res, ex := db.Prepare(`CREATE TABLE IF NOT EXISTS Printers(
	ID INTEGER PRIMARY KEY AUTOINCREMENT,
	PrinterName TEXT NOT NULL,
	TypeOfPrinter TEXT NOT NULL,
	UserNamePrinter TEXT NOT NULL,
	TypeOfConnection TEXT NOT NULL,
	ConnectionData TEXT NOT NULL,
	Version TEXT NOT NULL
	)`)
	if ex != nil {
		return ex
	}
	res.Exec()
	PR.Db = db
	return nil
}

func (PR *PrinterRepository) AddPrinter(CNC CNC.AnyCNC) error {
	printerData := CNC.GetDTO()
	query := `INSERT INTO Printers (
		PrinterName,
		UserNamePrinter,
		TypeOfConnection,
		ConnectionData,
		Version
		) VALUES(?,?,?,?,?)`
	statement, ex := PR.Db.Prepare(query)
	if ex != nil {
		return ex
	}
	_, ex = statement.Exec(printerData.TARGET_MACHINE_NAME,
		printerData.FIRMWARE_VERSION,
		printerData.ConnectionData,
		printerData.FIRMWARE_VERSION)
	if ex != nil {
		return ex
	}
	return nil
}

// TODO
func (PR *PrinterRepository) DeletePrinter(CNC CNC.AnyCNC) error {
	printerData := CNC.GetDTO()
	query := `DELETE INTO Printers (
		PrinterName,
		UserNamePrinter,
		TypeOfConnection,
		ConnectionData,
		Version
		) VALUES(?,?,?,?,?)`
	statement, ex := PR.Db.Prepare(query)
	if ex != nil {
		return ex
	}
	_, ex = statement.Exec(printerData.TARGET_MACHINE_NAME,
		printerData.ConnectionData,
		printerData.FIRMWARE_VERSION)
	if ex != nil {
		return ex
	}
	return nil
}

// func (PR *PrinterRepository) GetAllPrinters() ([]CNC.AnyCNC, error) {
// 	if PR.Db == nil {
// 		return nil, errors.New("data base is not connected")
// 	}
// 	rows, ex := PR.Db.Query(`SELECT
// 	ID,
// 	PrinterName,
// 	UserNamePrinter,
// 	TypeOfConnection,
// 	ConnectionData,
// 	Version
// 	FROM printers`)
// 	if ex != nil {
// 		return nil, ex
// 	}
// 	var result []CNC.AnyCNC
// 	for rows.Next() {
// 		var ID int
// 		var PrinterName, UserNamePrinter, TypeOfPrinter, TypeOfConnection, ConnectionData, Version string
// 		ex := rows.Scan(&ID, &PrinterName, &UserNamePrinter, &TypeOfPrinter, &TypeOfConnection, &ConnectionData, &Version)

// 		MainData := ThreeDPrinter.PrinterDTO{PrinterName: PrinterName,
// 			UserPrinterName:  UserNamePrinter,
// 			PrinterType:      TypeOfPrinter,
// 			Version:          Version,
// 			TypeOfConnection: TypeOfConnection, ConnectionData: ConnectionData}
// 		if ex != nil {
// 			return nil, ex
// 		}
// 		if Constr, ok := ThreeDPrinter.RegisteredPrinters[MainData.PrinterName]; ok {
// 			Printer := Constr()
// 			result = append(result, Printer)
// 		}
// 	}
// 	return result, nil
// }
