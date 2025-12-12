package FDM_Printer

import (
	"CNCManager/CNC"
	"CNCManager/CNC/CNCService"
	"context"
	"fmt"
	"log"
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

	Extruder1 struct {
		CurTemp  int
		NeedTemp int
	}
	Extruder2 struct {
		CurTemp  int
		NeedTemp int
	}
	Bed struct {
		CurTemp  int
		NeedTemp int
	}

	Fans map[int]uint8
}

func (FDM *FDMPrinterData) InitRealization() error {
	FDM.Fans = make(map[int]uint8)
	return nil
}

func (FDM *FDMPrinterData) GetJsonData() any {
	var jsonData struct {
		NozzleTemp string        `json:"nozzleTemp"`
		BedTemp    string        `json:"bedTemp"`
		Fans       map[int]uint8 `json:"fans"`
	}

	jsonData.NozzleTemp = strconv.Itoa(FDM.Extruder1.CurTemp) + " / " + strconv.Itoa(FDM.Extruder1.NeedTemp)
	jsonData.BedTemp = strconv.Itoa(FDM.Bed.CurTemp) + " / " + strconv.Itoa(FDM.Bed.NeedTemp)

	jsonData.Fans = make(map[int]uint8, len(FDM.Fans))
	for k, v := range FDM.Fans {
		jsonData.Fans[k] = v
	}
	return jsonData
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

func (FDM *FDMPrinterData) CheckTemps() {
	// FDM.
}

func (P *FDMPrinterData) ExecuteTask(file []byte, ctx context.Context) {
	P.WriteLog(CNCService.LogLevelInformation, "start printing!")
	for _, Data := range P.WorkFile {
		select {
		case <-ctx.Done():
			return
		default:
			res := CNCService.DeleteComments_GCode(Data)
			if res == "" || res == CNCService.EndOfData {
				continue
			}
			P.SendMessage([]byte(res + CNCService.EndOfData))
		}
	}
}

func (FDM *FDMPrinterData) ParseCommand(Prefix, dataStr string) {

	switch Prefix {
	case ExtruderTempPref:
		strs := strings.Split(dataStr, "/")
		FDM.Extruder1.CurTemp, _ = strconv.Atoi(strs[0])
		FDM.Extruder1.NeedTemp, _ = strconv.Atoi(strs[1])
		// _, err := fmt.Sscanf(dataStr, BedTemp, &FDM.Extruder1.CurTemp, &FDM.Extruder1.NeedTemp)
		// if err != nil {
		// 	log.Println(err)
		// 	FDM.WriteLog(CNCService.LogLevelError, err.Error())
		// }
	case BedTempPref:
		strs := strings.Split(dataStr, "/")
		FDM.Bed.CurTemp, _ = strconv.Atoi(strs[0])
		FDM.Bed.NeedTemp, _ = strconv.Atoi(strs[1])

		// _, err := fmt.Sscanf(dataStr, BedTemp, &FDM.Bed.CurTemp, &FDM.Bed.NeedTemp)
		// fmt.Printf("FDM.Bed.CurTemp: %v\n", FDM.Bed.CurTemp)
		// if err != nil {
		// 	log.Println(err)
		// 	FDM.WriteLog(CNCService.LogLevelError, err.Error())
		// }
	case FanSpeedPref:

		var index int
		var value uint8
		_, err := fmt.Sscanf(Prefix+dataStr, FanSpeed, &index, &value)
		if err != nil {
			log.Println(err)
			FDM.WriteLog(CNCService.LogLevelError, err.Error())
		}
		if FDM.Fans != nil {
			FDM.Fans[index] = uint8(value)
		}
	}
}
