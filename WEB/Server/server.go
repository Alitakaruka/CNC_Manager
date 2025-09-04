package Server

import (
	Service "CNCManager/Service"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
	_ "modernc.org/sqlite"
)

//go:embed files/index.html
var ServerFile embed.FS

type PrinterServer struct {
	Manager     Service.CNCManagerr
	mux         *http.ServeMux
	port        string
	Adrr        string
	PrinterRepo Service.PrinterRepository
}

func (PS *PrinterServer) InitServer() {
	//Config
	config := Service.GetConfig("config.yaml")
	var port, addr, sqlPath, logerPath string
	if config == nil {
		port = "8080"
		addr = "0.0.0.0"
		sqlPath = "CNCManagerDB.db"
		logerPath = "Logs.log"
	} else {
		port = config.Server.Port
		addr = config.Server.Addr
		sqlPath = config.Database.Path
		logerPath = config.Logs.Path
	}
	PS.port = port
	PS.Adrr = addr
	//
	startLoger(logerPath)
	Service.InitPrinters()
	PS.PrinterRepo.InitRepository(sqlPath)
	PS.mux = http.NewServeMux()
	fs := http.FS(ServerFile)
	PS.mux.Handle("/", http.FileServer(fs))
	//PS.mux.HandleFunc("/", mainHandled)
	PS.mux.HandleFunc("/connect", PS.ConnectPrinter)
	PS.mux.HandleFunc("/ws", PS.HandleWS)
	PS.mux.HandleFunc("/api/Printers", PS.GetPrintersInformation)
	PS.mux.HandleFunc("/api/ExecuteTask", PS.ExecuteTask)
	PS.mux.HandleFunc("/api/SaveSettings", PS.saveSettings)
	PS.mux.HandleFunc("/api/GetSettings", PS.getSettings)
	PS.mux.HandleFunc("/api/sendGCode", PS.SendGCode)
	// PS.mux.HandleFunc("/api/SetColor", PS.SetColor)
}

func CatchPanic(context string) {
	if r := recover(); r != nil {
		log.Printf("[PANIC in %s] %v", context, r)
	}
}

func (PS *PrinterServer) Serve() {
	defer CatchPanic("main")
	fmt.Printf("Server started: %s\n", PS.Adrr+":"+PS.port)
	go PS.Manager.LoggingAsync()
	err := http.ListenAndServe(PS.Adrr+":"+PS.port, PS.mux)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

// HTTP +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (PS *PrinterServer) ConnectPrinter(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "Failed to connect the printer due to: "+ex.Error(), http.StatusBadRequest)
		log.Println("Failed to connect the printer due to: " + ex.Error())
		return
	}
	w.Write([]byte("ok"))
}

func (PS *PrinterServer) SendGCode(w http.ResponseWriter, r *http.Request) {
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

// func (PS *PrinterServer) SetColor(w http.ResponseWriter, r *http.Request) {
// 	Red, err := strconv.Atoi(r.URL.Query().Get("R"))
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}
// 	Green, err := strconv.Atoi(r.URL.Query().Get("G"))
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}
// 	Blue, err := strconv.Atoi(r.URL.Query().Get("B"))
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}
// 	uniqueKey := r.URL.Query().Get("uniqueKey")

// 	if uniqueKey == "" {
// 		http.Error(w, "Unique key can not be empty", http.StatusBadRequest)
// 		return
// 	}
// 	err = PS.Manager.SetColor(byte(Red), byte(Green), byte(Blue), uniqueKey)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}
// 	w.Write([]byte("ok"))

// }

func (PS *PrinterServer) ExecuteTask(w http.ResponseWriter, r *http.Request) {
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
	PS.Manager.ExecuteTask(Key, fileBytes)
	w.Write([]byte("Start printing!"))
}

func (PS *PrinterServer) GetPrintersInformation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(PS.Manager.GetJson()))
}

func mainHandled(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "WEB/Server/files/index.html")
}

func (PS *PrinterServer) saveSettings(w http.ResponseWriter, r *http.Request) {

}

func (PS *PrinterServer) getSettings(w http.ResponseWriter, r *http.Request) {

}

func startLoger(filePath string) {
	file, ex := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if ex != nil {
		log.Fatal(ex)
	}
	muliWriter := io.MultiWriter(os.Stdout, file)
	log.SetFlags(log.Ltime | log.Ldate | log.Llongfile)
	log.SetOutput(muliWriter)
}

//HTTP ----------------------------------------------------------------------

// WEB SOCKET ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
func (PS *PrinterServer) HandleWS(w http.ResponseWriter, r *http.Request) {
	log.Println("connecting...")
	upgrader := websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
		return true
	}, EnableCompression: true}
	WebSoc, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	WebSoc.SetPingHandler(func(appdata string) error {
		return WebSoc.WriteControl(websocket.PongMessage, []byte(appdata), time.Now().Add(5))
	})
	WebSoc.SetPongHandler(func(appdata string) error {
		log.Println("Pong")
		if appdata != "Ping" {
			return errors.New("appdata error")
		}
		return nil
	})

	isClose := make(chan struct{})
	WebSoc.SetCloseHandler(func(code int, text string) error {
		close(isClose)
		return nil
	})
	go PS.WsRead(WebSoc, isClose)
	go PS.WSWriteData(WebSoc, isClose)
	go log.Printf("New socket!")
}

// TextMessage = 1
// BinaryMessage = 2
const PingTime = time.Second * 5

func (PS *PrinterServer) WsRead(WS *websocket.Conn, isClose chan struct{}) {
	ticker := time.NewTicker(PingTime)

	go func() {
		for {
			select {
			case <-ticker.C:
				WS.WriteControl(websocket.PingMessage, []byte("Ping"), time.Now().Add(PingTime))
				ticker.Reset(time.Second * 5)
			case <-isClose:
				return
			}
		}
	}()
	for {
		msgType, msg, err := WS.ReadMessage()
		ticker.Reset(PingTime)
		if err != nil {
			log.Println(err)
			return
		}
		_ = msgType
		if msgType == websocket.TextMessage {
			ParceWSMessage(string(msg))
		}
	}
}

func (PS *PrinterServer) WSWriteData(WS *websocket.Conn, isClose chan struct{}) {

}

func ParceWSMessage(msg string) {

}
