package CNCService

import "strings"

const Identification = "GET_FIRMWARE_DATA;"

// MACHINE TYPES
const (
	FDM_PRINTER = 1
	LASER       = 2
	SLA_PRINTER = 3
	SLS_PRINTER = 4
	MILLING     = 5
)

var MachinesTypes = map[int]string{
	FDM_PRINTER: "FDM",
	LASER:       "LASER",
	SLA_PRINTER: "SLA",
	SLS_PRINTER: "SLS",
	MILLING:     "MILLING",
}

const (
	FIRMWARE_NAME             = 1
	MACHINE_TYPE              = 2
	MACHINE_NAME              = 3
	FIRMWARE_VERSION          = 4
	EXCHANGE_PROTOCOL_VERSION = 5
)

var MachinesParametrs = map[int]string{
	FIRMWARE_NAME:             "FIRMWARE_NAME:",
	MACHINE_TYPE:              "MACHINE_TYPE:",
	MACHINE_NAME:              "TARGET_MACHINE_NAME:",
	FIRMWARE_VERSION:          "FIRMWARE_VERSION:",
	EXCHANGE_PROTOCOL_VERSION: "EXCHANGE_PROTOCOL_VERSION:",
}

//Identification example

//F_N_

const (
	EndOfData          = 0  // ";"
	Error              = 1  // "E_"
	StopPrint          = 2  // "!_"
	GetTemps           = 3  // "@_"
	GetAllInformation  = 4  // "#_"
	CheckConnection    = 5  // "%_"
	GetBaseInformation = 6  // "&_"
	Check              = 7  // "*_"
	NowTemperatureBed  = 8  // "B_"
	TemperatureNozzle  = 9  // "N_"
	IsPrinting         = 10 // "P_"
	ReadyToRead        = 11 // "R_"
	BufferCommandSize  = 12 // "S_"
	ItsGcodeCommand    = 13 // "G_"
	ClearBuffer        = 14 // "C_"
	SetLightStatus     = 15 // "L_"
)

// Команды от принтера к клиенту
const (
	ItsTemperatureN    = 16 // "N_"
	ItsTemperatureB    = 17 // "B_"
	CheckPrinter       = 18 // "*_"
	BufferACK          = 19 // "ok"
	ImPrinting         = 20 // "P_"
	MPositionX         = 21 // "X_"
	MPositionY         = 22 // "Y_"
	MPositionZ         = 23 // "Z_"
	MBufferCommandSize = 24 // "S_"
	MMaxBufferSize     = 25 // "^_"
	MWidth             = 26 // "W_"
	MLength            = 27 // "L_"
	MHeight            = 28 // "H_"
	MVersion           = 29 // "V_"
	MName              = 30 // "n_"
	MType              = 31 // "T_"
)

// Ошибки
const (
	ErrMemoryAlloc      = 32 // "0x01"
	ErrParseCommand     = 33 // "0x02"
	ErrUndefinedCommand = 34 // "0x03"
	ErrOutOfRange       = 35 // "0x04"
	ErrBufferOverflow   = 36 // "0x05"
	ErrTXBufferOverflow = 37 // "0x06"
	ErrRXBufferOverflow = 38 // "0x07"
	SYNC                = 40
)

// Exchange protocols
const (
	AliPri_GCode_V1 = iota
	AliPri_Images
)

// Таблица для сопоставления ID → строка
type ExchangeProtocol struct {
	Protocol int
}

func (EP *ExchangeProtocol) BuildTransmitData(Commands ...string) string {
	strResult := ""
	table := Protocols[EP.Protocol]
	sample := table[Tamlate]
	for _, Command := range Commands {
		for _, template := range Tamlates {
			strResult += strings.Replace(sample, template, Command, 1)
		}
	}
	return strResult
}

func (EP *ExchangeProtocol) BuildTransmitDataInt(Commands ...int) string {
	strResult := ""
	table := Protocols[EP.Protocol]
	sample := table[Tamlate]
	for _, Command := range Commands {
		for _, template := range Tamlates {
			strResult += strings.Replace(sample, template, table[Command], 1)
		}
	}
	return strResult
}

func (EP *ExchangeProtocol) Command(comm int) string {
	return Protocols[EP.Protocol][comm]
}

const Tamlate = 6666

var Tamlates = []string{"[COMMAND]", "[TYPE]"}

var Protocols = []map[int]string{
	0: map[int]string{
		Tamlate:            "[COMMAND];",
		EndOfData:          ";",
		Error:              "E_",
		StopPrint:          "!_",
		GetTemps:           "@_",
		GetAllInformation:  "#_",
		CheckConnection:    "%_",
		GetBaseInformation: "&_",
		Check:              "*_",
		NowTemperatureBed:  "B_",
		TemperatureNozzle:  "N_",
		IsPrinting:         "P_",
		ReadyToRead:        "R_",
		BufferCommandSize:  "S_",
		ItsGcodeCommand:    "G_",
		ClearBuffer:        "C_",
		SetLightStatus:     "L_",

		ItsTemperatureN:    "N_",
		ItsTemperatureB:    "B_",
		CheckPrinter:       "*_",
		BufferACK:          "ok",
		ImPrinting:         "P_",
		MPositionX:         "X_",
		MPositionY:         "Y_",
		MPositionZ:         "Z_",
		MBufferCommandSize: "S_",
		MMaxBufferSize:     "^_",
		MWidth:             "W_",
		MLength:            "L_",
		MHeight:            "H_",
		MVersion:           "V_",
		MName:              "n_",
		MType:              "T_",

		ErrMemoryAlloc:      "0x01",
		ErrParseCommand:     "0x02",
		ErrUndefinedCommand: "0x03",
		ErrOutOfRange:       "0x04",
		ErrBufferOverflow:   "0x05",
		ErrTXBufferOverflow: "0x06",
		ErrRXBufferOverflow: "0x07",

		SYNC: "+_",
	},
}

const PrinterTimeOut = 10
const InformationUpdateTime = 4
