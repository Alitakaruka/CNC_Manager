package Service

import (
	"CNCManager/CNC"
	CNCService "CNCManager/CNC/CNCService"
	"CNCManager/Service/DataBase"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
)

type ConnectionData struct {
	TypeOfConnection string
	ConnectionData   string
}

type CNCManager struct {
	IsChargeTable    chan []byte
	TimeIsCharged    chan []byte
	NewDataCncBuffer chan []byte
	IsStatusCharge   chan []byte
	Logs             chan CNCService.Log

	online  int
	ofline  int
	working int
	total   int

	CNC_Machines []*CNC.CNCCore
	DataBase     DataBase.CNCRepository
}

func (CNC_M *CNCManager) InitManager(sqlPath string) {
	CNC_M.IsChargeTable = make(chan []byte)
	CNC_M.TimeIsCharged = make(chan []byte)
	CNC_M.NewDataCncBuffer = make(chan []byte, 100)
	CNC_M.IsStatusCharge = make(chan []byte)
	CNC_M.Logs = make(chan CNCService.Log, 128)

	//TimerUpgrade
	go func() {
		seconds := 0

		Trig := make(chan struct{}, 1)
		ticker := time.NewTicker(time.Second * 60)
		defer ticker.Stop()

		go func() {
			for range ticker.C {
				Trig <- struct{}{}
			}
		}()
		for range Trig {
			seconds += 60
			curTime := seconds
			days := curTime / 86400
			curTime %= 86400
			hours := curTime / 3600
			curTime %= 3600
			minutes := curTime / 60
			curTime %= 60
			CNC_M.TimeIsCharged <- []byte(fmt.Sprintf("%dd %dh %dm", days, hours, minutes))
		}
	}()

	CNC_M.DataBase.InitRepository(sqlPath)
	machines := CNC_M.DataBase.GetAllMachines()
	for _, machine := range machines {
		machine := machine
		machine.Connection = CNC.GetConnector(machine.DTO.ConnectionData, machine.DTO.ConnectionString)
		if machine.Connection != nil {
			go func() {
				err := machine.Reconnect()
				if err != nil {
					log.Println(err)
					return
				}
				err = machine.InitDevice()
				if err != nil {
					log.Println(err)
					return
				}
				machine.CNCStart()
			}()

		} else {
			machine.WriteLog(CNCService.LogLevelWarning, "Failed to connect to the machine")
		}
		CNC_M.CNC_Machines = append(CNC_M.CNC_Machines, machine)
	}
}

func (CNC_M *CNCManager) Connect(conData ConnectionData) error {
	if index, find := CNC_M.findByConnectionData(conData); find {
		if CNC_M.IsConnected(index) {
			return errors.New("CNC is already connected")
		} else {
			return CNC_M.reconect(index)
		}
	}
	newCNC, ex := CNC.Connect(conData.TypeOfConnection, conData.ConnectionData)
	// log.Printf("Connected by type %v, data: %v\n", conData.TypeOfConnection, conData.ConnectionData)
	if ex != nil {
		return ex
	}
	err := newCNC.InitDevice()

	if err != nil {
		// log.Printf("Device init error:%v", err)
		newCNC.CloseConnection()
		return err
	}
	// Get DTO and set unique key if not set
	dto := newCNC.GetDTO()
	if dto.UniqueKey == "" {
		dto.UniqueKey = CNC_M.GenerateUniqueKey()
	}
	newCNC.SetDTO(dto)
	CNC_M.CNC_Machines = append(CNC_M.CNC_Machines, newCNC)

	go CNC_M.UpdateMachineData(newCNC)

	CNC_M.IsChargeTable <- []byte(CNC_M.GetJson())
	go newCNC.CNCStart()

	log.Println("New cnc start success!")
	newCNC.WriteLog(CNCService.LogLevelSuccess, "Successfully connected")

	return nil
}

func (CNC_M *CNCManager) UpdateLogs() {

}

func (CNC_M *CNCManager) GetTimeCharge() []byte {
	return <-CNC_M.TimeIsCharged
}

