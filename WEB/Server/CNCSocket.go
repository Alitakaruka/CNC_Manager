package Server

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	hub *Hub
	mut sync.Mutex
	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

func (C *Client) WriteMessage(message []byte) {
	C.mut.Lock()
	C.conn.WriteMessage(websocket.TextMessage, message)
	C.mut.Unlock()
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) ReadPump(callBack func(client *Client, message []byte)) {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	// c.conn.SetReadLimit(1024)
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(60 * time.Second)); return nil })

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, []byte{'\n'}, []byte{' '}, -1))
		callBack(c, message)
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(time.Second * 10)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(time.Second * 10))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.WriteMessage(message)
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(time.Second * 10))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
func ServeWs(hub *Hub, ReadCallBackFunc func(client *Client, message []byte), w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 1024)}
	client.hub.register <- client

	go client.WritePump()
	go client.ReadPump(ReadCallBackFunc)

	log.Println("New connection accepted sucsessfuly!")
}

type Hub struct {

	//Current active users
	ActiveUsers int
	//Byffer for new or reconnected clients
	ramBuffer [][]byte
	// Registered clients.
	clients map[*Client]bool
	// Inbound messages from the clients.
	broadcast chan []byte
	// Register requests from the clients.
	register chan *Client
	// Unregister requests from clients.
	unregister chan *Client
	memlen     int
}

func NewHub(MemoryLength int) *Hub {

	return &Hub{
		memlen:     MemoryLength,
		ramBuffer:  make([][]byte, 0),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			h.ActiveUsers++
			for _, val := range h.ramBuffer {
				select {
				case client.send <- val:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					h.ActiveUsers--
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

func (h *Hub) Send(msg []byte, saveInMemory bool) {
	if saveInMemory {
		h.ramBuffer = append(h.ramBuffer, msg)
		if len(h.ramBuffer) > h.memlen { //save last 100 massages
			h.ramBuffer = h.ramBuffer[1:]
		}
	}
	h.broadcast <- msg
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

type WSError struct {
	Type  string `json:"type"`
	ReqId string `json:"reqId"`
	Data  struct {
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
