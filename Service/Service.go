package Service

import (
	laser "CNCManager/CNC/LASERS"
	AtmegaPrinter "CNCManager/CNC/ThreeDPrinters/TypeOfPrinters/FMD/Printers"
	"errors"
	"log"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

const version = "1.1.0"

type Config struct {
	Server struct {
		Addr    string `yaml:"addr"`
		Port    string `yaml:"port"`
		MaxLogs int    `yaml:"maxLogs"`
	} `yaml:"server"`

	Database struct {
		Path string `yaml:"path"`
	} `yaml:"database"`
	Logs struct {
		Path string `yaml:"path"`
	} `yaml:"logs"`
}

func GetConfig(path string) *Config {
	file, ex := os.ReadFile(path)

	if ex != nil {
		log.Printf("Error to load config file:%v", ex)
		return nil
	}
	config := Config{}
	if ex = yaml.Unmarshal(file, &config); ex != nil {
		log.Printf("Error unmarshal config:%v", ex)
		return nil
	}
	return &config
}

func InitPrinters() {
	AtmegaPrinter.InitAtmegaPrinter()
	laser.InitStandartLaser()

	//TODO: other printers
}

func ValidVersion(FirmwareVersion string) error {
	split := strings.Split(FirmwareVersion, ".")
	if len(split) != 3 {
		return errors.New("strange version:" + FirmwareVersion)
	}
	Pmajor, _ := strconv.Atoi(split[0])
	Pminor, _ := strconv.Atoi(split[1])

	AppArr := strings.Split(version, ".")
	if len(split) != 3 {
		return errors.New("strange version:" + FirmwareVersion)
	}
	Amajor, _ := strconv.Atoi(AppArr[0])
	Aminor, _ := strconv.Atoi(AppArr[1])

	if Pmajor > Amajor {
		return errors.New("major is must")
	}
	if Pminor > Aminor {
		return errors.New("minor is must")
	}
	return nil
}

func CRC16(data []byte) uint16 {
	var crc uint16 = 0xFFFF
	for _, b := range data {
		crc ^= uint16(b)
		for i := 0; i < 8; i++ {
			if crc&1 != 0 {
				crc = (crc >> 1) ^ 0xA001
			} else {
				crc >>= 1
			}
		}
	}
	return crc
}
