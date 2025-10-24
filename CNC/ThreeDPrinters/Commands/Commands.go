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

var V1Commands = map[int]string{
	EndOfData:          "\r\n",
	Error:              "E\r_",
	StopPrint:          "!\r_",
	GetTemps:           "M105",
	CheckConnection:    "%\r_",
	GetBaseInformation: "&\r_",
	Check:              "*\r_",
	NowTemperatureBed:  "B\r_",
	TemperatureNozzle:  "N\r_",
	IsPrinting:         "P\r_",
	ReadyToRead:        "R\r_",
	BufferCommandSize:  "S\r_",
	ItsGcodeCommand:    "G\r_",
	ClearBuffer:        "C\r_",
	SetLightStatus:     "L\r_",
	ItsTemperatureN:    "N\r_",
	ItsTemperatureB:    "B\r_",
	CheckPrinter:       "*\r_",
	BufferACK:          "ok",
	ImPrinting:         "P\r_",
	MPositionX:         "X\r_",
	MPositionY:         "Y\r_",
	MPositionZ:         "Z\r_",
	MBufferCommandSize: "S\r_",
	MMaxBufferSize:     "^\r_",
	MWidth:             "W\r_",
	MLength:            "L\r_",
	MHeight:            "H\r_",
	MVersion:           "V\r_",
	MName:              "N\r_",
	MType:              "T\r_",

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
