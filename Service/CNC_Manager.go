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
)

type ConnectionData struct {
	TypeOfConnection string
	ConnectionData   string
}

type CNCManagerr struct {
	CNC_Machines []*CNC.CNCCore
	DataBase     DataBase.PrinterRepository
}

func (CNC_M *CNCManagerr) InitManager(sqlPath string) {
	CNC_M.DataBase.InitRepository(sqlPath)
	machines := CNC_M.DataBase.GetAllMachines()
	for _, machine := range machines {
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
		// CNC_M.CNC_Machines = append(CNC_M.CNC_Machines, &machine)
	}
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

		err := newCNC.InitDevice()

		if err != nil {
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
		// CNC_M.DataBase.AddMachine(newCNC.GetCore())

		newCNC.CNCStart()
		return nil
	}
}

// func copyCommonFields(src, dst any) {
// 	s := reflect.ValueOf(src).Elem()
// 	d := reflect.ValueOf(dst).Elem()

// 	for i := 0; i < s.NumField(); i++ {
// 		f := s.Type().Field(i)
// 		sv := s.Field(i)

// 		dv := d.FieldByName(f.Name)
// 		if dv.IsValid() && dv.Type() == sv.Type() && dv.CanSet() {
// 			dv.Set(sv)
// 		}
// 	}
// }

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

	} else {
		return errors.New("cnc not found")
	}
	return nil
}

func (CNC_M *CNCManagerr) reconect(index int) error {
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

func (CNC_M *CNCManagerr) Reconnect(uniqueKey string) error {
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

func (CNC_M *CNCManagerr) GetJson() string {
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
		fmt.Printf("val: %v\n", val)
		if ind, find := CNC_M.findByKey(Key); find {
			go CNC_M.CNC_Machines[ind].SendMessage([]byte(val + CNCService.EndOfData))
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

func (CNC_M *CNCManagerr) UploadFile(key, filename string, file []byte) error {
	index, ok := CNC_M.findByKey(key)
	if !ok {
		return errors.New("the device doesnt find")
	}
	CNC_M.CNC_Machines[index].UploadFile(filename, file)

	return nil
}
