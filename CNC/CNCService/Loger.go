package CNCService

const MaxLogsDefult = 100

type Loger struct {
	Data    []string
	MaxLogs int
}

func NewLoger(MaxLogs int) *Loger {
	if MaxLogs == 0 {
		MaxLogs = MaxLogsDefult
	}
	return &Loger{MaxLogs: MaxLogs, Data: make([]string, MaxLogs)}
}
func (LG *Loger) NewLog(log string) {
	LG.Data = append(LG.Data, log)
	LG.ClearLogs()
}

func (LG *Loger) ClearLogs() {
	Gap := len(LG.Data) - LG.MaxLogs
	if Gap > 0 {
		LG.Data = LG.Data[:len(LG.Data)-Gap]
	}
}

func (LG *Loger) GetLogs() []string {
	return LG.Data
}
