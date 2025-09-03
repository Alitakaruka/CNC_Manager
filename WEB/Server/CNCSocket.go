package Server

import (
	"crypto/sha1"
	"encoding/base64"
	"log"
	"net"
	"net/http"
)

const (
	MessageTypeText   = 1
	MessageTypeBinaty = 2
	MessageTypeClose  = 8
	MessageTypePing   = 9
	MessageTypePong   = 10
)

type CNCSocket struct {
	net.Conn
}

func (CNC_Soc *CNCSocket) ReadMessage() (string, error) {
	buffer := make([]byte, 512)
	n, err := CNC_Soc.Read(buffer)
	if err != nil {
		return "", err
	}
	if n != 0 {

	}
	return "", nil
}

func (CNC_Soc *CNCSocket) decodeFrame(frame []byte) string {
	payloadLen := int(frame[1] & 127)     // 131 & 127 = 3 (длина данных)
	maskKey := frame[2:6]                 // 4 байта маски
	maskedData := frame[6 : 6+payloadLen] // зашифрованные данные

	// Снимаем маску
	decoded := make([]byte, payloadLen)
	for i := 0; i < payloadLen; i++ {
		decoded[i] = maskedData[i] ^ maskKey[i%4]
	}
	return string(decoded)
}
func computeAcceptKey(key string) string {
	const magic = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
	h := sha1.New()
	h.Write([]byte(key + magic))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func HandleWS(w http.ResponseWriter, r *http.Request) *CNCSocket {
	CNC_Sock := CNCSocket{}
	if r.Header.Get("Upgrade") != "websocket" {
		return nil
	}
	key := r.Header.Get("Sec-WebSocket-Key")
	accept := computeAcceptKey(key)
	w.Header().Set("Upgrade", "websocket")
	w.Header().Set("Connection", "Upgrade")
	w.Header().Set("Sec-WebSocket-Accept", accept)
	w.WriteHeader(http.StatusSwitchingProtocols)
	hj, _ := w.(http.Hijacker)
	conn, _, _ := hj.Hijack()
	log.Println("New Socket")
	CNC_Sock.Conn = conn
	return &CNC_Sock
}
