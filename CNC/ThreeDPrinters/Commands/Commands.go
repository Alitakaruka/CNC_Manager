package Commands

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
)

// Таблица для сопоставления ID → строка
var V1Commands = map[int]string{
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
}

const PrinterTimeOut = 10
const InformationUpdateTime = 4

// #ifndef COMANDS_H

// #define COMANDS_H

// #define Debug_Flag1          UART_send_byte('1')

// #define Debug_Flag2          UART_send_byte('2')

// #define Debug_Flag3          UART_send_byte('3')

// #define Debug_Flag4          UART_send_byte('4')

// #define Debug_Flag5          UART_send_byte('5')

// #define Debug_Flag6          UART_send_byte('6')

// #define Debug_Flag7          UART_send_byte('7')

// #define Debug_Flag8          UART_send_byte('8')

// #define endComand            ";"

// #define EndComandByte        ';'

// #define UnknownValue         "__"

// //////////////////////////////////////////////////To the printer

// #define StopPrint            "!_"
// //#define StartPrint           "@_"

// #define GetAllInformation    "#_"

// #define GetBaseInformation   "&_"

// #define Check                "*_"

// #define GetTemps             "@_"

// #define NowTemperatureBed    "B_"

// #define NowTemperatureNozzle "N_"

// #define ReadyToRead          "R_"

// #define BufferComandSize     "S_"

// ///////////////////////////////////////////////////

// ///////////////////////////////////////////////////From the printer

// #define MTemperatureNozzle "N_%d"

// #define MTemperatureBed    "B_%d"

// #define M_PositionX        "X_%d"

// #define M_PositionY        "Y_%d"

// #define M_PositionZ        "Z_%d"

// #define M_BufferComandSize "S_%d"

// #define ErrorPrinter       "E_%d"

// #define M_Name             "n_%s"

// #define M_Version          "V_%s"

// #define M_Width            "W_%d"

// #define M_Length           "L_%d"

// #define M_Height           "H_%d"

// #define M_Type             "T_%s"

// ////////////////////////////////////////////////////GCode

// #define G0  "G0"

// #define G1  "G1"

// #define G2  "G2"

// #define G28 "G28"

// /////////////////////////////////////////////////////

// /////////////////////////////////////////////////////Errors code

// /////////////////////////////////////////////////////
