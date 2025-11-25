package Service

import (
	"CNCManager/CNC"
	CNCService "CNCManager/CNC/CNCService"
	FDMPrinter "CNCManager/CNC/ThreeDPrinters/TypeOfPrinters/FMD"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"
)

type ConnectionData struct {
	TypeOfConnection string
	ConnectionData   string
}

type CNCManagerr struct {
	CNC_Machines []CNC.AnyCNC
}

func (CNC_M *CNCManagerr) Connect(conData ConnectionData) error {
	if index, find := CNC_M.findByConnectionData(conData); find {
		if CNC_M.IsConnected(index) {
			return errors.New("CNC is already connected")
		} else {
			return CNC_M.reconect(index)
		}
	} else {
		newCNC, ex := CNC.Connect(conData.TypeOfConnection, conData.ConnectionData)
		if ex != nil {
			return ex
		}

		var targer CNC.AnyCNC
		var ok bool
		err := newCNC.InitDevice()
		if err != nil {
			newCNC.CloseConnection()
			return err
		}

		log.Println("CHIP:" + newCNC.GetDTO().Device_Chip_Name)
		if targer, ok = CNC.Machines[newCNC.GetDTO().Device_Chip_Name]; !ok {
			newCNC.CloseConnection()
			return errors.New("the device has not register")
		}

		// Get DTO and set unique key if not set
		dto := newCNC.GetDTO()
		if dto.UniqueKey == "" {
			dto.UniqueKey = CNC_M.GenerateUniqueKey()
		}
		targer.SetDTO(dto)

		CNC_M.CNC_Machines = append(CNC_M.CNC_Machines, targer)

		targer.CNCStart()
		return nil
	}
}

func (CNC_M *CNCManagerr) IsConnected(index int) bool {
	DTO := CNC_M.CNC_Machines[index].GetDTO()
	return DTO.Flags.Connected
}

func (CNC_M *CNCManagerr) findByConnectionData(ConData ConnectionData) (int, bool) {
	ConnectionString := ConData.TypeOfConnection + ":" + ConData.ConnectionData
	for ind, CNC := range CNC_M.CNC_Machines {
		DTO := CNC.GetDTO()
		if DTO.ConnectionString == ConnectionString {
			return ind, true
		}
	}
	return 0, false
}

func (CNC_M *CNCManagerr) findByKey(key string) (int, bool) {
	for ind, CNC := range CNC_M.CNC_Machines {
		DTO := CNC.GetDTO()
		if DTO.UniqueKey == key {
			return ind, true
		}
	}
	return 0, false
}

func (CNC_M *CNCManagerr) ExecuteTask(key string, byteFile []byte) error {
	if index, find := CNC_M.findByKey(key); find {
		cnc := CNC_M.CNC_Machines[index]
		err := cnc.StartTask(byteFile)
		if err != nil {
			return err
		}
		cntx, cancel := context.WithCancel(context.Background())
		go func() {
			dto := cnc.GetDTO()
			dto.Flags.ExecutingTask = true
			cnc.SetDTO(dto)

			fmt.Printf("cnc: %v\n", cnc.GetDTO())
			cnc.ExecuteTask(byteFile, cntx)
			dto.Flags.ExecutingTask = false
			cnc.SetDTO(dto)
		}()
		go func() {
			for cnc.CanExecuteTask() {
				time.Sleep(time.Millisecond * 100)
			}
			cancel()
		}()
	} else {
		return errors.New("cnc not found")
	}
	return nil
}

func (CNC_M *CNCManagerr) reconect(index int) error {
	CNC := CNC_M.CNC_Machines[index]
	_, err := CNC.Reconnect()
	if err != nil {
		return err
	}
	CNC.InitDevice()
	CNC.CNCStart()
	return nil
}

type CNC_JSON struct {
	CNC_Name         string `json:"CNCName"`
	Version          string `json:"version"`
	CncType          string `json:"CncType"`
	UniqueKey        string `json:"uniqueKey"`
	IsWorking        bool   `json:"isWorking"`
	ExecutingTask    bool   `json:"executingTask"`
	NozzleTemp       int    `json:"nozzleTemp"`
	BedTemp          int    `json:"bedTemp"`
	TypeOfConnection string `json:"typeOfConnection"`
	Progress         int    `json:"progress"`
	TimeRemaining    int    `json:"timeRemaining"`
	FileStorage      bool   `json:"fileStorage"`
	StorageFilesFMT  string `json:"storageFilesFMT"`

	Position struct {
		X float32 `json:"X"`
		Y float32 `json:"Y"`
		Z float32 `json:"Z"`
	} `json:"position"`
	Immutable struct {
		Width  int `json:"width"`
		Length int `json:"length"`
		Height int `json:"height"`
	} `json:"immutable"`
}

