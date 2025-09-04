package Service

import (
	"PrinterManager/CNC"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
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
		newPrinter, ex := CNC.Connect(conData.TypeOfConnection, conData.ConnectionData)
		if ex != nil {
			return ex
		}
		CNC_M.CNC_Machines = append(CNC_M.CNC_Machines, newPrinter)
		newPrinter.InitDevice()
		newPrinter.CNCStart()
		return nil
	}
}

func (CNC_M *CNCManagerr) IsConnected(index int) bool {
	DTO := CNC_M.CNC_Machines[index].GetDTO()
	return DTO.IsWorking
}

func (CNC_M *CNCManagerr) findByConnectionData(ConData ConnectionData) (int, bool) {
	for ind, printer := range CNC_M.CNC_Machines {
		DTO := printer.GetDTO()
		if DTO.ConnectionData == ConData.ConnectionData {
			return ind, true
		}
	}
	return 0, false
}

func (CNC_M *CNCManagerr) findByKey(key string) (int, bool) {
	for ind, printer := range CNC_M.CNC_Machines {
		DTO := printer.GetDTO()
		if DTO.UniqueKey == key {
			return ind, true
		}
	}
	return 0, false
}

func (CNC_M *CNCManagerr) ExecuteTask(key string, byteFile []byte) error {
	if index, find := CNC_M.findByKey(key); find {
		ex := CNC_M.CNC_Machines[index].ExecuteTask(byteFile)
		return ex
	}
	return errors.New("CNC not found")
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

func (CNC_M *CNCManagerr) GetJson() string {
	return ""
	// result := "["

	// for _, printer := range CNC_M.CNC_Machines {
	// 	result += (printer.GetLogs() + ",")
	// }
	// result, _ = strings.CutSuffix(result, ",")
	// result += "]"
	// return result
}

func (CNC_M *CNCManagerr) LoggingAsync() {
	for {
		for _, printer := range CNC_M.CNC_Machines {
			Logs := printer.GetLogs()
			for _, iLog := range Logs {
				log.Println(iLog)
			}
		}
		time.Sleep(time.Second * 5)
	}
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
	for _, printer := range CNC_M.CNC_Machines {
		DTO := printer.GetDTO()
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
