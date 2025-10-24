package CNCService

const Identification = "GET_FIRMWARE_DATA"

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

const (
	MyTemperatureN      = 16 // "N_"
	MyTemperatureB      = 17 // "B_"
	BufferACK           = 19 // "ok"
	ImPrinting          = 20 // "P_"
	MyPositionX         = 21 // "X_"
	MyPositionY         = 22 // "Y_"
	MyPositionZ         = 23 // "Z_"
	MyBufferCommandSize = 24 // "S_"
	MyMaxBufferSize     = 25 // "^_"
	MyWidth             = 26 // "W_"
	MyLength            = 27 // "L_"
	MyHeight            = 28 // "H_"
	MyName              = 30 // "n_"
	MyType              = 31 // "T_"
)

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

func GetCommand(command int) string {
	return Commands[command]
}

var Commands = map[int]string{
	EndOfData:          "\r\n",
	Error:              "E\r_",
	StopPrint:          "!\r_",
	GetTemps:           "M105",
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

	BufferACK: "ok",
	SYNC:      "+\r_",
}

var CNC_Data = map[int]string{
	MyPositionX:         "X\r_",
	MyPositionY:         "Y\r_",
	MyPositionZ:         "Z\r_",
	MyBufferCommandSize: "S\r_",
	MyMaxBufferSize:     "^\r_",
	MyWidth:             "W\r_",
	MyLength:            "L\r_",
	MyHeight:            "H\r_",
	MyName:              "n\r_",
	MyType:              "T\r_",
	MyTemperatureN:      "N\r_",
	MyTemperatureB:      "B\r_",
	ImPrinting:          "P\r_",
}

const PrinterTimeOut = 10
const InformationUpdateTime = 4
