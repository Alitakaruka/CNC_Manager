package CNC

import (
	"CNCManager/CNC/CNCService"
	"CNCManager/CNC/CNCService/Connectors"
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"time"
)

var Machines = map[string]RealizeCNC{}

const BaseTimeout = 11

type RealizeCNC interface {
	// AnyCNC
	ExecuteTask(file []byte)
	ParseCommand(Prefix, dataStr string)
	InitRealization() error
	GetJsonData() any
	SetCore(core *CNCCore)
}

type noCopy struct{}

func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}

type CNCCore struct {
	DTO CNC_DTO
	_   noCopy

	Realize RealizeCNC `json:"-"`
	// ReceiveBuffer  []byte                  `json:"-"`
	ReceiveBuffer chan byte `json:"-"` //todo small
	Commands      chan string

	// fileBytes   chan int                `json:"-"`
	mutex       sync.RWMutex            `json:"-"`
	WatchDog    *CNCService.WatchDog    `json:"-"`
	Checker     *time.Ticker            `json:"-"`
	Transmitter *CNCService.Transmitter `json:"-"`
	Connection  Connectors.CNCConnector `json:"-"`

	// FileTransmitter struct {
	// 	*CNCService.Transmitter
	// }

	FileTransmitter2 *CNCService.FileTransmitter
	FileTransmitter  *CNCService.Transmitter
	LogFile          *os.File `json:"-"`
	Progress         float32  `json:"_"`

	//General perpose
	//Flags
	// isInitEnd bool

	//events
	IsCharge  chan struct{}
	IsClose   chan struct{}
	IsTaskEnd chan struct{}

	Logs chan CNCService.Log
}

func NewCNCCore() *CNCCore {
	Core := CNCCore{IsCharge: make(chan struct{}, 1),
		ReceiveBuffer: make(chan byte, 1024),
		Logs:          make(chan CNCService.Log, 1024),
		IsClose:       make(chan struct{}, 1),
		IsTaskEnd:     make(chan struct{})}
	return &Core
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
	FIRMWARE            string `json:"FIRMWARE"`
	UniqueKey           string `json:"UniqueKey"`
	ConnectionData      string `json:"ConnectionData"`
	ConnectionType      string `json:"ConnectionType"`
}

func (cnc *CNCCore) StartTask(file []byte) error {

	// cntx, cancel := context.WithCancel(context.Background())
	go func() {
		cnc.mutex.Lock()
		dto := cnc.GetDTO()
		dto.Flags.ExecutingTask = true
		cnc.SetDTO(dto)
		cnc.mutex.Unlock()
		cnc.Realize.ExecuteTask(file)
	}()
	time.Sleep(time.Millisecond * 10) //stub
	go func() {
		<-cnc.IsTaskEnd
		cnc.WriteLog(CNCService.LogLevelInformation, "End of executing task!")
	}()

	cnc.WriteLog(CNCService.LogLevelSuccess, "task start!")
	return nil
}

func (cnc *CNCCore) CNCStart() {
	cnc.CreateLogFile()
	if !cnc.DTO.Switchable.Timeout {
		go cnc.StartWatchcDog()
		go cnc.CheckConnection_Async()
	}

	// go cnc.readConnectionAsync()
	go cnc.readResponces()
	go func() {
		cnc.Connection.WaitClosed()
		log.Println("The machine was disconect!" + cnc.DTO.TARGET_MACHINE_NAME)
		cnc.CloseConnection("Connection closed")
	}()

	cnc.SyncBuffers()
}

func (cnc *CNCCore) SyncBuffers() {
	cnc.SendMessage([]byte(CNCService.EndOfData + CNCService.SYNC + CNCService.EndOfData))
}

func (cnc *CNCCore) readResponces() {
	for {
		select {
		case <-cnc.IsClose:
			// fmt.Printf("\"readResponces\": %v\n", "readResponces close")
			return
		default:
			if Command := cnc.getNextByteStream(); Command != nil {
				cnc.parseCommand(string(Command))
			}
		}
	}
}

func (cnc *CNCCore) StartWatchcDog() {
	cnc.WatchDog = CNCService.NewWatchDog(10, nil)
	for {
		select {
		case <-cnc.WatchDog.Wait():
			cnc.WriteLog(CNCService.LogLevelError, "The machine timeot!")
			cnc.WatchDog.Close()
			cnc.CloseConnection("Watch dog timeout")
		case <-cnc.IsClose:
			cnc.WatchDog.Close()
			return
		}
	}
}

