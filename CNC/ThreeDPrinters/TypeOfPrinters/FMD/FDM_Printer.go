package FDM_Printer

import (
	PrinterService "PrinterManager/CNC/CNCService"
	ThreeDPrinter "PrinterManager/CNC/ThreeDPrinters"
	Commands "PrinterManager/CNC/ThreeDPrinters/Commands"
	"bufio"
	"errors"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	MaxTempNozzle = 300
	MinTempNozzle = -14
	MaxTempBed    = 100
	MinTempBed    = -14
)

type FDMPrinterData struct {
	ThreeDPrinter.PrinterDTO
	////////////////////////
	Width  int `json:"width"`
	Length int `json:"length"`
	Height int `json:"height"`
	//////////////////////// Printer state
	NowTempNozzle int  `json:"nowTempNozzle"`
	NowTempBed    int  `json:"nowTempBed"`
	IsPrinting    bool `json:"isPrinting"`

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
	Errors []error    `json:"-"`
	Buffer []string   `json:"-"`
	Mut    sync.Mutex `json:"-"`
}

func (FDM *FDMPrinterData) SetColorLight(R, G, B byte) {
	FDM.SendMessege([]byte(Commands.SetLightStatus + "R" + strconv.Itoa(int(R)) +
		"G" + strconv.Itoa(int(G)) +
		"B" + strconv.Itoa(int(B))))
}

func (FDM *FDMPrinterData) InitAndStart() {
	FDM.IsWorking = true
	FDM.Loger = PrinterService.NewLoger(100)
}

func Prepare_Command_to_printer(command string) string {
	command = Delete_comments(command)
	command = strings.TrimSpace(command)
	return command + Commands.EndOfData
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
	if len(FDM.Buffer) == 0 {
		return
	}
	FDM.Mut.Lock()
	bufferCopy := append([]string(nil), FDM.Buffer...)
	FDM.Buffer = FDM.Buffer[:0] //clear
	FDM.Mut.Unlock()
	for _, value := range bufferCopy {
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

	case Commands.Buffer_ACK:
		log.Println("ACK")
		FDM.PrinterBuffer.Increment()
	case Commands.CheckPrinter:
		log.Println("Check")
		return
	case Commands.ItsTemperatureN:
		PrinterService.SetIntValue(&FDM.NowTempNozzle, Command, &FDM.Mut)
	case Commands.ItsTemperatureB:
		PrinterService.SetIntValue(&FDM.NowTempBed, Command, &FDM.Mut)
	case Commands.M_PositionX:
		PrinterService.SetFloatValue(&FDM.MyXposition, Command, &FDM.Mut)
	case Commands.M_PositionY:
		PrinterService.SetFloatValue(&FDM.MyYposition, Command, &FDM.Mut)
	case Commands.M_PositionZ:
		PrinterService.SetFloatValue(&FDM.MyZposition, Command, &FDM.Mut)
	case Commands.M_BufferCommandSize:
		BufferSize, exe := strconv.Atoi(Command)
		if exe != nil {
			FDM.Errors = append(FDM.Errors, exe)
		} else {
			FDM.Mut.Lock()
			FDM.PrinterBuffer.SetBufferSize(uint(BufferSize))
			FDM.Mut.Unlock()
		}
	case Commands.Error:
		PrinterData, _ := strings.CutPrefix(Command, Commands.Error)
		FDM.Log_printer_error(PrinterData)
	case Commands.M_Version:
		PrinterService.SetStringValue(&FDM.Version, Command, &FDM.Mut)
	case Commands.M_Name:
		PrinterService.SetStringValue(&FDM.PrinterName, Command, &FDM.Mut)
	case Commands.M_Type:
		PrinterService.SetStringValue(&FDM.PrinterType, Command, &FDM.Mut)
	case Commands.M_Length:
		PrinterService.SetIntValue(&FDM.Length, Command, &FDM.Mut)
	case Commands.M_Height:
		PrinterService.SetIntValue(&FDM.Height, Command, &FDM.Mut)
	case Commands.M_Width:
		PrinterService.SetIntValue(&FDM.Width, Command, &FDM.Mut)
	default:
		log.Printf("Undefined command:%v ,Len: %v", Command, len(FDM.Buffer))
	}
}

func (FDM *FDMPrinterData) SetFieldValue(Field, value *any) {
	FDM.Mut.Lock()
	*Field = *value
	FDM.Mut.Unlock()
}

func (FDM *FDMPrinterData) SendMessege(message []byte) error {
	log.Printf("BufferSize:%v", FDM.PrinterBuffer.GetValueData())
	FDM.PrinterBuffer.WaitForNonZero()
	FDM.PrinterBuffer.Decrement()

	if FDM.Connection == nil {
		return errors.New("printer does not connected")
	}
	_, ex := FDM.Connection.Write(message)
	log.Printf("i send:%v\n", string(message))
	if ex != nil {
		log.Printf("Error sending:%v", ex.Error())
	}
	return ex
}

func (FDM *FDMPrinterData) ReadConnectionAsync() {
	clear(FDM.Buffer)
	reader := bufio.NewReader(FDM.Connection)
	for FDM.IsWorking {
		line, ex := reader.ReadString(Commands.EndCommandByte)
		if ex != nil {
			FDM.Errors = append(FDM.Errors, ex)
		} else {
			FDM.Mut.Lock()
			FDM.WathcDog.Reset(time.Second * Commands.PrinterTimeOut)
			if line == "" {
				continue
			}
			res, _ := strings.CutSuffix(line, Commands.EndOfData)
			FDM.Buffer = append(FDM.Buffer, res)
			FDM.Mut.Unlock()
		}
	}
}

func (FDM *FDMPrinterData) StartWathcDog() {
	FDM.WathcDog = time.NewTimer(time.Second * Commands.PrinterTimeOut)

	<-FDM.WathcDog.C
	log.Println("Printer timeot!")
	FDM.IsWorking = false
	FDM.Connection.Close()
}

func (FDM *FDMPrinterData) Log_printer_error(errorCode string) {
	// Data.mutex.TryLock()
	FDM.Mut.Lock()
	switch errorCode {
	case Commands.MemoryAllocError:
		FDM.Errors = append(FDM.Errors, errors.New("memory alloc is faled"))
	case Commands.ParseCommandError:
		FDM.Errors = append(FDM.Errors, errors.New("the printer could not parse the command"))
	case Commands.UndefinedCommand:
		FDM.Errors = append(FDM.Errors, errors.New("the printer did not understand the command"))
	case Commands.OutOfRange:
		FDM.Errors = append(FDM.Errors, errors.New("the printer has exceeded its print limits"))
	default:
		FDM.Errors = append(FDM.Errors, errors.New("indefined error:"+errorCode))
	}
	FDM.Mut.Unlock()
	// Data.mutex.Unlock()
}
