package CNC

import (
	"CNCManager/CNC/CNCService"
	"CNCManager/CNC/CNCService/Connectors"
	"bufio"
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	LogErrorPrefix        = "E_"
	LogWarningPrefix      = "W_"
	LogInformationgPrefix = "I_"
)
const BaseTimeout = 10

type AnyCNC interface {
	GetDTO() CNC_DTO
	SetDTO(DTO CNC_DTO)
	Reconnect() (bool, error)
	GetLogs() []string
	SendMessage(message []byte)
	InitDevice() error
	CNCStart()
	ExecuteTask(file []byte) error
}

type CNCCore struct {
	DTO           *CNC_DTO
	Connection    Connectors.CNCConnector
	Transmitter   *CNCService.Transmitter
	ReceiveBuffer []byte `json:"-"`

	Mutex    sync.Mutex
	Protocol CNCService.ExchangeProtocol
	WatchDog *time.Timer
	Logs     []string
	WorkFile []string `json:"-"`
}

type CNC_DTO struct {
	TARGET_MACHINE_NAME       string `json:"TARGET_MACHINE_NAME"`
	MACHINE_TYPE              string `json:"MACHINE_TYPE"`
	FIRMWARE_VERSION          string `json:"FIRMWARE_VERSION"`
	IsWorking                 bool   `json:"isWorking"`
	Connected                 bool   `json:"Connected"`
	EXCHANGE_PROTOCOL_VERSION int    `json:"EXCHANGE_PROTOCOL_VERSION"`
	ConnectionData            string `json:"-"`
	UniqueKey                 string `json "UniqueKey"`
}

func (cnc *CNCCore) ExecuteTask(file []byte) error {
	return errors.New("this CNC can not be executing tasks")
	//stub
}

func (cnc *CNCCore) CNCStart() {
	cnc.ReceiveBuffer = make([]byte, 512)
	cnc.Transmitter = CNCService.NewTransmitter()
	cnc.Transmitter.SyncBuffers(cnc.Connection, cnc.Protocol)
	go cnc.StartWatchcDog()
	go cnc.ReadConnectionAsync()
}

func (cnc *CNCCore) StartWatchcDog() {
	cnc.WatchDog = time.NewTimer(time.Second * BaseTimeout)

	<-cnc.WatchDog.C
	log.Println(cnc.DTO.TARGET_MACHINE_NAME + " " +
		cnc.DTO.MACHINE_TYPE + " timeot!")
	cnc.writeLog(cnc.DTO.TARGET_MACHINE_NAME+" "+
		cnc.DTO.MACHINE_TYPE+" timeot!", LogErrorPrefix)
	cnc.DTO.IsWorking = false
	cnc.Connection.Close()
}

func (cnc *CNCCore) InitDevice() error {
	reader := CNCService.NewTimeoutReader(cnc.Connection, time.Second*2)
	cnc.Connection.Write([]byte(CNCService.Identification))
	res := reader.Read()
	if res == "" {
		return errors.New("the device did not respond to the request")
	}

	err := (json.Unmarshal([]byte(res), &cnc.DTO))
	if err != nil {
		return err
	}
	cnc.DTO.Connected = true
	cnc.Protocol.Protocol = cnc.DTO.EXCHANGE_PROTOCOL_VERSION
	return nil
}

func (cnc *CNCCore) FillDeviceData(str string) error {

	return nil
}
func (cnc *CNCCore) writeLog(log, logLevel string) {
	if log != "" {
		cnc.Logs = append(cnc.Logs, logLevel+log)
	}
}

func (cnc *CNCCore) ReadConnectionAsync() {
	cnc.ReceiveBuffer = cnc.ReceiveBuffer[:0]
	reader := bufio.NewReader(cnc.Connection)
	for cnc.DTO.IsWorking {
		Byte, ex := reader.ReadByte()
		if ex != nil {
			cnc.writeLog(ex.Error(), LogErrorPrefix)
		} else {
			cnc.Mutex.Lock()
			cnc.WatchDog.Reset(time.Second * BaseTimeout)
			cnc.ReceiveBuffer = append(cnc.ReceiveBuffer, Byte)
			cnc.Mutex.Unlock()
		}
	}
}

func (cnc *CNCCore) GetDTO() CNC_DTO {
	return *cnc.DTO
}

func (cnc *CNCCore) SetDTO(DTO CNC_DTO) {
	*cnc.DTO = DTO
}

func (cnc *CNCCore) GetNextByteStream(delim byte) ([]byte, bool) {
	result := []byte{}
	cnc.Mutex.Lock()
	defer cnc.Mutex.Unlock()
	for ind, val := range cnc.ReceiveBuffer {
		if val != delim {
			result = append(result, val)
		} else {
			cnc.ReceiveBuffer = cnc.ReceiveBuffer[ind:]
			return result, true
		}
	}
	return result, false
}

func (cnc *CNCCore) SendMessage(message []byte) {
	cnc.Transmitter.WaitForNonZero()
	cnc.Transmitter.Decrement()

	if cnc.Connection == nil {
		log.Println("CNC does not connected")
	}
	_, ex := cnc.Connection.Write(message)
	if ex != nil {
		cnc.writeLog(ex.Error(), LogErrorPrefix)
	}
}

func (cnc *CNCCore) GetLogs() []string {
	copy := cnc.Logs[:len(cnc.Logs)]
	cnc.Logs = cnc.Logs[:0]
	return copy
}

func (cnc *CNCCore) Reconnect() (bool, error) {
	ok, err := cnc.Connection.Reconnect()
	if err != nil {
		return ok, err
	}
	cnc.CNCStart()
	return true, nil
}

func Connect(typeOfConnection string, connectionData string) (AnyCNC, error) {
	Core := CNCCore{}
	switch typeOfConnection {
	case "COM":
		strs := strings.Split(connectionData, ":")
		var port, Baud string
		if len(strs) == 2 {
			port = strs[0]
			Baud = strs[1]
			BaudRate, err := strconv.Atoi(Baud)
			if err != nil {
				return nil, err
			}
			Core.Connection = Connectors.NewSerialConnector(port, BaudRate)
		} else if len(strs) == 1 {
			Core.Connection = Connectors.NewSerialConnector(connectionData, 9600)
		}
	case "IP":
		strs := strings.Split(connectionData, ":")
		var ip, port string
		if len(strs) == 2 {
			ip = strs[0]
			port = strs[1]
		} else {
			return nil, errors.New("invalid IP address format")
		}
		Core.Connection = Connectors.NewIpConnector(ip, port)
	case "later":

	default:
		return nil, errors.New("undefined type of connection")
	}
	err := Core.Connection.Connect()
	if err != nil {
		return nil, err
	} else {
		return &Core, nil
	}
}

func (cnc *CNCCore) LoadFileForWork(file []byte) error {
	clear(cnc.WorkFile)
	DataFile := string(file)
	if cnc.Connection == nil {
		return errors.New("device is not connected")
	}

	if cnc.DTO.IsWorking {
		return errors.New("printer is already print")
	}
	cnc.WorkFile = strings.Split(DataFile, "\n")
	return nil
}
