package Server

import (
	Service "CNCManager/Service"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	_ "modernc.org/sqlite"
)

//go:embed files/index.html
var ServerFile embed.FS

type CNCServer struct {
	SecondsWork uint32
	Connections uint32

	Hub *Hub

	TimerUpdader  *WSTransmiterr
	StatusUpdader *WSTransmiterr
	TableUpdader  *WSTransmiterr
	LogsUpdader   *WSTransmiterr

	Manager Service.CNCManagerr
	mux     *http.ServeMux
	port    string
	Adrr    string
}

func (PS *CNCServer) InitServer(port, addr, sqlPath string) {

	PS.port = port
	PS.Adrr = addr

	PS.Hub = NewHub()

	Service.InitPrinters()
	PS.Manager.InitManager(sqlPath)
	PS.mux = http.NewServeMux()
	//fs := http.FS(ServerFile)
	//PS.mux.Handle("/", http.FileServer(fs))
	PS.mux.HandleFunc("/", mainHandled)
	// PS.mux.HandleFunc("/connect", PS.ConnectPrinter)
	PS.mux.HandleFunc("/ws", PS.HandleWS)
	// PS.mux.HandleFunc("/api/Printers", PS.GetPrintersInformation)
	// PS.mux.HandleFunc("/api/ExecuteTask", PS.ExecuteTask)
	// PS.mux.HandleFunc("/api/SaveSettings", PS.saveSettings)
	// PS.mux.HandleFunc("/api/GetSettings", PS.getSettings)
	// PS.mux.HandleFunc("/api/sendGCode", PS.SendGCode)

	//
	PS.TimerUpdader = NewWSTransmiterr()
	PS.StatusUpdader = NewWSTransmiterr()
	PS.TableUpdader = NewWSTransmiterr()
	PS.LogsUpdader = NewWSTransmiterr()
}

func CatchPanic(context string) {
	if r := recover(); r != nil {
		log.Printf("[PANIC in %s] %v", context, r)
	}
}

func (PS *CNCServer) CountTime() {
	for {
		<-time.After(time.Second)
		PS.SecondsWork++
	}
}

func (PS *CNCServer) Serve() {
	defer CatchPanic("main")

	go PS.CountTime()
	go PS.UpdateLogs()

	// go PS.UpdateCNCTable()
	go PS.UpdateStatus()
	go PS.UpdateTimer()

	log.Printf("Server started: %s\n", PS.Adrr+":"+PS.port)
	err := http.ListenAndServe(PS.Adrr+":"+PS.port, PS.mux)
	if err != nil {
		log.Println(err)
		panic(err)
	}

}

// HTTP +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (PS *CNCServer) ConnectPrinter(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "undefined query", http.StatusBadRequest)
		log.Println("undefined query")
		return
	}
	information := Service.ConnectionData{}
	ex := json.NewDecoder(r.Body).Decode(&information)
	if ex != nil {
		http.Error(w, "Failed to decode json: "+ex.Error(), http.StatusBadRequest)
		log.Println("Failed to decode json: " + ex.Error())
		return
	}
	ex = PS.Manager.Connect(information)
	if ex != nil {
		http.Error(w, "Failed to connect the CNC due to: "+ex.Error(), http.StatusBadRequest)
		log.Println("Failed to connect the CNC due to: " + ex.Error())
		return
	}
	w.Write([]byte("ok"))
}

func (PS *CNCServer) SendGCode(w http.ResponseWriter, r *http.Request) {
	Gcode := r.URL.Query().Get("GCode")
	uniqueKey := r.URL.Query().Get("uniqueKey")
	log.Println(r.URL.Query())
	if Gcode == "" || uniqueKey == "" {
		http.Error(w, "void parametrs", http.StatusBadRequest)
		return
	}
	ex := PS.Manager.SendGCode(Gcode, uniqueKey)
	if ex != nil {
		http.Error(w, ex.Error(), http.StatusBadRequest)
		return
	}
	w.Write([]byte("ok"))
}

func (PS *CNCServer) ExecuteTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.ParseMultipartForm(32 << 20) // 32MB max
	file, _, ex := r.FormFile("PrintFile")
	if ex != nil {
		http.Error(w, "Failed to get file: "+ex.Error(), http.StatusBadRequest)
		log.Println("Failed to get file: " + ex.Error())
		return
	}
	defer file.Close()

	fileBytes, ex := io.ReadAll(file)
	if ex != nil {
		http.Error(w, "Failed to read file: "+ex.Error(), http.StatusBadRequest)
		log.Println("Failed to read file: " + ex.Error())
		return
	}

	Key := r.URL.Query().Get("uniqueKey")
	if Key == "" {
		http.Error(w, "uniqueKey parameter is required", http.StatusBadRequest)
		return
	}

	ex = PS.Manager.ExecuteTask(Key, fileBytes)
	if ex != nil {
		http.Error(w, "Failed to execute task: "+ex.Error(), http.StatusBadRequest)
		log.Println("Failed to execute task: " + ex.Error())
		return
	}

	w.Write([]byte("Start printing!"))
}

