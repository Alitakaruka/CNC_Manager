package DataBase

import (
	"CNCManager/CNC"
	"database/sql"
	"errors"
	"fmt"
	"log"
)

type CNCRepository struct {
	Db *sql.DB
}

func (PR *CNCRepository) InitRepository(sqlPath string) error {
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
	uniqueKey TEXT NOT NULL,
	machineType 	 INTEGER NOT NULL,
	chip  TEXT,
	customName  TEXT,
	connectionType   TEXT,
	connectionData TEXT,
	firmwate INTEGER.
	)`)
	if ex != nil {
		return ex
	}
	res.Exec()

	res, ex = db.Prepare(`CREATE TABLE IF NOT EXIST Settings(
	id INTEGER PRIMARY KEY AUTOUNCREMENT,
	deviceId INTEGER NOT NULL,
	x_step_mm INTEGER,
	y_step_mm INTEGER,
	z_step_mm INTEGER,
	e_step_mm INTEGER,
	z_offset  INTEGER
	)`)

	if ex != nil {
		return ex
	}
	res.Exec()

	PR.Db = db
	return nil
}

func (PR *CNCRepository) AddMachine(CNC *CNC.CNCCore) error {
	if PR.Db == nil {
		return errors.New("null database")
	}

	if dto := PR.FindMachine(CNC); dto != nil {
		return nil
	}

	statement, err := PR.Db.Prepare(`INSERT INTO Machines(
		machineType,
		chip,
		customName,
		connectionType,
		connectionData,
		firmwate)
		VALUES(?,?,?,?,?,?)
	)`)

	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}
	_, err = statement.Exec(CNC.DTO.UniqueKey,
		CNC.DTO.MACHINE_TYPE,
		CNC.DTO.TARGET_MACHINE_NAME,
		CNC.DTO.ConnectionString,
		CNC.DTO.ConnectionData,
		CNC.DTO.Device_Chip_Name)
	if err != nil {
		fmt.Printf("ex: %v\n", err)
		return err
	}
	return nil
}

func (PR *CNCRepository) FindMachine(CNC *CNC.CNCCore) *CNC.CNC_DTO {
	if PR.Db == nil {
		return nil
	}

	query := `SELECT 
	machineType,
	chip,
	customName,
	connectionType,
	connectionData, 
	firmwate
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
			return nil
		}
		log.Println(err)
		return nil
	}
	return nil
}

func (PR *CNCRepository) GetAllMachines() []*CNC.CNCCore {
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
		core := CNC.CNCCore{}
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
func (PR *CNCRepository) DeletePrinter(CNC *CNC.CNCCore) error {
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
