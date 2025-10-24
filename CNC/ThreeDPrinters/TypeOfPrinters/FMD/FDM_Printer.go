package FDM_Printer

import (
	"CNCManager/CNC"
	"CNCManager/CNC/CNCService"
	PrinterService "CNCManager/CNC/CNCService"
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
	Commands := strings.Split(string(bufferCopy), CNCService.Commands[CNCService.EndOfData])
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

	case CNCService.Commands[CNCService.BufferACK]:
		FDM.Transmitter.Increment()
	case CNCService.Commands[CNCService.Check]:
		log.Println("Check")
		return
	case CNCService.Commands[CNCService.MyTemperatureN]:
		PrinterService.SetIntValue(&FDM.NowTempNozzle, Command, &FDM.Mutex)
	case CNCService.Commands[CNCService.MyTemperatureB]:
		PrinterService.SetIntValue(&FDM.NowTempBed, Command, &FDM.Mutex)
	case CNCService.Commands[CNCService.MyPositionX]:
		PrinterService.SetFloatValue(&FDM.MyXposition, Command, &FDM.Mutex)
	case CNCService.Commands[CNCService.MyPositionY]:
		PrinterService.SetFloatValue(&FDM.MyYposition, Command, &FDM.Mutex)
	case CNCService.Commands[CNCService.MyPositionZ]:
		PrinterService.SetFloatValue(&FDM.MyZposition, Command, &FDM.Mutex)
	case CNCService.Commands[CNCService.Error]:
		// PrinterData, _ := strings.CutPrefix(Command, CNCService.Commands[CNCService.Error))
		// FDM.Log_printer_error(PrinterData)
	case CNCService.Commands[CNCService.MyLength]:
		PrinterService.SetIntValue(&FDM.Length, Command, &FDM.Mutex)
	case CNCService.Commands[CNCService.MyHeight]:
		PrinterService.SetIntValue(&FDM.Height, Command, &FDM.Mutex)
	case CNCService.Commands[CNCService.MyWidth]:
		PrinterService.SetIntValue(&FDM.Width, Command, &FDM.Mutex)
	default:
		log.Printf("Undefined command:%v ,Len: %v", Command, len(FDM.ReceiveBuffer))
	}
}