func (PS *CNCServer) UploadFile(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	file, _, ex := r.FormFile("PrintFile")
	if ex != nil {
		http.Error(w, ex.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()
	fileBytes, ex := io.ReadAll(file)
	if ex != nil {
		http.Error(w, ex.Error(), http.StatusBadRequest)
	}
	Key := r.URL.Query().Get("uniqueKey")
	PS.Manager.UploadFile(Key, "test.gcode", fileBytes)
	w.Write([]byte("Start printing!"))
}

func (PS *CNCServer) GetPrintersInformation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(PS.Manager.GetJson()))
}

func mainHandled(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "WEB/Server/files/index.html")
}

//HTTP ----------------------------------------------------------------------

// WEB SOCKET ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (PS *CNCServer) HandleWS(w http.ResponseWriter, r *http.Request) {
	ServeWs(PS.Hub, w, r)

	// upgrader := websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
	// 	return true
	// }, EnableCompression: true}
	// WebSoc, err := upgrader.Upgrade(w, r, nil)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// }
	// WebSoc.SetPingHandler(func(appdata string) error {
	// 	return WebSoc.WriteControl(websocket.PongMessage, []byte(appdata), time.Now().Add(5))
	// })
	// WebSoc.SetPongHandler(func(appdata string) error {
	// 	if appdata != "Ping" {
	// 		return errors.New("appdata error")
	// 	}
	// 	return nil
	// })

	// isClose := make(chan struct{})
	// WebSoc.SetCloseHandler(func(code int, text string) error {
	// 	PS.Connections--
	// 	close(isClose)
	// 	return nil
	// })
	// PS.Connections++
	// var mut sync.Mutex
	// //Reader
	// go func() {
	// 	ticker := time.NewTicker(PingTime)

	// 	go func() {
	// 		for {
	// 			select {
	// 			case <-ticker.C:
	// 				WebSoc.WriteControl(websocket.PingMessage, []byte("Ping"), time.Now().Add(PingTime))
	// 				ticker.Reset(time.Second * 5)
	// 			case <-isClose:
	// 				return
	// 			}
	// 		}
	// 	}()
	// 	for {
	// 		select {
	// 		case <-isClose:
	// 			return
	// 		default:
	// 			msgType, msg, err := WebSoc.ReadMessage()
	// 			if err != nil {
	// 				return
	// 			}
	// 			_ = msgType
	// 			if msgType == websocket.TextMessage {
	// 				Responce := PS.ExecuteWSMessage(string(msg), WebSoc)
	// 				mut.Lock()
	// 				WebSoc.WriteMessage(websocket.TextMessage, Responce)
	// 				mut.Unlock()
	// 			}
	// 		}
	// 	}
	// }()

	// go PS.WSWriteData(WebSoc, isClose, &mut)
}

const PingTime = time.Second * 5

func (PS *CNCServer) DataUpdater() {
	go func() {
		PS.TimerUpdader.WaitNewData()
		PS.Hub.Send([]byte(PS.TimerUpdader.GetNowData()))
	}()
	go func() {
		PS.StatusUpdader.WaitNewData()
		PS.Hub.Send([]byte(PS.StatusUpdader.GetNowData()))
	}()
	go func() {
		<-PS.Manager.IsCharge
		cncsJson := PS.Manager.GetJson()
		PS.Hub.Send([]byte(cncsJson))
	}()
	go func() {
		PS.LogsUpdader.WaitNewData()
		PS.Hub.Send([]byte(PS.LogsUpdader.GetNowData()))
	}()
}

