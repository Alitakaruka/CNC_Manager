package CNC

import (
	"CNCManager/CNC/CNCService"
	"CNCManager/CNC/CNCService/Connectors"
	"bufio"
	"context"
	"errors"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

var Machines = map[string]AnyCNC{}

const BaseTimeout = 11

type AnyCNC interface {
	GetDTO() CNC_DTO
	SetDTO(DTO CNC_DTO)
	Reconnect() (bool, error)
	GetLogs() []CNCService.Log
	SendMessage(message []byte)
	InitDevice() error
	CNCStart()
	StartTask(file []byte) error
	ExecuteTask(file []byte, ctx context.Context)
	UploadFile(filename string, file []byte)
	CloseConnection()
	ParseCommand(Prefix, dataStr string)
	CanExecuteTask() bool
}

type CNCCore struct {
	DTO CNC_DTO

	//Transmitter   *CNCService.Transmitter
	ReceiveBuffer []byte                  `json:"-"`
	fileBytes     chan int                `json:"-"`
	Mutex         sync.Mutex              `json:"-"`
	WatchDog      *CNCService.WatchDog    `json:"-"`
	Checker       *time.Ticker            `json:"-"`
	Connection    Connectors.CNCConnector `json:"-"`

	Logs     []CNCService.Log
	WorkFile []string `json:"-"`
}

type CNC_DTO struct {
	Position struct {
		X float32 `json:"X"`
		Y float32 `json:"Y"`
		Z float32 `json:"Z"`
	} `json:"Position"`
	Immutable struct {
		Width  int `json:"Width"`
		Length int `json:"Length"`
		Height int `json:"Height"`
	}
	Flags struct {
		WIFI          bool `json:"WIFI"`
		Connected     bool `json:"Connected"`
		ExecutingTask bool `json:"ExecutingTask"`
	} `json:"Flags"`

	Memory struct {
		Buffer         uint32
		FileStorageFMT string
		FileStorage    bool
	}
	Switchable struct {
		Timeout bool
		Light   bool
		RGB     bool
	}

	Device_Chip_Name    string `json:"-"`
	TARGET_MACHINE_NAME string `json:"TARGET_MACHINE_NAME"`
	MACHINE_TYPE        int    `json:"MACHINE_TYPE"`
	FIRMWARE_VERSION    string `json:"FIRMWARE_VERSION"`
	UniqueKey           string `json:"UniqueKey"`
	ConnectionData      string `json:"ConnectionData"`
	ConnectionString    string `json:"-"`
}

func (cnc *CNCCore) GetConnection() {

}

func (cnc *CNCCore) StartTask(file []byte) error {
	err := cnc.LoadFileForWork(file)
	if err != nil {
		return err
	}
	return nil
}

func (cnc *CNCCore) ExecuteTask(file []byte, ctx context.Context) {
	cnc.WriteLog("this CNC can not be executing tasks", CNCService.LogLevelError) //stub
}

func (cnc *CNCCore) CNCStart() {
	cnc.ReceiveBuffer = make([]byte, 512)
	if !cnc.DTO.Switchable.Timeout {
		go cnc.StartWatchcDog()
		go cnc.CheckConnection_Async()
	}
	go cnc.ReadConnectionAsync()
	cnc.WriteLog("Successfully connected", CNCService.LogLevelInformation)
}

func (cnc *CNCCore) StartWatchcDog() {
	cnc.WatchDog = CNCService.NewWatchDog(11, nil)

	if close := cnc.WatchDog.Wait(); close {
		return
	}

	if cnc.isConnected() {
		cnc.WriteLog("The machine timeot!", CNCService.LogLevelError)
		cnc.CloseConnection()
	}

}

func (cnc *CNCCore) InitDevice() error {
	reader := CNCService.NewTimeoutReader(cnc.Connection, time.Second*2)
	cnc.Connection.Write([]byte(CNCService.Identification + CNCService.EndOfData))
	res := reader.Read()
	if res == "" {
		return errors.New("the device did not respond to the request")
	}
	commands := strings.Split(res, CNCService.EndOfData)

	for _, comm := range commands {
		if comm == CNCService.BufferACK {
			continue
		}
		cnc.ParseBaseCommand(comm)
	}

	if cnc.DTO.TARGET_MACHINE_NAME == "" || cnc.DTO.MACHINE_TYPE == 0 {
		return errors.New("the device did not respond as expected")
	}
	cnc.DTO.Flags.Connected = true
	return nil
}

func (cnc *CNCCore) FillDeviceData(str string) error {
	return nil
}

func (cnc *CNCCore) WriteLog(log, logLevel string) {
	if log != "" {
		log = cnc.DTO.TARGET_MACHINE_NAME + ":" + log
		cnc.Logs = append(cnc.Logs, CNCService.Log{Level: logLevel, Message: log})
	}
}

func (cnc *CNCCore) isConnected() bool {
	cnc.Mutex.Lock()
	defer cnc.Mutex.Unlock()
	return cnc.DTO.Flags.Connected
}

func (cnc *CNCCore) CanExecuteTask() bool {
	cnc.Mutex.Lock()
	defer cnc.Mutex.Unlock()
	return cnc.DTO.Flags.Connected && cnc.DTO.Flags.ExecutingTask
}

func (cnc *CNCCore) isExecutingTask() bool {
	cnc.Mutex.Lock()
	defer cnc.Mutex.Unlock()
	return cnc.DTO.Flags.ExecutingTask
}

func (cnc *CNCCore) ReadConnectionAsync() {
	cnc.ReceiveBuffer = cnc.ReceiveBuffer[:0]
	reader := bufio.NewReader(cnc.Connection)
	for cnc.isConnected() {
		Byte, ex := reader.ReadByte()
		if ex != nil {
			cnc.CloseConnection()
			cnc.WriteLog(ex.Error(), CNCService.LogLevelError)
		} else {
			cnc.Mutex.Lock()

			if cnc.WatchDog != nil {
				cnc.WatchDog.Alive()
			}

			cnc.ReceiveBuffer = append(cnc.ReceiveBuffer, Byte)
			cnc.Mutex.Unlock()
		}
	}
}

func (cnc *CNCCore) CheckConnection_Async() {
	cnc.Checker = time.NewTicker(5 * time.Second)
	defer cnc.Checker.Stop()

	for range cnc.Checker.C {

		if !cnc.isConnected() {
			return
		}
		cnc.SendCommand([]byte(CNCService.EndOfData))
	}
}

func (cnc *CNCCore) LoadFileForWork(file []byte) error {
	clear(cnc.WorkFile)
	DataFile := string(file)
	if cnc.Connection == nil {
		return errors.New("device is not connected")
	}
	cnc.WorkFile = strings.Split(DataFile, "\n")
	return nil
}

func (cnc *CNCCore) GetDTO() CNC_DTO {
	return cnc.DTO
}

func (cnc *CNCCore) SetDTO(DTO CNC_DTO) {
	cnc.DTO = DTO
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

	if cnc.Connection == nil {
		log.Println("CNC does not connected")
	}
	_, ex := cnc.Connection.Write(message)
	if ex != nil {
		cnc.WriteLog(ex.Error(), CNCService.LogLevelError)
	}
}

func (cnc *CNCCore) SendCommand(Command []byte) {

	if cnc.Connection == nil {
		log.Println("CNC does not connected")
	}
	_, ex := cnc.Connection.Write(Command)
	if ex != nil {
		cnc.WriteLog(ex.Error(), CNCService.LogLevelError)
	}
	_, ex = cnc.Connection.Write([]byte(CNCService.EndOfData))
	if ex != nil {
		cnc.WriteLog(ex.Error(), CNCService.LogLevelError)
	}

}

func (cnc *CNCCore) GetLogs() []CNCService.Log {
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
			Core.DTO.ConnectionData = "COM"
		}
	case "IP":
		strs := strings.Split(connectionData, ":")
		var ip, port string
		if len(strs) == 2 {
			ip = strs[0]
			port = strs[1]

		} else {
			ip = strings.TrimSpace(connectionData)
			port = "8080"
		}
		Core.Connection = Connectors.NewIpConnector(ip, port)
		Core.DTO.ConnectionData = "WIFI"
	case "later":

	default:
		return nil, errors.New("undefined type of connection")
	}

	err := Core.Connection.Connect()
	Core.DTO.ConnectionString = typeOfConnection + ":" + connectionData
	if err != nil {
		return nil, err
	} else {
		return &Core, nil
	}
}