func (cnc *CNCCore) InitDevice() error {
	cnc.Transmitter = CNCService.NewTransmitter()
	cnc.FileTransmitter = CNCService.NewTransmitter()
	cnc.FileTransmitter.SetLimits(1, 1)

	go cnc.readConnectionAsync()
	cnc.SendMessage([]byte(CNCService.EndOfData + CNCService.Identification + CNCService.EndOfData))

	var Data []byte
	stop := false

	timeout := time.After(time.Second * 5)
	for !stop {
		select {
		case <-time.After(time.Second * 2):
			stop = true
		case <-timeout:
			stop = true
		case b := <-cnc.ReceiveBuffer:
			Data = append(Data, b)
		}
	}
	// fmt.Println("Stop ident!")

	res := string(Data)
	// fmt.Printf("res: %v\n", res)
	// fmt.Printf("res: %v\n", []byte(res))

	if res == "" {
		cnc.CloseConnection("the device did not respond to the request")
		return errors.New("the device did not respond to the request")
	}
	commands := strings.Split(res, CNCService.EndOfData)
	// fmt.Printf("commands: %v\n", commands)
	for _, comm := range commands {
		cnc.parseCommand(comm)
	}

	// fmt.Println("Parce end!")
	if cnc.DTO.TARGET_MACHINE_NAME == "" || cnc.DTO.MACHINE_TYPE == 0 {
		cnc.CloseConnection("the device did not respond as expected")
		return errors.New("the device did not respond as expected")
	}

	if targer, ok := Machines[cnc.DTO.Device_Chip_Name]; !ok {
		cnc.CloseConnection("the device dint register")
		return errors.New("the device dint register")
	} else {

		cnc.ModifyCharge()
		targer.SetCore(cnc)
		cnc.Realize = targer
		cnc.DTO.Flags.Connected = true
		// fmt.Println("realization init start")
		err := cnc.Realize.InitRealization() //todo это потом поправить
		if err != nil {
			cnc.CloseConnection(err.Error())
			return err
		}
		// cnc.isInitEnd = true
		return err
	}
}

func (cnc *CNCCore) WriteLog(logLevel, Log string) {
	if Log != "" {
		Log = cnc.DTO.TARGET_MACHINE_NAME + ":" + Log
		cnc.Logs <- CNCService.Log{Level: logLevel, Message: Log}
	}
	cnc.ModifyCharge()
}

func (cnc *CNCCore) readConnectionAsync() {
	// cnc.ReceiveBuffer = cnc.ReceiveBuffer[:0]
	reader := bufio.NewReader(cnc.Connection)

	for {
		select {
		case <-cnc.IsClose:
			// fmt.Println("stop reading!")
			return
		default:
			Byte, ex := reader.ReadByte()

			if ex != nil && ex != io.EOF {
				cnc.WriteLog(CNCService.LogLevelError, ex.Error())
				log.Println(ex)
				cnc.CloseConnection(ex.Error())
				// fmt.Printf("Byte: %v\n", Byte)
			}

			if cnc.WatchDog != nil {
				cnc.WatchDog.Alive()
			}
			select {
			case cnc.ReceiveBuffer <- Byte:
			default:
				log.Println("ReceiveBuffer overflow!")
				time.Sleep(time.Second)
			}

		}
	}

}

