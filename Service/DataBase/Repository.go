package DataBase

import (
	"CNCManager/CNC"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sync"
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

	res, ex := db.Prepare(`CREATE TABLE IF NOT EXISTS Machines(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	UniqueKey     	 TEXT NOT NULL,
	MACHINE_TYPE 	 INTEGER NOT NULL,
	TARGET_MACHINE_NAME  TEXT NOT NULL,
	Device_Chip_Name  TEXT NOT NULL,
	ConnectionString   TEXT NOT NULL,
	ConnectionData TEXT NOT NULL
	)`)
	if ex != nil {
		return ex
	}
	res.Exec()

	// //Later
	// res, err := db.Prepare(`CREATE TABLE IF NOT EXIST CNCSettings(
	// 	FOREIGN KEY(CNC_id) REFERENCES Machines(id),
	// 	Width INTEGER NOT NULL,
	// 	Length INTEGER NOT NULL,
	// 	Height INTEGER NOT NULL,

	// )`)

	// if err != nil {
	// 	log.Println(err)
	// 	return err
	// }
	PR.Db = db
	return nil
}

func (PR *PrinterRepository) AddMachine(CNC *CNC.CNCCore) error {
	if PR.Db == nil {
		return errors.New("null database")
	}
	if PR.FindMachine(CNC) {
		return nil
	}
	query := `INSERT INTO Machines (
		UniqueKey,
		MACHINE_TYPE,
		Device_Chip_Name,
		TARGET_MACHINE_NAME,
		ConnectionString,
		ConnectionData
		) VALUES(?,?,?,?,?,?)`
	statement, ex := PR.Db.Prepare(query)
	if ex != nil {
		fmt.Printf("ex: %v\n", ex)
		return ex
	}
	_, ex = statement.Exec(CNC.DTO.UniqueKey,
		CNC.DTO.MACHINE_TYPE,
		CNC.DTO.TARGET_MACHINE_NAME,
		CNC.DTO.ConnectionString,
		CNC.DTO.ConnectionData,
		CNC.DTO.Device_Chip_Name)
	if ex != nil {
		fmt.Printf("ex: %v\n", ex)
		return ex
	}
	return nil
}

func (PR *PrinterRepository) FindMachine(CNC *CNC.CNCCore) bool {
	if PR.Db == nil {
		return false
	}

	query := `SELECT 
	UniqueKey,
	MACHINE_TYPE  as Type,
	TARGET_MACHINE_NAME  as Name,
	ConnectionString,
	ConnectionData 
	from Machines
	where MACHINE_TYPE = (?) 
	and TARGET_MACHINE_NAME = (?) 
	and ConnectionString = (?)
	and ConnectionData = (?)`

	row := PR.Db.QueryRow(query,
		CNC.DTO.MACHINE_TYPE,
		CNC.DTO.TARGET_MACHINE_NAME,
		CNC.DTO.ConnectionString,
		CNC.DTO.ConnectionData,
	)
	var dum int

	err := row.Scan(&dum)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
		log.Println(err)
		return false
	}
	return true
}

func (PR *PrinterRepository) GetAllMachines() []*CNC.CNCCore {
	var result = make([]*CNC.CNCCore, 0)

	if PR.Db == nil {
		return result
	}

	query := `SELECT 
        UniqueKey,
        MACHINE_TYPE,
        TARGET_MACHINE_NAME,
        ConnectionString,
        ConnectionData,
		Device_Chip_Name
    FROM Machines`

	rows, err := PR.Db.Query(query)

	if err != nil {
		if err != sql.ErrNoRows {
			log.Println(err)
		}
		log.Println(err) //TODO
		return result
	}

	var Scan struct {
		UniqueKey           string
		MACHINE_TYPE        int
		TARGET_MACHINE_NAME string
		ConnectionString    string
		ConnectionData      string
		Device_Chip_Name    string
	}

	for rows.Next() {
		err := rows.Scan(&Scan.UniqueKey,
			&Scan.MACHINE_TYPE,
			&Scan.TARGET_MACHINE_NAME,
			&Scan.ConnectionString,
			&Scan.ConnectionData,
			&Scan.Device_Chip_Name)

		if err != nil {
			log.Println(err)
			return result
		}
		core := CNC.CNCCore{Mutex: &sync.RWMutex{}}
		core.DTO.UniqueKey = Scan.UniqueKey
		core.DTO.MACHINE_TYPE = Scan.MACHINE_TYPE
		core.DTO.TARGET_MACHINE_NAME = Scan.TARGET_MACHINE_NAME
		core.DTO.ConnectionString = Scan.ConnectionString
		core.DTO.ConnectionData = Scan.ConnectionData
		core.DTO.Device_Chip_Name = Scan.Device_Chip_Name

		result = append(result, &core)
	}
	return result
}

// TODO
func (PR *PrinterRepository) DeletePrinter(CNC *CNC.CNCCore) error {
	if PR.Db == nil {
		return errors.New("")
	}

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
