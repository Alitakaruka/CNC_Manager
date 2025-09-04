package AtmegaPrinter

import (
	"PrinterManager/CNC/CNCService"
	FDMPrinter "PrinterManager/CNC/ThreeDPrinters/TypeOfPrinters/FMD"
)

type AtmegaPrinter struct {
	FDMPrinter.FDMPrinterData
}

func (P *AtmegaPrinter) ExecuteTask(file []byte) error {
	err := P.LoadFileForWork(file)
	if err != nil {
		return err
	}
	P.DTO.IsWorking = true
	go func() {
		for _, Data := range P.WorkFile {
			if Data == "" {
				continue
			}
			if !P.DTO.IsWorking {
				break
			}
			res := "" //:= (FDMPrinter.Prepare_Command_to_printer(Data))
			if res == "" || res == P.Protocol.Command(CNCService.EndOfData) {
				continue
			}
			P.SendMessage([]byte(res))
		}
		P.DTO.IsWorking = false
	}()
	return nil
}

func InitAtmegaPrinter() {

}
