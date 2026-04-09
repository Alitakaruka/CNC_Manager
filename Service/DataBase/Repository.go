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
	uniqueKey TEXT NOT NULL UNIQUE,
	machineType 	 INTEGER NOT NULL,
	chip  TEXT,
	customName  TEXT,
	connectionType   TEXT,
	connectionData TEXT,
	firmwate INTEGER)`)

	if ex != nil {
		log.Println(ex)
		return ex
	}
	_, ex = res.Exec()

	if ex != nil {
		log.Println(ex)
	}

	res, ex = db.Prepare(`CREATE TABLE IF NOT EXISTS Settings(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	deviceId INTEGER NOT NULL,
	x_step_mm INTEGER,
	y_step_mm INTEGER,
	z_step_mm INTEGER,
	e_step_mm INTEGER,
	z_offset  INTEGER
	)`)

	if ex != nil {
		log.Println(ex)
		return ex
	}
	_, ex = res.Exec()

	if ex != nil {
		log.Println(ex)
	}
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
		uniqueKey,
		machineType,
		chip,
		customName,
		connectionType,
		connectionData,
		firmwate)
		VALUES(?,?,?,?,?,?,?)`)

	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}
	_, err = statement.Exec(CNC.DTO.UniqueKey,
		CNC.DTO.MACHINE_TYPE,
		CNC.DTO.Device_Chip_Name,
		CNC.DTO.TARGET_MACHINE_NAME,
		CNC.DTO.ConnectionType,
		CNC.DTO.ConnectionData,
		"")
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
	where machineType = (?) 
	and customName = (?) 
	and connectionType = (?)
	and connectionData = (?)`

	row := PR.Db.QueryRow(query,
		CNC.DTO.MACHINE_TYPE,
		CNC.DTO.TARGET_MACHINE_NAME,
		CNC.DTO.ConnectionType,
		CNC.DTO.ConnectionData,
	)
	dto := CNC.DTO
	err := row.Scan(&dto.MACHINE_TYPE,
		&dto.Device_Chip_Name,
		&dto.TARGET_MACHINE_NAME,
		&dto.ConnectionType,
		&dto.ConnectionData,
		&dto.FIRMWARE)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		log.Println(err)
		return nil
	}
	return &dto
}

func (PR *CNCRepository) FindByKey(Key string) *CNC.CNC_DTO {
	row := PR.Db.QueryRow(`SELECT 
	machineType,
	chip,
	customName,
	connectionType,
	connectionData
	from Machines 
	where uniqueKey = (?)`, Key)

	DTO := CNC.CNC_DTO{}
	err := row.Scan(&DTO.MACHINE_TYPE,
		&DTO.Device_Chip_Name,
		&DTO.TARGET_MACHINE_NAME,
		&DTO.ConnectionType,
		&DTO.ConnectionData)

	if err != nil {
		log.Println(err)
		return nil
	}
	return &DTO
}

func (PR *CNCRepository) GetAllMachines() []CNC.CNC_DTO {
	// id INTEGER PRIMARY KEY AUTOINCREMENT,
	// uniqueKey TEXT NOT NULL,
	// machineType 	 INTEGER NOT NULL,
	// chip  TEXT,
	// customName  TEXT,
	// connectionType   TEXT,
	// connectionData TEXT,
	// firmwate INTEGER.
	// )`)
	var result = make([]CNC.CNC_DTO, 0)

	if PR.Db == nil {
		return result
	}

	query := `SELECT 
        UniqueKey,
        machineType,
        chip,
        customName,
		connectionType,
		connectionData
    FROM Machines`

	rows, err := PR.Db.Query(query)

	if err != nil {
		if err != sql.ErrNoRows {
			log.Println(err)
		}
		log.Println(err) //TODO
		return result
	}

	// var Scan struct {
	// 	UniqueKey      string
	// 	machineType    int
	// 	chip           string
	// 	customName     string
	// 	connectionType string
	// 	connectionData string
	// }

	DTO := CNC.CNC_DTO{}

	for rows.Next() {
		err := rows.Scan(
			&DTO.UniqueKey,
			&DTO.MACHINE_TYPE,
			&DTO.Device_Chip_Name,
			&DTO.TARGET_MACHINE_NAME,
			&DTO.ConnectionType,
			&DTO.ConnectionData)

		if err != nil {
			log.Println(err)
			return result
		}

		result = append(result, DTO)
	}
	// fmt.Printf("result: %v\n", result)
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
		printerData.FIRMWARE)
	if ex != nil {
		return ex
	}
	return nil
}
