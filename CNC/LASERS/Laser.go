package laser

import "CNCManager/CNC"

type Laser struct {
	CNC.CNCCore
}

func InitStandartLaser() {
	Mfunck := func() CNC.AnyCNC {
		return &Laser{}
	}
	CNC.RegisterCNC("Esp_32_Laser", Mfunck)
}
