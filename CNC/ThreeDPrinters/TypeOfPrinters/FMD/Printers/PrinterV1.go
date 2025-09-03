package AtmegaPrinter

import (
	"encoding/json"
	"errors"
	"log"
	"strings"
	"time"

	ThreeDPrinter "PrinterManager/CNC/ThreeDPrinters"
	Commands "PrinterManager/CNC/ThreeDPrinters/Commands"
	FDMPrinter "PrinterManager/CNC/ThreeDPrinters/TypeOfPrinters/FMD"
)

type AtmegaPrinter struct {
	FDMPrinter.FDMPrinterData
	CommandsFile []string `json:"-"`
}

func (P *AtmegaPrinter) StartWork() {
	P.InitAndStart()

	go P.ReadConnectionAsync()
	go P.StartWathcDog()
	go P.updatePrinterInformationAsync()
	go P.readBufferAsync()
	P.getAllData()
}

func (P *AtmegaPrinter) CheckTemp() {

	tiker := time.NewTicker(time.Second * 5)
	for P.IsWorking {
		<-tiker.C
		if P.NowTempBed > FDMPrinter.MaxTempBed {
			P.Panic("maximum bed temperature exceeded")
		}
		if P.NowTempBed <= FDMPrinter.MinTempBed {
			P.Panic("minimum bed temperature exceeded")
		}
		if P.NowTempNozzle > FDMPrinter.MaxTempNozzle {
			P.Panic("maximum nozzle temperature exceeded")
		}
		if P.NowTempNozzle <= FDMPrinter.MinTempNozzle {
			P.Panic("minimum nozzle temperature exceeded")
		}
	}
}

func (P *AtmegaPrinter) print() {
	for _, Data := range P.CommandsFile {
		if Data == "" {
			continue
		}
		if !P.IsWorking {
			break
		}
		res := (FDMPrinter.Prepare_Command_to_printer(Data))
		if res == "" || res == Commands.EndOfData {
			continue
		}
		P.SendMessege([]byte(res))
	}
	P.IsPrinting = false
}

func (P *AtmegaPrinter) Panic(err string) {
	P.IsPrinting = false
	P.SendMessege([]byte("TODO:Panic"))
	P.Errors = append(P.Errors, errors.New("Panic error:"+err))
	log.Printf("Panic error:%v \n", err)
}

// func (P *AtmegaPrinter) Connect() error {
// 	coinfig := Serial.OpenOptions{PortName: P.MainData.ConnectionData,
// 		BaudRate:              9600,
// 		DataBits:              8,
// 		StopBits:              1,
// 		ParityMode:            Serial.PARITY_NONE,
// 		RTSCTSFlowControl:     false,
// 		InterCharacterTimeout: 100,
// 	}
// 	port, ex := Serial.Open(coinfig)
// 	if ex != nil {
// 		return ex
// 	}
// 	P.Port = port
// 	P.StartWork()
// 	return nil
// }

// func (P *AtmegaPrinter) FillData(data ThreeDPrinter.PrinterDTO, ConnectionData io.ReadWriteCloser) {
// 	P.Connection = ConnectionData
// 	P.MainData = data
// }

func InitAtmegaPrinter() {
	Mfunck := func() ThreeDPrinter.AnyPrinter {
		return &AtmegaPrinter{}
	}
	ThreeDPrinter.RegisterPrinters("ATM16", Mfunck)
	ThreeDPrinter.RegisterPrinters("ATM32", Mfunck)
}

func (P *AtmegaPrinter) Get_JsonData() string {
	jsonStr, ex := json.Marshal(&P)
	if ex != nil {
		return ex.Error()
	}
	return string(jsonStr)
}

func (P *AtmegaPrinter) StartPrint(BinaryData []byte) error {
	clear(P.CommandsFile)
	DataFile := string(BinaryData)
	if P.Connection == nil {
		return errors.New("port is not connected")
	}

	if P.IsPrinting {
		return errors.New("printer is already print")
	}

	P.CommandsFile = strings.Split(DataFile, "\n")
	go P.print()
	P.IsPrinting = true
	return nil
}

func (P *AtmegaPrinter) updatePrinterInformationAsync() {
	ticker := time.NewTicker(time.Second * Commands.InformationUpdateTime)
	for P.IsWorking {
		<-ticker.C
		P.SendMessege([]byte(Commands.GetAllInformation + Commands.EndOfData))
	}
}

func (P *AtmegaPrinter) readBufferAsync() {
	for P.IsWorking {
		P.CommandsInData()
	}
}

func (P *AtmegaPrinter) getAllData() {
	P.SendMessege([]byte(Commands.GetAllInformation + Commands.EndOfData))
}

func (P *AtmegaPrinter) Get_Logs() []string {
	allLogs := make([]string, 0, len(P.Errors))
	P.Mut.Lock()
	defer P.Mut.Unlock()
	for _, LOG := range P.Errors {
		if LOG == nil {
			return allLogs
		}
		allLogs = append(allLogs, LOG.Error())
	}
	P.Errors = P.Errors[:0] //clear logs
	return allLogs
}
