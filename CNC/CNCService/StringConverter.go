package CNCService

import (
	"log"
	"strconv"
	"sync"
)

func SetIntValue(field *int, val string, mut *sync.Mutex) {
	mut.Lock()
	parsed, err := strconv.Atoi(val)
	if err != nil {
		log.Printf("Error conversion:%v", err)
	}
	*field = parsed
	mut.Unlock()
}

func SetFloatValue(field *float32, val string, mut *sync.Mutex) {
	mut.Lock()
	parsed, err := strconv.ParseFloat(val, 32)
	if err != nil {
		log.Printf("Error conversion:%v", err)
	}
	*field = float32(parsed)
	mut.Unlock()
}

func SetUintValue(field *uint, val string, mut *sync.Mutex) {
	mut.Lock()
	parsed, err := strconv.Atoi(val)
	if err != nil {
		log.Printf("Error conversion:%v", err)
	}
	*field = uint(parsed)
	mut.Unlock()
}

func SetStringValue(field *string, val string, mut *sync.Mutex) {
	mut.Lock()
	*field = val
	mut.Unlock()
}
