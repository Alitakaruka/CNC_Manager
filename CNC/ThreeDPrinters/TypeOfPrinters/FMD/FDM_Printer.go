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
	data := strings.Split(string(file), "\n")
	for _, Data := range data {
		select {
		case <-ctx.Done():
			return
		default:
			res := CNCService.DeleteComments_GCode(Data)

			if strings.HasPrefix(res, "M104") || //TODO DEBUG!
				strings.HasPrefix(res, "M109") ||
				strings.HasPrefix(res, "M140") ||
				strings.HasPrefix(res, "M190") {
				continue // ← пропускаем нагрев
			}

			if res == "" || res == CNCService.EndOfData {
				continue
			}
			// fmt.Printf("res: %v\n", res)
			if ok := P.SendMessage([]byte(res + CNCService.EndOfData)); !ok {
				P.WriteLog(CNCService.LogLevelError, "Printing aborted! Command cant send!")
				return
			}
		}
	}
}

func SkipHeatingCommands(gcode string) string {
	lines := strings.Split(gcode, "\n")
	out := make([]string, 0, len(lines))

	for _, line := range lines {
		l := strings.TrimSpace(line)

		// комментарии пропускаем как есть
		if strings.HasPrefix(l, ";") || l == "" {
			out = append(out, line)
			continue
		}

		// берём только первую команду (до комментария)
		cmd := l
		if i := strings.Index(cmd, ";"); i >= 0 {
			cmd = cmd[:i]
		}
		cmd = strings.ToUpper(strings.TrimSpace(cmd))

		if strings.HasPrefix(cmd, "M104") ||
			strings.HasPrefix(cmd, "M109") ||
			strings.HasPrefix(cmd, "M140") ||
			strings.HasPrefix(cmd, "M190") {
			continue // ← пропускаем нагрев
		}

		out = append(out, line)
	}

	return strings.Join(out, "\n")
}

func (FDM *FDMPrinterData) SetCore(core *CNC.CNCCore) {
	FDM.CNCCore = *core
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
