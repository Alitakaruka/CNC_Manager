package Server

import (
	"encoding/json"
	"sync"

	"github.com/gorilla/websocket"
)

type WSTransmiterr struct {
	Data  string
	Flag  bool
	mutex sync.Mutex
	cond  *sync.Cond
}

func NewWSTransmiterr() *WSTransmiterr {
	ws := &WSTransmiterr{}
	ws.cond = sync.NewCond(&ws.mutex)
	return ws
}

func (ws *WSTransmiterr) WaitNewData() {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()

	old := ws.Data
	for old == ws.Data {
		ws.cond.Wait()
	}
}

func (ws *WSTransmiterr) SetNewData(data string) {
	ws.mutex.Lock()
	ws.Data = data
	ws.cond.Broadcast()
	ws.mutex.Unlock()
}

func (ws *WSTransmiterr) GetNowData() string {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()
	return ws.Data
}

type WSConnection struct {
	WS    websocket.Conn
	Mutex sync.Mutex
}

type WSMessage struct {
	Type  string          `json:"type"`
	ReqId string          `json:"reqId"`
	Data  json.RawMessage `json:"data"`
}

type jsonWsGcode struct {
	gcode     string
	uniqueKey string
	reqId     string
}

type WSError struct {
	Type  string `json:"type"`
	ReqId string `json:"reqId"`
	data  struct {
	} `json:"data"`
}

func WEB_Socket_ACK(reqId string, ok bool) []byte {
	responceACk := struct {
		Type  string `json:"type"`
		ReqId string `json:"reqId"`
		Data  struct {
			Ok bool `json:"ok"`
		} `json:"data"`
	}{Type: "ack", ReqId: reqId}
	responceACk.Data.Ok = ok
	js, err := json.Marshal(responceACk)
	if err != nil {
		return []byte{}
	}
	return js
}

func WEB_Socket_ERROR(reqId, msg string) []byte {

	responceError := struct {
		Type    string `json:"type"`
		ReqId   string `json:"reqId"`
		Message string `json:"message"`
	}{Type: "error", ReqId: reqId, Message: msg}
	js, err := json.Marshal(responceError)
	if err != nil {
		return []byte{}
	}
	return js
}

func WEB_Socket_LOG(id, Timestamp uint32, level, msg string) []byte {
	template := struct {
		Type string `json:"type"`
		Data struct {
			Id        uint32 `json:"id"`
			Timestamp uint32 `json:"timestamp"`
			Level     string `json:"level"`
			Message   string `json:"message"`
		} `json:"data"`
	}{Type: "log"}
	template.Data.Id = id
	template.Data.Level = level
	template.Data.Timestamp = Timestamp
	template.Data.Message = msg

	js, err := json.Marshal(template)
	if err != nil {
		return []byte{}
	}
	return js
}
