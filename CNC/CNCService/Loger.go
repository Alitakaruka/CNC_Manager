package CNCService

const MaxLogsDefult = 100

const (
	LogLevelError       = "error"
	LogLevelWarning     = "warning"
	LogLevelInformation = "info"
	LogLevelSuccess     = "success"
)

type Log struct {
	Code    int
	Level   string
	Message string
}