func (CNC_M *CNCManagerr) GetJson() string {
	CNCs := make([]CNC_JSON, 0, len(CNC_M.CNC_Machines))

	for _, machine := range CNC_M.CNC_Machines {
		dto := machine.GetDTO()
		// log.Println(dto)
		CNC := CNC_JSON{
			CNC_Name:         dto.TARGET_MACHINE_NAME,
			Version:          dto.FIRMWARE_VERSION,
			CncType:          getMachineTypeName(dto.MACHINE_TYPE),
			UniqueKey:        dto.UniqueKey,
			IsWorking:        dto.Flags.Connected,
			ExecutingTask:    dto.Flags.ExecutingTask,
			TypeOfConnection: dto.ConnectionData,
			FileStorage:      dto.Memory.FileStorage,
			StorageFilesFMT:  dto.Memory.FileStorageFMT,
			Progress:         0,
			TimeRemaining:    0,
		}

		CNC.Position.X = dto.Position.X
		CNC.Position.Y = dto.Position.Y
		CNC.Position.Z = dto.Position.Z

		CNC.Immutable.Width = dto.Immutable.Width
		CNC.Immutable.Length = dto.Immutable.Length
		CNC.Immutable.Height = dto.Immutable.Height

		CNC.NozzleTemp = 0
		CNC.BedTemp = 0

		// Try type assertion for FDMPrinterData directly
		if fdmMachine, ok := machine.(*FDMPrinter.FDMPrinterData); ok {
			CNC.NozzleTemp = fdmMachine.ExtruderTemp
			CNC.BedTemp = fdmMachine.TempBed
		} else {
			// Try to access via reflection for embedded structs
			// This will work for AtmegaPrinter which embeds FDMPrinterData
			tempData := getFDMData(machine)
			if tempData != nil {
				CNC.NozzleTemp = tempData.ExtruderTemp
				CNC.BedTemp = tempData.TempBed
			}
		}

		CNCs = append(CNCs, CNC)
	}

	jsonData, err := json.Marshal(CNCs)
	if err != nil {
		log.Printf("Error marshaling CNCs to JSON: %v", err)
		return "[]"
	}
	return string(jsonData)
}

func getMachineTypeName(machineType int) string {
	if name, ok := CNCService.MachinesTypes[machineType]; ok {
		return name
	}
	return "Unknown"
}

// Helper function to extract FDM data using reflection
// This handles both direct FDMPrinterData and embedded (like AtmegaPrinter)
func getFDMData(machine CNC.AnyCNC) *FDMPrinter.FDMPrinterData {
	// Direct type assertion
	if fdm, ok := machine.(*FDMPrinter.FDMPrinterData); ok {
		return fdm
	}

	// Use reflection to access embedded FDMPrinterData
	v := reflect.ValueOf(machine)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// Try to find FDMPrinterData field (embedded)
	fdmField := v.FieldByName("FDMPrinterData")
	if fdmField.IsValid() {
		// Get pointer to the field to avoid copying mutex
		if fdmField.CanAddr() {
			fdmPtr := fdmField.Addr().Interface()
			if fdm, ok := fdmPtr.(*FDMPrinter.FDMPrinterData); ok {
				return fdm
			}
		}
	}

	return nil
}

func (CNC_M *CNCManagerr) GetAllLogs() []CNCService.Log {
	result := make([]CNCService.Log, 0)
	for _, CNC := range CNC_M.CNC_Machines {
		Logs := CNC.GetLogs()
		result = append(result, Logs...)
	}
	return result
}

func (CNC_M *CNCManagerr) SendGCode(GCode, Key string) error {
	Commands := strings.Split(GCode, "\n")
	for _, val := range Commands {
		if ind, find := CNC_M.findByKey(Key); find {
			go CNC_M.CNC_Machines[ind].SendMessage([]byte(val))
			return nil
		}
	}
	return errors.New("CNC not found")
}

func (CNC_M *CNCManagerr) SaveSettings() {

}

func (CNC_M *CNCManagerr) GenerateUniqueKey() string {
	for {
		key := GenerateRandomKey()
		if CNC_M.isUnique(key) {
			return key
		}
	}
}

func (CNC_M *CNCManagerr) isUnique(key string) bool {
	for _, CNC := range CNC_M.CNC_Machines {
		DTO := CNC.GetDTO()
		if DTO.UniqueKey == key {
			return false
		}
	}
	return true
}

func GenerateRandomKey() string {
	bytes := make([]byte, 8)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}

func (CNC_M *CNCManagerr) UploadFile(key, filename string, file []byte) {
	index, ok := CNC_M.findByKey(key)

	if !ok {
		return
	}
	CNC_M.CNC_Machines[index].UploadFile(filename, file)
}
