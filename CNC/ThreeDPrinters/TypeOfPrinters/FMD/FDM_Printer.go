package FDM_Printer

import (
	"CNCManager/CNC"
	"strconv"
	"strings"
)

const (
	MaxTempNozzle = 300
	MinTempNozzle = -14
	MaxTempBed    = 100
	MinTempBed    = -14
)

type FDMPrinterData struct {
	CNC.CNCCore
	//////////////////////// Printer state
	ExtruderTemp  int `json:"nowTempNozzle"`
	ExtruderTemp1 int `json:"nowTempNozzle1"`
	TempBed       int `json:"nowTempBed"`
	////////////////////////
	HasLight bool
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

func (FDM *FDMPrinterData) ParseCommand(Prefix, dataStr string) {
	switch Prefix {
	case ExtruderTemp:
		FDM.ExtruderTemp, _ = strconv.Atoi(dataStr)
	case Extruder1Temp:
		FDM.ExtruderTemp1, _ = strconv.Atoi(dataStr)
	case BedTemp:
		FDM.TempBed, _ = strconv.Atoi(dataStr)
	}
}