func (cnc *CNCCore) UploadFile(filename string, file []byte) {
	strCommandStart :=
		CNCService.StartOfTransmision +
			CNCService.FILE_NAME + filename + string('\n') +
			CNCService.FILE_SIZE + strconv.Itoa(len(file)) + string('\n') +
			CNCService.EndOfData

	cnc.SendMessage([]byte(strCommandStart))
	// cnc.WatchDog.Close() todo
	for len(file) != 0 {
		select {
		case bytes := <-cnc.fileBytes:
			if bytes > len(file) {
				bytes = len(file)
			}
			transfer := file[:bytes]
			cnc.SendMessage(transfer)
			file = file[bytes:]
		case <-time.After(time.Second * 5):

		}
	}
	strRes := CNCService.EndOfTransmision + string(CNCService.EndOfData)
	cnc.SendMessage([]byte(strRes))
	// cnc.WatchDog.Reset(time.Second * BaseTimeout)
}

func (cnc *CNCCore) CloseConnection() {
	cnc.Mutex.Lock()
	if cnc.DTO.Flags.Connected {
		cnc.Connection.Close()
		cnc.DTO.Flags.Connected = false

		if cnc.WatchDog != nil {
			cnc.WatchDog.Close()
		}
	}
	cnc.Mutex.Unlock()
}