func (cnc *CNCCore) CheckConnection_Async() {
	cnc.Checker = time.NewTicker(5 * time.Second)
	defer cnc.Checker.Stop()
	for {
		select {
		case <-cnc.IsClose:
			return
		case <-cnc.Checker.C:
			if cnc.DTO.Flags.ExecutingTask {
				continue
			}
			cnc.SendMessage([]byte(CNCService.EndOfData))
		}
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
	// result := make([]byte, 0)
	var result []byte
	// cnc.mutex.RLock()

	for {
		select {
		case <-cnc.IsClose:
			return nil
		case b := <-cnc.ReceiveBuffer:
			// fmt.Println(b)
			result = append(result, b)
			if bytes.HasSuffix(result, []byte(CNCService.EndOfData)) {
				return result
			}
		}
	}
	// for Data := range cnc.ReceiveBuffer {
	// 	// cnc/ReceiveBuffer
	// 	result = append(result, Data)
	// 	if strings.HasSuffix(string(result), CNCService.EndOfData) {
	// 		return result
	// 	}
	// }
	// return result
}

func (cnc *CNCCore) SendMessage(message []byte) {
	if cnc.Transmitter.Wait(len(message)) {
		cnc.Transmitter.Trainsmit(len(message))
		_, ex := cnc.Connection.Write(message)
		// if len(message) > 0 {
		// 	// log.Printf("I send:%v", string(message))
		// 	// fmt.Printf("cnc.Transmitter.CurrentFreeBytes: %v\n", cnc.Transmitter.CurrentFreeBytes)
		// 	// fmt.Printf("cnc.Transmitter.MaxBytes: %v\n", cnc.Transmitter.MaxBytes)
		// }
		if ex != nil {
			cnc.CloseConnection(ex.Error())
			cnc.WriteLog(CNCService.LogLevelError, ex.Error())
		}
	} else {
		cnc.WriteLog(CNCService.LogLevelWarning, "Device command size too large!The command will be ignored!")
	}
}

func (cnc *CNCCore) Reconnect() error {
	// New chans
	select {
	case <-cnc.IsClose:
		// cnc.Logs = make(chan CNCService.Log, 1024)
		cnc.IsClose = make(chan struct{}, 1)
		cnc.IsTaskEnd = make(chan struct{})
		cnc.IsCharge = make(chan struct{})
	default:

	}
	err := cnc.Connection.Connect()
	if err != nil {
		return err
	}
	log.Println(cnc.DTO.TARGET_MACHINE_NAME + ": connection restored!")
	return nil
}

func Connect(typeOfConnection string, connectionData string) (*CNCCore, error) {
	// fmt.Println("Connect(typeOfConnection string, connectionData string)")
	Core := NewCNCCore()
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
	case "later":

	default:
		return nil, errors.New("undefined type of connection")
	}

	Core.DTO.ConnectionType = typeOfConnection
	Core.DTO.ConnectionData = connectionData
	err := Core.Connection.Connect()
	if err != nil {
		log.Println("Connection error:" + err.Error())
		return nil, err
	} else {
		return Core, nil
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

func (cnc *CNCCore) UploadFile(filename string, base64file []byte) {

}

func (cnc *CNCCore) UploadFileInMemory(base64file []byte) {

	pg := make(chan string)
	go cnc.CreatePages(pg, base64file)

	for {
		data := <-pg
		if data != "" {
			if cnc.FileTransmitter.Wait(1) {
				cnc.FileTransmitter.Trainsmit(1)
				log.Println(string(debug.Stack()))
				log.Printf("cnc.FileTransmitter.MaxBytes: %v\n", cnc.FileTransmitter.MaxBytes)
				cnc.SendMessage([]byte(CNCService.Filedata + data + CNCService.EndOfData))
			} else {
				log.Println(debug.Stack())
			}
		}
	}
}

func (cnc *CNCCore) CreatePages(pg chan string, base64Str []byte) {
	MaxBytes := cnc.Transmitter.MaxBytes - len(CNCService.Filedata) - len(CNCService.EndOfData)
	MaxBytes -= MaxBytes % 4 //base64  encoded data = 4

	fmt.Printf("MaxBytes: %v\n", MaxBytes)
	log.Println(len(base64Str) / MaxBytes)

	Page := CNCService.Filedata
	for _, val := range base64Str {
		Page += string(val)
		if len(Page) == MaxBytes {
			pg <- Page
			Page = CNCService.Filedata
		}
	}
	close(pg)
}

func (cnc *CNCCore) CloseConnection(cause string) {
	log.Printf("Connection closed. Cause: %v\n", cause)
	cnc.mutex.Lock()
	defer cnc.mutex.Unlock()
	select {
	case <-cnc.IsClose:
		return
	default:

	}
	cnc.Connection.Close()
	cnc.DTO.Flags.Connected = false

	if cnc.WatchDog != nil {
		cnc.WatchDog.Close()
	}
	close(cnc.IsClose)
	close(cnc.IsCharge)
	close(cnc.IsTaskEnd)
	cnc.Progress = 0
	// close(cnc.Logs)
	cnc.WriteLog(CNCService.LogLevelError, "The device was close!")
	// cnc.isInitEnd = false
}

func RegisterCNC(name string, f func() RealizeCNC) {
	Machines[name] = f()
}

func (cnc *CNCCore) ModifyCharge() {
	select {
	case <-cnc.IsClose:
		return
	default:
	}
	select {
	case cnc.IsCharge <- struct{}{}:
	default:
	}
}

func (cnc *CNCCore) sendFilePage(Command string) {
	var DataLen, Offset int
	_, err := fmt.Sscanf(Command, CNCService.GetNewFileData, &DataLen, &Offset)
	if err != nil {
		fileData := cnc.FileTransmitter2.GetNewPage(DataLen, Offset)
		resultStr := fmt.Sprintf(CNCService.FileDataRecieve, DataLen, Offset, string(fileData))
		go cnc.SendMessage([]byte(resultStr))
	}
}

func (cnc *CNCCore) parseCommand(Command string) {
	Copy := cnc.GetDTO()
	Command, _ = strings.CutSuffix(Command, CNCService.EndOfData)
	if len(Command) == 0 {
		return
	}
	if Command == CNCService.BufferACK {
		cnc.Transmitter.ACK()
		return
	}
	if Command == CNCService.FileACK {
		cnc.FileTransmitter.ACK()
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
	case CNCService.SwitchHasLight:
		cnc.DTO.Switchable.RGB = (dataInt == 1)
	case CNCService.SwitchRGBLight:
		cnc.DTO.Switchable.Light = (dataInt == 1)
	case CNCService.Error:
		cnc.WriteLog(CNCService.LogLevelError, dataStr)
	case CNCService.Warning:
		cnc.WriteLog(CNCService.LogLevelWarning, dataStr)
	case CNCService.Information:
		cnc.WriteLog(CNCService.LogLevelInformation, dataStr)
	case CNCService.Success:
		cnc.WriteLog(CNCService.LogLevelSuccess, dataStr)
	case CNCService.MyBufferLen:
		cnc.Transmitter.SetLimits(dataInt, dataInt)
	case CNCService.FileReadPrefix:
		cnc.sendFilePage(Command)
	default:
		if cnc.Realize != nil {
			cnc.Realize.ParseCommand(prefix, dataStr)
		}
	}
	if Copy == cnc.GetDTO() {
		cnc.ModifyCharge()
	}
}