func (CNC_M *CNCManager) UpdateMachineData(Machine *CNC.CNCCore) {
	CNC_M.online++
	CNC_M.chargeStatus()
	for {
		select {
		case <-Machine.IsCharge:
			json := CNC_M.GetMachineJson(Machine)
			CNC_M.NewDataCncBuffer <- json

		case log := <-Machine.Logs:
			CNC_M.Logs <- log
		case <-Machine.IsClose:
			CNC_M.online--
			CNC_M.chargeStatus()
			CNC_M.NewDataCncBuffer <- CNC_M.GetMachineJson(Machine)
			return
			// default:

		}
	}
}

func (CNC_M *CNCManager) chargeStatus() {
	CNC_M.IsStatusCharge <- []byte(fmt.Sprintf(`{"type":"status","data":{"online":%d,"printing":%d,"offline":%d,"total":%d}}`, CNC_M.online, CNC_M.working, CNC_M.ofline, CNC_M.total))
}

func (CNC_M *CNCManager) GetStatus() []byte {
	return <-CNC_M.IsStatusCharge
}

func (CNC_M *CNCManager) GetTable() []byte {
	return <-CNC_M.IsChargeTable
}

func (M *CNCManager) GetNewCncData() []byte {
	return <-M.NewDataCncBuffer
}

func (CNC_M *CNCManager) IsConnected(index int) bool {
	DTO := CNC_M.CNC_Machines[index].GetDTO()
	return DTO.Flags.Connected
}

func (CNC_M *CNCManager) findByConnectionData(ConData ConnectionData) (int, bool) {
	ConnectionString := ConData.TypeOfConnection + ":" + ConData.ConnectionData
	for ind, CNC := range CNC_M.CNC_Machines {
		DTO := CNC.GetDTO()
		if DTO.ConnectionString == ConnectionString {
			return ind, true
		}
	}
	return 0, false
}

func (CNC_M *CNCManager) findByKey(key string) (int, bool) {
	for ind, CNC := range CNC_M.CNC_Machines {
		DTO := CNC.GetDTO()
		if DTO.UniqueKey == key {
			return ind, true
		}
	}
	return 0, false
}

func (CNC_M *CNCManager) ExecuteTask(key string, byteFile []byte) error {
	if index, find := CNC_M.findByKey(key); find {
		cnc := CNC_M.CNC_Machines[index]
		err := cnc.StartTask(byteFile)
		if err != nil {
			return err
		}

	} else {
		return errors.New("cnc not found")
	}
	return nil
}

func (CNC_M *CNCManager) reconect(index int) error {
	CNC := CNC_M.CNC_Machines[index]
	err := CNC.Reconnect()
	if err != nil {
		return err
	}
	err = CNC.InitDevice()
	if err != nil {
		return err
	}
	CNC.CNCStart()

	return nil
}

func (CNC_M *CNCManager) Reconnect(uniqueKey string) error {
	if ind, ok := CNC_M.findByKey(uniqueKey); ok {
		return CNC_M.reconect(ind)
	}
	return errors.New("the device not found")
}

type CNC_JSON struct {
	CNC_Name string `json:"CNCName"`
	// Version          string `json:"version"`
	CncType          string `json:"CncType"`
	UniqueKey        string `json:"uniqueKey"`
	IsWorking        bool   `json:"isWorking"`
	ExecutingTask    bool   `json:"executingTask"`
	TypeOfConnection string `json:"typeOfConnection"`
	Progress         int    `json:"progress"`
	TimeRemaining    int    `json:"timeRemaining"`
	FileStorage      bool   `json:"fileStorage"`
	StorageFilesFMT  string `json:"storageFilesFMT"`

	TDP any `json:"TDP"` //Specifical device data

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

	Light struct {
		HasLight bool `json:"hasLight"`
		RGBLight bool `json:"rgbLight"`
	} `json:"light"`
}

