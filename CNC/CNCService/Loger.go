package CNCService

const MaxLogsDefult = 100

const (
	LogLevelError       = "error"
	LogLevelWarning     = "warning"
	LogLevelInformation = "info"
)

type Log struct {
	Code    int
	Level   string
	Message string
}
