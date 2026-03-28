package CNCService

import "strings"

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
	FDM_PRINTER: "FDM 3D PRINTER",
	LASER:       "LASER",
	SLA_PRINTER: "SLA 3D PRINTER",
	SLS_PRINTER: "SLS 3D PRINTER",
	MILLING:     "MILLING",
}

const (
	EndOfData = "\r\n" // ";"
//	StopPrint          = "!_"   // "!_"
//	GetAllInformation  = "#_"   // "#_"
//	CheckConnection    = "%_"   // "%_"
//	GetBaseInformation = "&_"   // "&_"
//	Check              = "*_"   // "*_"
//	NowTemperatureBed  = "B_"   // "B_"
//	TemperatureNozzle  = "N_"   // "N_"
//	IsPrinting         = "P_"   // "P_"
//	ReadyToRead        = "R_"   // "R_"
//	BufferCommandSize  = "S_"   // "S_"
//	ItsGcodeCommand    = "G_"   // "G_"
//	ClearBuffer        = "C_"   // "C_"
//	SetLightStatus     = "L_"   // "L_"
//
// )
)
const (
	MyTemperatureN = "N:"          // "N_"
	MyTemperatureB = "B:"          // "B_"
	BufferACK      = "ok"          // "ok"
	ImPrinting     = "IsPrinting:" // "P_"
	MyPositionX    = "X:"          // "X_"
	MyPositionY    = "Y:"          // "Y_"
	MyPositionZ    = "Z:"          // "Z_"
	MyBufferLen    = "M_Buff_Len:"

	Error            = "Error:"
	Warning          = "Warning:"
	Information      = "Info:"
	Success          = "Success:"
	MyWidth          = "M_Width:"  // "W_"
	MyLength         = "M_Length:" // "L_"
	MyHeight         = "M_Height:" // "H_"
	MyName           = "M_Name:"   // "n_"
	MyType           = "M_Type:"   // "T_"
	DEVICE_CHIP_NAME = "Device_chip_name:"

	SwitchTimeout  = "Switch_Timeout:"
	SwitchHasLight = "HasLight:"
	SwitchRGBLight = "RGBLight:"

	// ConnectionType = "ConnectionType:"
	SYNC = "SYNC"
)

// Immutable
const (
	WIFI = "WIFI:"
)

const (
	StartOfTransmision = "F\\1"
	EndOfTransmision   = "F\\4"
	FILE_NAME          = "FILENAME:"
	FILE_SIZE          = "SIZE:"
	GET_FILE_FATA      = "GET_FILE_FATA:"
	SET_LIGHT          = "SET_LIGHT:"
	RED                = "RED:"
	GREEN              = "GREEN:"
	BLUE               = "BLUE:"
)
const PrinterTimeOut = 10
const InformationUpdateTime = 4

func DeleteGCodeComments(command string) string {
	before, _, _ := strings.Cut(command, ";")
	return before
}

// DeleteComments_GCode is kept for backward compatibility.
func DeleteComments_GCode(command string) string {
	return DeleteGCodeComments(command)
}