func (CNC_M *CNCManager) GetMachineJson(machine *CNC.CNCCore) []byte {
	// js := CNC_JSON{}
	dto := machine.GetDTO()
	CNC := CNC_JSON{
		CNC_Name: dto.TARGET_MACHINE_NAME,
		// Version:          dto.FIRMWARE_VERSION,
		CncType:          getMachineTypeName(dto.MACHINE_TYPE),
		UniqueKey:        dto.UniqueKey,
		IsWorking:        dto.Flags.Connected,
		ExecutingTask:    dto.Flags.ExecutingTask,
		TypeOfConnection: dto.ConnectionData,
		FileStorage:      dto.Memory.FileStorage,
		StorageFilesFMT:  dto.Memory.FileStorageFMT,
		Progress:         machine.Progress,
		TimeRemaining:    0,
	}

	CNC.Position.X = dto.Position.X
	CNC.Position.Y = dto.Position.Y
	CNC.Position.Z = dto.Position.Z

	CNC.Immutable.Width = dto.Immutable.Width
	CNC.Immutable.Length = dto.Immutable.Length
	CNC.Immutable.Height = dto.Immutable.Height

	// CNC.TDP.NozzleTemp = "0 / 0"
	// CNC.TDP.BedTemp = "0 / 0"
	if realize := machine.Realize; realize != nil {
		CNC.TDP = machine.Realize.GetJsonData()
	}

	CNC.Light.HasLight = dto.Switchable.Light
	CNC.Light.RGBLight = dto.Switchable.RGB

	jsonData, err := json.Marshal(CNC)
	if err != nil {
		log.Printf("Error marshaling CNCs to JSON: %v", err)
		return []byte("[]")
	}
	return (jsonData)

}

func (CNC_M *CNCManager) GetJson() string {
	CNCs := make([]CNC_JSON, 0, len(CNC_M.CNC_Machines))

	for _, machine := range CNC_M.CNC_Machines {
		dto := machine.GetDTO()
		CNC := CNC_JSON{
			CNC_Name: dto.TARGET_MACHINE_NAME,
			// Version:          dto.FIRMWARE_VERSION,
			CncType:          getMachineTypeName(dto.MACHINE_TYPE),
			UniqueKey:        dto.UniqueKey,
			IsWorking:        dto.Flags.Connected,
			ExecutingTask:    dto.Flags.ExecutingTask,
			TypeOfConnection: dto.ConnectionData,
			FileStorage:      dto.Memory.FileStorage,
			StorageFilesFMT:  dto.Memory.FileStorageFMT,
			Progress:         machine.Progress,
			TimeRemaining:    0,
		}

		CNC.Position.X = dto.Position.X
		CNC.Position.Y = dto.Position.Y
		CNC.Position.Z = dto.Position.Z

		CNC.Immutable.Width = dto.Immutable.Width
		CNC.Immutable.Length = dto.Immutable.Length
		CNC.Immutable.Height = dto.Immutable.Height

		// CNC.TDP.NozzleTemp = "0 / 0"
		// CNC.TDP.BedTemp = "0 / 0"
		if realize := machine.Realize; realize != nil {
			CNC.TDP = machine.Realize.GetJsonData()
		}

		CNC.Light.HasLight = dto.Switchable.Light
		CNC.Light.RGBLight = dto.Switchable.RGB

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

func (CNC_M *CNCManager) GetLog() CNCService.Log {
	return <-CNC_M.Logs
}

func (CNC_M *CNCManager) SendGCode(GCode, Key string) error {
	ind, find := CNC_M.findByKey(Key)
	if !find {
		return errors.New("CNC not found")
	}

	cnc := CNC_M.CNC_Machines[ind]
	commands := strings.Split(GCode, "\n")
	go func() {
		for _, val := range commands {
			val = strings.TrimSpace(val)
			if val == "" {
				continue
			}
			cnc.SendMessage([]byte(val + CNCService.EndOfData))
		}
	}()

	return nil
}

func (CNC_M *CNCManager) GenerateUniqueKey() string {
	for {
		key := GenerateRandomKey()
		if CNC_M.isUnique(key) {
			return key
		}
	}
}

func (CNC_M *CNCManager) isUnique(key string) bool {
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

func (CNC_M *CNCManager) UploadFile(key, filename string, file []byte) error {
	index, ok := CNC_M.findByKey(key)
	if !ok {
		return errors.New("the device doesnt find")
	}
	CNC_M.CNC_Machines[index].UploadFile(filename, file)

	return nil
}
