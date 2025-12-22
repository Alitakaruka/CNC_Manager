package CNC

import (
	"CNCManager/CNC/CNCService"
	"CNCManager/CNC/CNCService/Connectors"
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var Machines = map[string]RealizeCNC{}

const BaseTimeout = 11

type RealizeCNC interface {
	// AnyCNC
	ExecuteTask(file []byte, ctx context.Context)
	ParseCommand(Prefix, dataStr string)
	InitRealization() error
	GetJsonData() any
	SetCore(core *CNCCore)
}

type noCopy struct{}

func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}

type CNCCore struct {
	DTO     CNC_DTO
	_       noCopy
	Realize RealizeCNC `json:"-"`
	// ReceiveBuffer  []byte                  `json:"-"`
	ReceiveBuffer chan byte               `json:"-"`
	fileBytes     chan int                `json:"-"`
	Mutex         *sync.RWMutex           `json:"-"`
	WatchDog      *CNCService.WatchDog    `json:"-"`
	Checker       *time.Ticker            `json:"-"`
	Transmitter   *CNCService.Transmitter `json:"-"`
	Connection    Connectors.CNCConnector `json:"-"`

	Logs     []CNCService.Log
	LogFile  *os.File `json:"-"`
	Progress int      `json:"_"`
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

func (cnc *CNCCore) StartTask(file []byte) error {
	// err := cnc.LoadFileForWork(file)
	// if err != nil {
	// 	return err
	// }
	cntx, cancel := context.WithCancel(context.Background())
	go func() {
		cnc.Mutex.Lock()
		dto := cnc.GetDTO()
		dto.Flags.ExecutingTask = true
		cnc.SetDTO(dto)
		cnc.Mutex.Unlock()

		cnc.Realize.ExecuteTask(file, cntx)

		if cnc.CanExecuteTask() {
			cnc.WriteLog(CNCService.LogLevelSuccess, "Task executing successeful!")
		} else {
			cnc.WriteLog(CNCService.LogLevelError, "Task failed successfully!")
		}
		cnc.Mutex.Lock()
		dto.Flags.ExecutingTask = false
		cnc.SetDTO(dto)
		cnc.Mutex.Unlock()
	}()
	time.Sleep(time.Millisecond * 10) //stub
	go func() {
		for cnc.CanExecuteTask() {
			time.Sleep(time.Millisecond * 100)
		}
		cancel()
		cnc.WriteLog(CNCService.LogLevelInformation, "End of executing task!")
	}()

	cnc.WriteLog(CNCService.LogLevelSuccess, "task start!")
	return nil
}

func (cnc *CNCCore) CNCStart() {
	cnc.Transmitter.SyncBuffers(cnc.Connection)
	cnc.ReceiveBuffer = make(chan byte, 1024)
	cnc.CreateLogFile()
	if !cnc.DTO.Switchable.Timeout {
		go cnc.StartWatchcDog()
		go cnc.CheckConnection_Async()
	}

	go cnc.readConnectionAsync()
	go cnc.readResponces()
	cnc.WriteLog(CNCService.LogLevelSuccess, "Successfully connected")
}

func (cnc *CNCCore) readResponces() {
	for cnc.isConnected() {
		if Command := cnc.getNextByteStream(); Command != nil {
			cnc.parseCommand(string(Command))
		}
	}
}

func (cnc *CNCCore) StartWatchcDog() {
	cnc.WatchDog = CNCService.NewWatchDog(11, nil)
	if close := cnc.WatchDog.Wait(); close {
		return
	}
	if cnc.isConnected() {
		cnc.WriteLog(CNCService.LogLevelError, "The machine timeot!")
		cnc.CloseConnection()
	}

}

func (cnc *CNCCore) InitDevice() error {
	cnc.Transmitter = CNCService.NewTransmitter()

	reader := CNCService.NewTimeoutReader(cnc.Connection, time.Second*2)
	cnc.SendMessage([]byte(CNCService.Identification + CNCService.EndOfData))
	res := reader.Read()
	if res == "" {
		cnc.Connection.Close()
		return errors.New("the device did not respond to the request")
	}
	commands := strings.Split(res, CNCService.EndOfData)
	fmt.Printf("commands: %v\n", commands)
	for _, comm := range commands {
		cnc.parseCommand(comm)
	}

	if cnc.DTO.TARGET_MACHINE_NAME == "" || cnc.DTO.MACHINE_TYPE == 0 {
		cnc.Connection.Close()
		return errors.New("the device did not respond as expected")
	}

	if targer, ok := Machines[cnc.DTO.Device_Chip_Name]; !ok {
		cnc.Connection.Close()
		return errors.New("the device dint register")
	} else {
		targer.SetCore(cnc)
		cnc.Realize = targer
		cnc.DTO.Flags.Connected = true
		return cnc.Realize.InitRealization()
	}
}

func (cnc *CNCCore) WriteLog(logLevel, Log string) {
	if Log != "" {
		Log = cnc.DTO.TARGET_MACHINE_NAME + ":" + Log
		cnc.Logs = append(cnc.Logs, CNCService.Log{Level: logLevel, Message: Log})
	}
}

func (cnc *CNCCore) isConnected() bool {
	cnc.Mutex.RLock()
	con := cnc.DTO.Flags.Connected
	cnc.Mutex.RUnlock()
	return con
}

func (cnc *CNCCore) CanExecuteTask() bool {
	cnc.Mutex.RLock()
	can := cnc.DTO.Flags.Connected && cnc.DTO.Flags.ExecutingTask
	cnc.Mutex.RUnlock()
	return can
}

func (cnc *CNCCore) readConnectionAsync() {
	// cnc.ReceiveBuffer = cnc.ReceiveBuffer[:0]
	reader := bufio.NewReader(cnc.Connection)
	for cnc.isConnected() {
		Byte, ex := reader.ReadByte()
		if ex != nil {
			cnc.CloseConnection()
			cnc.WriteLog(CNCService.LogLevelError, ex.Error())
		} else {
			cnc.Mutex.Lock()
			if cnc.WatchDog != nil {
				cnc.WatchDog.Alive()
			}
			cnc.ReceiveBuffer <- Byte
			// cnc.ReceiveBuffer = append(cnc.ReceiveBuffer, Byte)
			cnc.Mutex.Unlock()
		}
	}
}

func (cnc *CNCCore) CheckConnection_Async() {
	Checker := time.NewTicker(5 * time.Second)
	defer Checker.Stop()

	for range Checker.C {
		if cnc.DTO.Flags.ExecutingTask {
			continue
		}
		if !cnc.isConnected() {
			return
		}
		cnc.SendMessage([]byte(CNCService.EndOfData))
	}
}

func (cnc *CNCCore) CreateLogFile() {
	var err error
	cnc.LogFile, err = os.OpenFile(cnc.DTO.TARGET_MACHINE_NAME+".log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0666)
	if err != nil {
		// panic(err)
	}
}

func (cnc *CNCCore) LoadFileForWork(file []byte) error {
	// clear(cnc.WorkFile)
	// DataFile := string(file)
	// if cnc.Connection == nil {
	// 	return errors.New("device is not connected")
	// }
	// cnc.WorkFile = strings.Split(DataFile, "\n")
	return nil
}

func (cnc *CNCCore) GetDTO() CNC_DTO {
	return cnc.DTO
}

func (cnc *CNCCore) SetDTO(DTO CNC_DTO) {
	cnc.DTO = DTO
}

func (cnc *CNCCore) getNextByteStream() []byte {
	result := make([]byte, 0)
	// cnc.Mutex.RLock()

	for Data := range cnc.ReceiveBuffer {
		// cnc/ReceiveBuffer
		result = append(result, Data)
		if strings.HasSuffix(string(result), CNCService.EndOfData) {
			return result
		}
	}
	return result
	// ind := strings.Index(string(cnc.ReceiveBuffer), CNCService.EndOfData)
	// if ind == -1 {
	// 	// fmt.Println("index -1")
	// 	cnc.Mutex.RUnlock()
	// 	return nil
	// }
	// Command := cnc.ReceiveBuffer[:ind]
	// // fmt.Printf("Command: %v\n", Command)
	// cnc.Mutex.RUnlock()
	// cnc.Mutex.Lock()
	// cnc.ReceiveBuffer = cnc.ReceiveBuffer[ind+len([]rune(CNCService.EndOfData)):]
	// // fmt.Printf("cnc.ReceiveBuffer: %v\n", cnc.ReceiveBuffer)
	// cnc.Mutex.Unlock()
	// return Command
}

func (cnc *CNCCore) SendMessage(message []byte) bool {
	if cnc.Transmitter.Wait(len(message)) {
		cnc.Transmitter.Trainsmit(len(message))
		_, ex := cnc.Connection.Write(message)
		if ex != nil {
			cnc.WriteLog(CNCService.LogLevelError, ex.Error())
			return false
		}
		return true
	} else {
		cnc.WriteLog(CNCService.LogLevelWarning, "Device command size too large!The command will be ignored!")
		return false
	}
}

func (cnc *CNCCore) GetLogs() []CNCService.Log {
	cnc.Mutex.Lock()
	copy := cnc.Logs[:len(cnc.Logs)]
	cnc.Logs = cnc.Logs[:0]
	cnc.Mutex.Unlock()
	return copy
}

func (cnc *CNCCore) Reconnect() error {
	// fmt.Println("Reconnect!")
	err := cnc.Connection.Reconnect()
	if err != nil {
		return err
	}
	return nil
}

func Connect(typeOfConnection string, connectionData string) (*CNCCore, error) {
	Core := CNCCore{Mutex: &sync.RWMutex{}}
	strs := strings.Split(connectionData, ":")

	switch typeOfConnection {
	case "COM":
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
	case "IP", "WIFI":
		var ip, port string
		if len(strs) == 2 {
			ip = strs[0]
			port = strs[1]

		} else {
			ip = strings.TrimSpace(connectionData)
			port = "8080"
		}
		Core.Connection = Connectors.NewIpConnector(ip, port)
		Core.DTO.ConnectionData = "IP"
	case "later":

	default:
		return nil, errors.New("undefined type of connection")
	}

	err := Core.Connection.Connect()
	Core.DTO.ConnectionString = connectionData
	if err != nil {
		return nil, err
	} else {
		return &Core, nil
	}
}

func GetConnector(ConData, ConString string) Connectors.CNCConnector {
	switch ConData {
	case "COM":
		strs := strings.Split(ConString, ":")
		var port, Baud string
		if len(strs) == 2 {
			port = strs[0]
			Baud = strs[1]
			BaudRate, err := strconv.Atoi(Baud)
			if err != nil {
				return nil
			}
			return Connectors.NewSerialConnector(port, BaudRate)
		} else if len(strs) == 1 {
			return Connectors.NewSerialConnector(ConString, 9600)
		}
	case "IP":
		strs := strings.Split(ConString, ":")
		var ip, port string
		if len(strs) == 2 {
			ip = strs[0]
			port = strs[1]

		} else {
			ip = strings.TrimSpace(ConString)
			port = "8080"
		}
		return Connectors.NewIpConnector(ip, port)
	case "later":
	}
	return nil
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

func RegisterCNC(name string, f func() RealizeCNC) {
	Machines[name] = f()
}

func (cnc *CNCCore) parseCommand(Command string) {
	Command, _ = strings.CutSuffix(Command, CNCService.EndOfData)
	if len(Command) == 0 {
		return
	}
	if Command == CNCService.BufferACK {
		cnc.Transmitter.ACK()
	}

	prefix := Command[:strings.Index(Command, ":")+1]
	dataStr := strings.TrimSpace(Command[strings.Index(Command, ":")+1:])
	dataF32, _ := strconv.ParseFloat(dataStr, 32)
	dataInt, _ := strconv.Atoi(dataStr)

	// fmt.Printf("Command: %v\n", Command)
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
	case CNCService.SwitchHasLight:
		cnc.DTO.Switchable.RGB = (dataInt == 1)
	case CNCService.SwitchRGBLight:
		cnc.DTO.Switchable.Light = (dataInt == 1)
	case CNCService.Error:
		cnc.WriteLog(CNCService.LogLevelError, dataStr)
		cnc.LogFile.Write([]byte(time.Now().Format("dd.mm.yy") + "  Error:" + dataStr + "\n"))
	case CNCService.Warning:
		cnc.WriteLog(CNCService.LogLevelWarning, dataStr)
		cnc.LogFile.Write([]byte(time.Now().Format("dd.mm.yy") + ":  Warning:" + dataStr + "\n"))
	case CNCService.Information:
		cnc.LogFile.Write([]byte(time.Now().Format("dd.mm.yy") + "  Info:" + dataStr + "\n"))
		cnc.WriteLog(CNCService.LogLevelInformation, dataStr)
	case CNCService.Success:
		cnc.LogFile.Write([]byte(time.Now().Format("dd.mm.yy") + "  Success:" + dataStr + "\n"))
		cnc.WriteLog(CNCService.LogLevelSuccess, dataStr)
	default:
		if cnc.Realize != nil {
			cnc.Realize.ParseCommand(prefix, dataStr)
		}
	}
}
