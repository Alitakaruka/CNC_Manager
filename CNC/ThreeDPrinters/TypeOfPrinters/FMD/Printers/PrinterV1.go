package AtmegaPrinter

import (
	"CNCManager/CNC"
	FDMPrinter "CNCManager/CNC/ThreeDPrinters/TypeOfPrinters/FMD"
)

type AtmegaPrinter struct {
	FDMPrinter.FDMPrinterData
}

func InitAtmegaPrinter() {
	Mfunck := func() CNC.RealizeCNC {
		return &AtmegaPrinter{}
	}
	CNC.RegisterCNC("ATM_16", Mfunck)
	CNC.RegisterCNC("ATM_32", Mfunck)
}
