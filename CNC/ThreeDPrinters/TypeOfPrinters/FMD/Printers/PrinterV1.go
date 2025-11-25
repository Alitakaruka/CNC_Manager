package AtmegaPrinter

import (
	"CNCManager/CNC"
	"CNCManager/CNC/CNCService"
	FDMPrinter "CNCManager/CNC/ThreeDPrinters/TypeOfPrinters/FMD"
)

type AtmegaPrinter struct {
	FDMPrinter.FDMPrinterData
}

func (P *AtmegaPrinter) ExecuteTask(file []byte) error {
	err := P.LoadFileForWork(file)
	if err != nil {
		return err
	}
	P.DTO.Flags.Connected = true
	go func() {
		for _, Data := range P.WorkFile {
			if Data == "" {
				continue
			}
			if !P.DTO.Flags.Connected {
				break
			}
			res := "" //:= (FDMPrinter.Prepare_Command_to_printer(Data))
			if res == "" || res == CNCService.EndOfData {
				continue
			}
			P.SendMessage([]byte(res))
		}
		P.DTO.Flags.Connected = false
	}()
	return nil
}

func InitAtmegaPrinter() {
	Mfunck := func() CNC.AnyCNC {
		return &AtmegaPrinter{}
	}
	CNC.RegisterCNC("ATM16", Mfunck)
	CNC.RegisterCNC("ATM32", Mfunck)
}
