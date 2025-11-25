package AtmegaPrinter

import (
	"CNCManager/CNC"
	"CNCManager/CNC/CNCService"
	FDMPrinter "CNCManager/CNC/ThreeDPrinters/TypeOfPrinters/FMD"
	"context"
)

type AtmegaPrinter struct {
	FDMPrinter.FDMPrinterData
}

func (P *AtmegaPrinter) ExecuteTask(file []byte, ctx context.Context) {
	P.WriteLog("start printing!", CNCService.LogLevelInformation)
	for _, Data := range P.WorkFile {
		select {
		case <-ctx.Done():
			P.WriteLog("print stoped!", CNCService.LogLevelError)
			return
		default:
			if Data == "" {
				continue
			}
			res := "" //:= (FDMPrinter.Prepare_Command_to_printer(Data)) //todo
			if res == "" || res == CNCService.EndOfData {
				continue
			}
			P.SendMessage([]byte(res))
		}
	}
}

func InitAtmegaPrinter() {
	Mfunck := func() CNC.AnyCNC {
		return &AtmegaPrinter{}
	}
	CNC.RegisterCNC("ATM_16", Mfunck)
	CNC.RegisterCNC("ATM_32", Mfunck)
}