func RegisterCNC(name string, f func() AnyCNC) {
	Machines[name] = f()
}

func (cnc *CNCCore) ParseBaseCommand(Command string) {
	if len(Command) == 0 {
		return
	}
	prefix := Command[:strings.Index(Command, ":")+1]
	dataStr := strings.TrimSpace(Command[strings.Index(Command, ":")+1:])
	dataF32, _ := strconv.ParseFloat(dataStr, 32)
	dataInt, _ := strconv.Atoi(dataStr)

	switch prefix {
	case CNCService.DEVICE_CHIP_NAME:
		cnc.DTO.Device_Chip_Name = dataStr
	case CNCService.MyName:
		cnc.DTO.TARGET_MACHINE_NAME = dataStr
	case CNCService.MyType:
		cnc.DTO.MACHINE_TYPE = dataInt
	case CNCService.MyPositionX:
		cnc.DTO.Position.X = float32(dataF32)
	case CNCService.MyPositionY:
		cnc.DTO.Position.Y = float32(dataF32)
	case CNCService.MyPositionZ:
		cnc.DTO.Position.Z = float32(dataF32)
	case CNCService.MyWidth:
		cnc.DTO.Immutable.Width = dataInt
	case CNCService.MyLength:
		cnc.DTO.Immutable.Length = dataInt
	case CNCService.MyHeight:
		cnc.DTO.Immutable.Height = dataInt
	case CNCService.SwitchTimeout:
		cnc.DTO.Switchable.Timeout = (dataInt == 1)
	case CNCService.Error:
		cnc.WriteLog(CNCService.LogLevelError, dataStr)
	default:
		cnc.ParseCommand(prefix, dataStr)
	}
}

func (cnc *CNCCore) ParseCommand(Prefix, DataStr string) {
	log.Println("This device cant parce your command! Overload `ParseCommand`!")
}
