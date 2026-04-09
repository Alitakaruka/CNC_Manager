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

	_ "modernc.org/sqlite"
)

//go:embed files/index.html
var ServerFile embed.FS

type CNCServer struct {
	Hub     *Hub
	Manager Service.CNCManager
	mux     *http.ServeMux
	port    string
	Adrr    string
}

func (PS *CNCServer) InitServer(port, addr, sqlPath string, maxlogs int) {

	PS.port = port
	PS.Adrr = addr

	PS.Hub = NewHub(maxlogs)
	go PS.Hub.Run()

	Service.InitPrinters()
	PS.Manager.InitManager(sqlPath)
	PS.mux = http.NewServeMux()
	//fs := http.FS(ServerFile)
	//PS.mux.Handle("/", http.FileServer(fs))
	PS.mux.HandleFunc("/", mainHandled)
	PS.mux.HandleFunc("/ws", PS.HandleWS)

	//
}

func CatchPanic(context string) {
	if r := recover(); r != nil {
		log.Printf("[PANIC in %s] %v", context, r)
	}
}

func (PS *CNCServer) Serve() {
	defer CatchPanic("main")

	go PS.UpdateLogs()
	go PS.UpdateTable()
	go PS.UpdateCNCData()
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
	ServeWs(PS.Hub, PS.WsCallBack, w, r)
}

func (PS *CNCServer) WsCallBack(client *Client, message []byte) {
	Responce := PS.ExecuteWSMessage(message)
	client.WriteMessage(Responce)
}

const PingTime = time.Second * 5

// func (PS *CNCServer) DataUpdater() {
// 	go func() {
// 		PS.TimerUpdader.WaitNewData()
// 		PS.Hub.Send([]byte(PS.TimerUpdader.GetNowData()),false)
// 	}()
// 	go func() {
// 		PS.StatusUpdader.WaitNewData()
// 		PS.Hub.Send([]byte(PS.StatusUpdader.GetNowData()),false)
// 	}()
// 	go func() {
// 		PS.LogsUpdader.WaitNewData()
// 		PS.Hub.Send([]byte(PS.LogsUpdader.GetNowData()))
// 	}()
// }

func (PS *CNCServer) ExecuteWSMessage(msg []byte) []byte {
	var mas WSMessage
	err := json.Unmarshal([]byte(msg), &mas)
	if err != nil {
		log.Println(err)
	}

	// fmt.Printf("mas: %v\n", mas)
	// log.Println("WebSocket message" + msg)
	switch mas.Type {
	case "connect":
		con := Service.ConnectionData{}
		json.Unmarshal(mas.Data, &con)
		err := PS.Manager.Connect(con)
		if err != nil {
			log.Printf("The device: %v, connected by %v, cant connected by reason: %v", con.ConnectionData, con.TypeOfConnection, err.Error())
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
		cncsJson := PS.Manager.GetJson()
		if cncsJson != "" && cncsJson != "[]" {
			cncsMsg := fmt.Sprintf(`{"type":"printers","data":%s}`, cncsJson)
			PS.Hub.Send([]byte(cncsMsg), false)
		}

		return WEB_Socket_ACK(mas.ReqId, true)

	case "GetRegistry":

		// cncsJson := PS.Manager.GetJson()
		// if cncsJson != "" && cncsJson != "[]" {
		// 	cncsMsg := fmt.Sprintf(`{"type":"printers","data":%s}`, cncsJson)
		// 	PS.Hub.Send([]byte(cncsMsg), false)
		// }

		reg := PS.Manager.GetRegistry()
		PS.Hub.Send(reg, false)
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

func (PS *CNCServer) UpdateCNCData() {
	for {
		// Send CNC machines data
		Data := PS.Manager.GetNewCncData()
		if len(Data) != 0 && string(Data) != "[]" {
			cncsMsg := fmt.Sprintf(`{"type":"printers","data":%s}`, Data)
			PS.Hub.Send([]byte(cncsMsg), false)
		}
	}
}

func (PS *CNCServer) UpdateTable() {
	for {
		// Send CNC machines data
		Data := PS.Manager.GetTable()
		if len(Data) != 0 && string(Data) != "[]" {
			cncsMsg := fmt.Sprintf(`{"type":"printers","data":%s}`, Data)
			PS.Hub.Send([]byte(cncsMsg), false)
		}
	}
}

func (PS *CNCServer) UpdateTimer() {
	for {
		JsonData := PS.Manager.GetTimeCharge()
		dataSec := struct {
			Uptime            string `json:"uptime"`
			ActiveConnections int    `json:"activeConnections"`
		}{Uptime: string(JsonData),
			ActiveConnections: int(PS.Hub.ActiveUsers)}
		jsonData, _ := json.Marshal(dataSec)

		// fmt.Printf("PS.Hub.ActiveUsers: %v\n", PS.Hub.ActiveUsers)
		resp := WSMessage{
			Type: "systemInfo",
			Data: jsonData,
		}
		res, _ := json.Marshal(resp)
		// fmt.Printf("res: %v\n", string(res))
		PS.Hub.Send(res, false)
	}
}

func (PS *CNCServer) UpdateStatus() {
	for {
		status := PS.Manager.GetStatus()
		// fmt.Printf("status: %v\n", string(status))
		PS.Hub.Send(status, false)
	}
}

func (PS *CNCServer) UpdateLogs() {
	id := uint32(1)
	for {
		Log := PS.Manager.GetLog()
		js := WEB_Socket_LOG(id, uint32(time.Now().Unix()), Log.Level, Log.Message)
		id++
		PS.Hub.Send(js, true)
	}
}
