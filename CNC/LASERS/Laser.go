package laser

import (
	"CNCManager/CNC"
	"CNCManager/CNC/CNCService"
	"context"
)

type Laser struct {
	CNC.CNCCore
}

func InitStandartLaser() {
	Mfunck := func() CNC.RealizeCNC {
		return &Laser{}
	}
	CNC.RegisterCNC("Esp_32_Laser", Mfunck)
}

func (L *Laser) SetCore(core *CNC.CNCCore) {
	L.CNCCore = *core
}

func (L *Laser) InitRealization() error {
	// FDM.Fans = make(map[int]uint8)
	return nil
}

func (L *Laser) GetJsonData() any {

	return ""
}

func (L *Laser) ExecuteTask(file []byte, ctx context.Context) {
	L.WriteLog(CNCService.LogLevelInformation, "start printing!")
	for _, Data := range L.WorkFile {
		select {
		case <-ctx.Done():
			return
		default:
			if Data == "" {
				continue
			}
			res := CNCService.DeleteComments_GCode(Data) //:= (FDMPrinter.Prepare_Command_to_printer(Data)) //todo
			if res == "" || res == CNCService.EndOfData {
				continue
			}
			L.SendMessage([]byte(res))
		}
	}
}

func (L *Laser) ParseCommand(Prefix, dataStr string) {

}
