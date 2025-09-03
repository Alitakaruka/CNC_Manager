package ThreeDPrinter

import (
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"time"

	"PrinterManager/CNC/CNCService"
	Commands "PrinterManager/CNC/ThreeDPrinters/Commands"

	Serial "github.com/jacobsa/go-serial/serial"
)

// //////////////////Connection type
const (
	connect_com = "COM"
	connect_ip  = "IP"
	connect_usb = "USB"
)

// ////////////////// Type of rprinter

type AnyPrinter interface {
	StartPrint(BinaryData []byte) error
	GetDTO() *PrinterDTO
	SetDTO(dto *PrinterDTO)
	Get_JsonData() string
	StartWork()
	Get_Logs() []string
	SetColorLight(R, G, B byte)
	SendMessege(message []byte) error
}

var RegisteredPrinters = map[string]func() AnyPrinter{}

func RegisterPrinters(name string, constructor func() AnyPrinter) {
	RegisteredPrinters[name] = constructor
}

type PrinterDTO struct {
	PrinterName      string                    `json:"printerName"`
	UserPrinterName  string                    `json:"userPrinterName"`
	PrinterType      string                    `json:"printerType"`
	Version          string                    `json:"version"`
	IsWorking        bool                      `json:"isWorking"`
	TypeOfConnection string                    `json:"typeOfConnection"`
	ConnectionData   string                    `json:"-"`
	UniqueKey        string                    `json:"uniqueKey"`
	PrinterBuffer    *CNCService.PrinterBuffer `json:"-"`
	Connection       io.ReadWriteCloser
}

func NewPrinterDTO() *PrinterDTO {
	return &PrinterDTO{PrinterBuffer: CNCService.NewPrinterBuffer()}
}
func (DTO *PrinterDTO) GetDTO() *PrinterDTO {
	return DTO
}

func (DTO *PrinterDTO) SetDTO(dto *PrinterDTO) {
	*DTO = *dto
}

func ConnectPrinter(typeOfConnection string, connectionData string) (AnyPrinter, error) {
	if typeOfConnection == connect_com {
		return connectCom(connectionData)
	} else if typeOfConnection == connect_ip {
		return connectIP(connectionData)
	}

	return nil, errors.New("")
}

func connectCom(connectionData string) (AnyPrinter, error) {
	options := Serial.OpenOptions{
		PortName:              connectionData,
		BaudRate:              9600,
		DataBits:              8,
		StopBits:              1,
		MinimumReadSize:       1,
		ParityMode:            Serial.PARITY_NONE,
		RTSCTSFlowControl:     false,
		InterCharacterTimeout: 100,
	}

	port, ex := Serial.Open(options)
	if ex != nil {
		return nil, ex
	}
	basePrinterData, ex := readStartData(port)
	if ex != nil {
		return nil, ex
	}
	fmt.Println(basePrinterData)
	basePrinterData.TypeOfConnection = connect_com
	basePrinterData.ConnectionData = connectionData

	if basePrinterData.PrinterName == "" ||
		basePrinterData.Version == "" ||
		basePrinterData.PrinterType == "" {
		port.Close()
		return nil, errors.New("printer initialization error")
	}
	if constructor, ok := RegisteredPrinters[basePrinterData.PrinterName]; ok {
		printer := constructor()
		log.Printf("We set DTO:%v", printer.GetDTO())
		printer.SetDTO(basePrinterData)
		log.Printf(`New printer connected. 
		Name:%v 
		Version:%v 
		Type:%v
		DData buffer size:%v`,
			printer.GetDTO().PrinterName,
			printer.GetDTO().Version,
			printer.GetDTO().PrinterType,
			printer.GetDTO().PrinterBuffer.GetBufferSize())

		return printer, nil

	} else {
		port.Close()
		return nil, errors.New("printer does not registered")
	}
}

func readStartData(Connection io.ReadWriteCloser) (*PrinterDTO, error) {

	basePrinterData := NewPrinterDTO()

	basePrinterData.Connection = Connection

	reader := CNCService.NewTimeoutReader(Connection, time.Second)
	_, ex := Connection.Write([]byte(Commands.GetBaseInformation + Commands.EndOfData))
	result := reader.Read()
	log.Printf("I read:%v\n", result)

	commands := strings.Split(result, Commands.EndOfData)
	if ex != nil {
		return nil, ex
	}
	for _, val := range commands {
		if val == "" {
			continue
		}
		if strings.HasPrefix(val, Commands.M_Version) {
			basePrinterData.Version, _ = strings.CutPrefix(val, Commands.M_Version)
		} else if strings.HasPrefix(val, Commands.M_Name) {
			basePrinterData.PrinterName, _ = strings.CutPrefix(val, Commands.M_Name)
		} else if strings.HasPrefix(val, Commands.M_Type) {
			basePrinterData.PrinterType, _ = strings.CutPrefix(val, Commands.M_Type)
		} else if strings.HasPrefix(val, Commands.BufferCommandSize) {
			str, _ := strings.CutPrefix(val, Commands.BufferCommandSize)
			BufferSize, err := strconv.Atoi(str)
			if err != nil {
				return nil, err
			}
			basePrinterData.PrinterBuffer.SetBufferSize(uint(BufferSize))
		} else if strings.HasPrefix(val, Commands.M_MaxBufferSize) {
			str, _ := strings.CutPrefix(val, Commands.M_MaxBufferSize)
			MaxSize, err := strconv.Atoi(str)
			if err != nil {
				return nil, err
			}
			basePrinterData.PrinterBuffer.SetMaxBufferSize(uint(MaxSize))
		}
	}

	return basePrinterData, nil
}

func connectIP(connectionData string) (AnyPrinter, error) {
	return nil, errors.New("") //TODO
}
