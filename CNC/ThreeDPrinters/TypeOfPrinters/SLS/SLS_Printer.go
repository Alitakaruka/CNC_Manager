package SLS_Printer

import (
	"CNCManager/CNC"
	"CNCManager/CNC/CNCService"
	"archive/zip"
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/png"
	"log"
)

type SLSPrinterData struct {
	Core   *CNC.CNCCore
	Images *CNCService.Transmitter
	Matrix struct {
		X    int
		Y    int
		inch float32
	}
}

func (SLS *SLSPrinterData) SetCore(core *CNC.CNCCore) {
	SLS.Core = core
}

func (SLS *SLSPrinterData) InitRealization() error {
	SLS.Images = CNCService.NewTransmitter()
	SLS.Images.SetLimits(1, 1)

	SLS.Core.SendMessage([]byte(SyncImageBuffer + CNCService.EndOfData))
	SLS.Core.WriteLog(CNCService.LogLevelInformation, "Init!")
	// SLS.Fans = make(map[int]uint8)
	return nil
}

func (SLS *SLSPrinterData) GetJsonData() any {
	// var jsonData struct {
	// 	NozzleTemp string        `json:"nozzleTemp"`
	// 	BedTemp    string        `json:"bedTemp"`
	// 	Fans       map[int]uint8 `json:"fans"`
	// }

	// SLS.mut.Lock()
	// jsonData.NozzleTemp = strconv.Itoa(SLS.Extruder1.CurTemp) + " / " + strconv.Itoa(SLS.Extruder1.NeedTemp)
	// jsonData.BedTemp = strconv.Itoa(SLS.Bed.CurTemp) + " / " + strconv.Itoa(SLS.Bed.NeedTemp)

	// jsonData.Fans = make(map[int]uint8, len(SLS.Fans))
	// for k, v := range SLS.Fans {
	// 	jsonData.Fans[k] = v
	// }
	// SLS.mut.Unlock()
	// return jsonData
	return nil
}

func (P *SLSPrinterData) ExecuteTask(file []byte) {

	reader, err := zip.NewReader(bytes.NewReader(file), int64(len(file)))

	if err != nil {
		log.Println(err)
		P.Core.WriteLog(CNCService.LogLevelError, err.Error())
	}

	P.Core.WriteLog(CNCService.LogLevelInformation, "start printing!")

	// fmt.Printf("file: %v\n", file)
	for _, f := range reader.File {

		fmt.Println("Файл:", f.Name)
		if f.Name == "config.json" {
			fmt.Println("TODO CONFIG!")
			continue
		}

		file, err := f.Open()
		if err != nil {
			P.Core.WriteLog(CNCService.LogLevelError, err.Error())
			continue
		}
		image, err := png.Decode(file)

		if err != nil {
			P.Core.WriteLog(CNCService.LogLevelError, err.Error())
			continue
		}

		bitMap := ConvertImage2BitMap(image)
		ImageResult := base64.StdEncoding.EncodeToString(bitMap)
		P.Core.UploadFileInMemory([]byte(ImageResult))
	}

	fmt.Println("End!")
	P.Core.IsTaskEnd <- struct{}{}
}

func (SLS *SLSPrinterData) SendImage(base64Str []byte) {

}

func (SLS *SLSPrinterData) CreatePages(pg chan string, base64Str []byte) {
	// MaxBytes := SLS.Core.Transmitter.MaxBytes - len(ImageData)
	// MaxBytes -= MaxBytes % 4 //base64  encoded data = 4

	// log.Println(len(base64Str) / MaxBytes)
	// Page := ImageData
	// for _, val := range base64Str {
	// 	Page += string(val)
	// 	if len(Page) == MaxBytes {
	// 		pg <- Page
	// 		Page = ImageData
	// 	}
	// }
	// close(pg)
}

func ConvertImage2BitMap(img image.Image) []byte {
	X := img.Bounds().Dx()
	Y := img.Bounds().Dy()

	BitVector := CNCService.NewBitSet(X * Y)

	i := 0
	for y := 0; y < Y; y++ {
		for x := 0; x < X; x++ {
			color := img.At(x, y)
			if r, g, b, _ := color.RGBA(); r >= 127 && g >= 127 && b >= 127 {
				BitVector.Set(i)
			} else {
				BitVector.Clear(i)
			}
			i++
		}
	}
	return BitVector.GetData()
}

func InitSLSPrinter() {
	Mfunck := func() CNC.RealizeCNC {
		return &SLSPrinterData{}
	}
	CNC.RegisterCNC("ESP32_S3_SLS", Mfunck)
	CNC.RegisterCNC("ESP32_SLS", Mfunck)
}

func (SLS *SLSPrinterData) ParseCommand(Prefix, dataStr string) {

	switch Prefix {
	case MatrixParams:
		var x, y int
		var inch float32
		if _, err := fmt.Sscanf(dataStr, MatrixParams, &x, &y, &inch); err != nil {
			SLS.Matrix.X = x
			SLS.Matrix.Y = y
			SLS.Matrix.inch = inch
		}

	case MaxImagePrefix:
		var max int
		if _, err := fmt.Sscanf(dataStr, MaxImages, &max); err != nil {
			SLS.Images.SetLimits(max, max)
		}
	}

}
