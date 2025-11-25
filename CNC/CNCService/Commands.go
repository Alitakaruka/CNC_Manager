package CNCService

const Identification = "Identification"

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
	EndOfData          = "\r\n" // ";"
	StopPrint          = "!_"   // "!_"
	GetAllInformation  = "#_"   // "#_"
	CheckConnection    = "%_"   // "%_"
	GetBaseInformation = "&_"   // "&_"
	Check              = "*_"   // "*_"
	NowTemperatureBed  = "B_"   // "B_"
	TemperatureNozzle  = "N_"   // "N_"
	IsPrinting         = "P_"   // "P_"
	ReadyToRead        = "R_"   // "R_"
	BufferCommandSize  = "S_"   // "S_"
	ItsGcodeCommand    = "G_"   // "G_"
	ClearBuffer        = "C_"   // "C_"
	SetLightStatus     = "L_"   // "L_"
)

const (
	MyTemperatureN      = "N:"          // "N_"
	MyTemperatureB      = "B:"          // "B_"
	BufferACK           = "ok"          // "ok"
	ImPrinting          = "IsPrinting:" // "P_"
	MyPositionX         = "X:"          // "X_"
	MyPositionY         = "Y:"          // "Y_"
	MyPositionZ         = "Z:"          // "Z_"
	MyBufferCommandSize = "Buf:"        // "S_"
	MyMaxBufferSize     = "^_"          // "^_"

	Error            = "Error:"    // "E_"
	MyWidth          = "M_Width:"  // "W_"
	MyLength         = "M_Length:" // "L_"
	MyHeight         = "M_Height:" // "H_"
	MyName           = "M_Name:"   // "n_"
	MyType           = "M_Type:"   // "T_"
	DEVICE_CHIP_NAME = "Device_chip_name:"

	SwitchTimeout = "Switch_Timeout:"

	ConnectionType = "ConnectionType:"
	SYNC           = "+\r_"
)

// Immutable
const (
	WIFI = "WIFI:"
)

const (
	ErrMemoryAlloc      = 32 // "0x01"
	ErrParseCommand     = 33 // "0x02"
	ErrUndefinedCommand = 34 // "0x03"
	ErrOutOfRange       = 35 // "0x04"
	ErrBufferOverflow   = 36 // "0x05"
	ErrTXBufferOverflow = 37 // "0x06"
	ErrRXBufferOverflow = 38 // "0x07"
)

const (
	StartOfTransmision = "F\\1"
	EndOfTransmision   = "F\\4"
	FILE_NAME          = "FILENAME:"
	FILE_SIZE          = "SIZE:"
	GET_FILE_FATA      = "GET_FILE_FATA:%d"
)
const PrinterTimeOut = 10
const InformationUpdateTime = 4