func (PS *CNCServer) ExecuteWSMessage(msg string, WS *websocket.Conn) []byte {
	var mas WSMessage
	err := json.Unmarshal([]byte(msg), &mas)
	if err != nil {
		log.Println(err)
	}

	switch mas.Type {
	case "connect":
		con := Service.ConnectionData{}
		json.Unmarshal(mas.Data, &con)
		err := PS.Manager.Connect(con)

		if err != nil {
			return WEB_Socket_ERROR(mas.ReqId, err.Error())
		}
		return WEB_Socket_ACK(mas.ReqId, true)
	case "reconnect":
		con := struct {
			UniqueKey string `json:"uniqueKey"`
		}{}
		json.Unmarshal(mas.Data, &con)
		err := PS.Manager.Reconnect(con.UniqueKey)
		if err != nil {
			return WEB_Socket_ERROR(mas.ReqId, err.Error())
		} else {
			return WEB_Socket_ACK(mas.ReqId, true)
		}
	case "GetMachines":
		PS.TableUpdader.SetNewData("") //todo костыль
		// cncsJson += string(bytes)
		return WEB_Socket_ACK(mas.ReqId, true)
	case "command":
		Command := struct {
			Gcode     string `json:"gcode"`
			UniqueKey string `json:"uniqueKey"`
		}{}

		err := json.Unmarshal(mas.Data, &Command)
		if err != nil {
			return WEB_Socket_ERROR(mas.ReqId, err.Error())
		}
		err = PS.Manager.SendGCode(Command.Gcode, Command.UniqueKey)
		if err != nil {
			return WEB_Socket_ERROR(mas.ReqId, err.Error())
		}
		return WEB_Socket_ACK(mas.ReqId, true)

	case "executeTask":
		task := struct {
			UniqueKey string `json:"uniqueKey"`
			FileName  string `json:"fileName"`
			FileData  []byte `json:"fileData"`
		}{}
		err := json.Unmarshal(mas.Data, &task)
		if err != nil {
			return WEB_Socket_ERROR(mas.ReqId, err.Error())
		}

		/////////////////////////////
		// fmt.Printf("task.FileData: %v\n", string(task.FileData))
		// enc := base64.StdEncoding.EncodeToString(task.FileData)

		// dec, err := base64.StdEncoding.DecodeString(enc)
		// if err != nil {
		// 	log.Fatal(err)
		// }

		////////////////////////////

		// fmt.Printf("task.FileData: %v\n", task.FileData)
		// fmt.Printf("base64.StdEncoding.EncodeToString((task.FileData)): %v\n", base64.StdEncoding.EncodeToString((task.FileData)))
		err = PS.Manager.ExecuteTask(task.UniqueKey, (task.FileData))
		if err != nil {
			return WEB_Socket_ERROR(mas.ReqId, err.Error())
		}
		return WEB_Socket_ACK(mas.ReqId, true)
	}

	return []byte{}
}

// func (PS *CNCServer) UpdateCNCTable() {
// 	for {
// 		// Send CNC machines data
// 		<-PS.Manager.IsCharge
// 		cncsJson := PS.Manager.GetJson()
// 		if cncsJson != "" && cncsJson != "[]" {
// 			cncsMsg := fmt.Sprintf(`{"type":"printers","data":%s}`, cncsJson)
// 			PS.TableUpdader.SetNewData(cncsMsg)
// 		}
// 		time.Sleep(50 * time.Millisecond)
// 	}
// }

func (PS *CNCServer) UpdateTimer() {
	for {
		curTime := PS.SecondsWork
		days := curTime / 86400
		curTime %= 86400
		hours := curTime / 3600
		curTime %= 3600
		minutes := curTime / 60
		curTime %= 60

		dataSec := struct {
			Uptime            string `json:"uptime"`
			ActiveConnections int    `json:"activeConnections"`
		}{Uptime: fmt.Sprintf("%dd %dh %dm", days, hours, minutes),
			ActiveConnections: int(PS.Connections)}
		jsonData, _ := json.Marshal(dataSec)
		resp := WSMessage{
			Type: "systemInfo",
			Data: jsonData,
		}
		res, _ := json.Marshal(resp)
		PS.TimerUpdader.SetNewData(string(res))
		time.Sleep(500 * time.Millisecond)
	}
}

func (PS *CNCServer) UpdateStatus() {
	for {
		// Send status
		online := 0
		printing := 0
		for _, machine := range PS.Manager.CNC_Machines {
			dto := machine.GetDTO()
			if dto.Flags.Connected {
				online++
				if dto.Flags.ExecutingTask {
					printing++
				}
			}
		}

		statusMsg := fmt.Sprintf(
			`{"type":"status","data":{"online":%d,"printing":%d,"offline":%d,"total":%d}}`,
			online,
			printing,
			len(PS.Manager.CNC_Machines)-online,
			len(PS.Manager.CNC_Machines),
		)
		PS.StatusUpdader.SetNewData(statusMsg)

	}
}

func (PS *CNCServer) UpdateLogs() {
	id := uint32(1)
	for {
		Logs := PS.Manager.GetAllLogs()
		for _, val := range Logs {
			js := WEB_Socket_LOG(id, uint32(time.Now().Unix()), val.Level, val.Message)
			id++
			PS.LogsUpdader.SetNewData(string(js))
		}
		time.Sleep(time.Second)
	}
}
