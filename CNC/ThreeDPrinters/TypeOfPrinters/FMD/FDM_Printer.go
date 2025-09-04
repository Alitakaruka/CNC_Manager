package FDM_Printer

import (
	"PrinterManager/CNC"
	"PrinterManager/CNC/CNCService"
	PrinterService "PrinterManager/CNC/CNCService"
	"log"
	"strings"
	"time"
)

const (
	MaxTempNozzle = 300
	MinTempNozzle = -14
	MaxTempBed    = 100
	MinTempBed    = -14
)

type FDMPrinterData struct {
	CNC.CNCCore
	////////////////////////
	Width  int `json:"width"`
	Length int `json:"length"`
	Height int `json:"height"`
	//////////////////////// Printer state
	NowTempNozzle int `json:"nowTempNozzle"`
	NowTempBed    int `json:"nowTempBed"`

	//////////////////////// Printer position
	MyXposition float32 `json:"myXposition"`
	MyYposition float32 `json:"myYposition"`
	MyZposition float32 `json:"myZposition"`

	////////////////////////
	HasLight bool
	////////////////////////
	WathcDog *time.Timer `json:"-"`
	//////////////////////// Information
	Loger  *PrinterService.Loger
	Errors []error `json:"-"`
}

func Delete_comments(command string) string {
	strArr := strings.Split(command, ";")
	if len(strArr) == 0 {
		command = strings.TrimSpace(command)
	}
	if len(strArr) >= 2 {
		return strings.TrimSpace(strArr[0])
	}
	return command
}

func (FDM *FDMPrinterData) CommandsInData() {
	if len(FDM.ReceiveBuffer) == 0 {
		return
	}
	FDM.Mutex.Lock()
	bufferCopy := append([]byte(nil), FDM.ReceiveBuffer...)
	FDM.ReceiveBuffer = FDM.ReceiveBuffer[:0] //clear
	FDM.Mutex.Unlock()
	Commands := strings.Split(string(bufferCopy), FDM.Protocol.Command(CNCService.EndOfData))
	for _, value := range Commands {
		if value == "" {
			continue
		}
		r := []rune(value)
		Prefix := string(r[:2])
		Command, _ := strings.CutPrefix(value, Prefix)
		FDM.ParseCommand(Prefix, Command)
	}

}

func (FDM *FDMPrinterData) ParseCommand(Prefix, Command string) {
	switch Prefix {

	case FDM.Protocol.Command(CNCService.BufferACK):
		FDM.Transmitter.Increment()
	case FDM.Protocol.Command(CNCService.Check):
		log.Println("Check")
		return
	case FDM.Protocol.Command(CNCService.ItsTemperatureN):
		PrinterService.SetIntValue(&FDM.NowTempNozzle, Command, &FDM.Mutex)
	case FDM.Protocol.Command(CNCService.ItsTemperatureB):
		PrinterService.SetIntValue(&FDM.NowTempBed, Command, &FDM.Mutex)
	case FDM.Protocol.Command(CNCService.MPositionX):
		PrinterService.SetFloatValue(&FDM.MyXposition, Command, &FDM.Mutex)
	case FDM.Protocol.Command(CNCService.MPositionY):
		PrinterService.SetFloatValue(&FDM.MyYposition, Command, &FDM.Mutex)
	case FDM.Protocol.Command(CNCService.MPositionZ):
		PrinterService.SetFloatValue(&FDM.MyZposition, Command, &FDM.Mutex)
	case FDM.Protocol.Command(CNCService.Error):
		// PrinterData, _ := strings.CutPrefix(Command, FDM.Protocol.Command(CNCService.Error))
		// FDM.Log_printer_error(PrinterData)
	case FDM.Protocol.Command(CNCService.MLength):
		PrinterService.SetIntValue(&FDM.Length, Command, &FDM.Mutex)
	case FDM.Protocol.Command(CNCService.MHeight):
		PrinterService.SetIntValue(&FDM.Height, Command, &FDM.Mutex)
	case FDM.Protocol.Command(CNCService.MWidth):
		PrinterService.SetIntValue(&FDM.Width, Command, &FDM.Mutex)
	default:
		log.Printf("Undefined command:%v ,Len: %v", Command, len(FDM.ReceiveBuffer))
	}
}
