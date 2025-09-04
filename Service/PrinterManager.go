package Service

import (
	"PrinterManager/CNC"
	"PrinterManager/CNC/CNCService/Connectors"
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
	connectors   []Connectors.CNCConnector
	CNC_Machines []CNC.AnyCNC
	Printers     []ThreeDPrinter.AnyPrinter
}

func (PM *CNCManagerr) Connect(conData ConnectionData) error {
	if index, find := PM.findByConnectionData(conData); find {
		if PM.IsConnected(index) {
			return errors.New("CNC is already connected")
		} else {
			return PM.reconect(index)
		}
	} else {
		newPrinter, ex := CNC.Connect(conData.TypeOfConnection, conData.ConnectionData)
		if ex != nil {
			return ex
		}
		// key := PM.GenerateUniqueKey()
		// DTO := newPrinter.GetDTO()
		PM.CNC_Machines = append(PM.CNC_Machines, newPrinter)
		// PM.Printers = append(PM.Printers, newPrinter)
		// newPrinter.StartWork()
		newPrinter.InitDevice()
		newPrinter.CNCStart()
		return nil
	}
}

func (PM *CNCManagerr) IsConnected(index int) bool {
	DTO := PM.Printers[index].GetDTO()
	return DTO.IsWorking
}

func (PM *CNCManagerr) findByConnectionData(ConData ConnectionData) (int, bool) {
	for ind, printer := range PM.Printers {
		DTO := printer.GetDTO()
		if (DTO.ConnectionData == ConData.ConnectionData) &&
			(DTO.TypeOfConnection == ConData.TypeOfConnection) {
			return ind, true
		}
	}
	return 0, false
}

func (PM *CNCManagerr) findByKey(key string) (int, bool) {
	for ind, printer := range PM.Printers {
		DTO := printer.GetDTO()
		if DTO.UniqueKey == key {
			return ind, true
		}
	}
	return 0, false
}

func (PM *CNCManagerr) StartPrint(key string, byteFile []byte) error {
	if index, find := PM.findByKey(key); find {
		ex := PM.Printers[index].StartPrint(byteFile)
		return ex
	}
	return errors.New("printer not found")
}

func (PM *CNCManagerr) reconect(index int) error {
	DTO := PM.Printers[index].GetDTO()
	newPrinter, ex := ThreeDPrinter.ConnectPrinter(DTO.TypeOfConnection, DTO.ConnectionData)
	if ex != nil {
		return ex
	}
	newPrinter.GetDTO().UniqueKey = DTO.UniqueKey
	PM.Printers[index] = newPrinter
	newPrinter.StartWork()
	return nil
}

func (PM *CNCManagerr) GetJson() string {
	result := "["

	for _, printer := range PM.Printers {
		result += (printer.Get_JsonData() + ",")
	}
	result, _ = strings.CutSuffix(result, ",")
	result += "]"
	return result
}

func (PM *CNCManagerr) LoggingAsync() {
	for {
		for _, printer := range PM.Printers {
			Logs := printer.Get_Logs()
			for _, iLog := range Logs {
				log.Println(iLog)
			}
		}
		time.Sleep(time.Second * 5)
	}
}

func (PM *CNCManagerr) SendGCode(GCode, Key string) error {
	Commands := strings.Split(GCode, "\n")
	for _, val := range Commands {
		log.Printf("Gcode sended:%v", val)
		if ind, find := PM.findByKey(Key); find {
			go PM.Printers[ind].SendMessege([]byte(GCode))
			return nil
		}
	}
	return errors.New("printer not found")
}

func (PM *CNCManagerr) SaveSettings() {

}

func (PM *CNCManagerr) GenerateUniqueKey() string {
	for {
		key := GenerateRandomKey()
		if PM.isUnique(key) {
			return key
		}
	}
}

func (PM *CNCManagerr) isUnique(key string) bool {
	for _, printer := range PM.Printers {
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

func (PM *CNCManagerr) SetColor(R, G, B byte, key string) error {
	if index, ok := PM.findByKey(key); ok {
		go PM.Printers[index].SetColorLight(R, G, B)
		return nil
	} else {
		return errors.New("Printer not found")
	}

}
