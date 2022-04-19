package model

type ReqMsg struct {
	FuncName     string   `json:"func_name"`
	State        State    `json:"state"`
	Msg          string   `json:"msg"`
	Severity     Severity `json:"severity"`
	ResourceType string   `json:"resource_type"`
}

type State string
type Severity string

const (
	None       = State("NONE")
	Pending    = State("PENDING")
	InProgress = State("INPROGRESS")
	Faliure    = State("FAIILURE")
	Timeout    = State("TIMEOUT")
	Success    = State("SUCCESS")
	Unknown    = Severity("UNKNOWN")
	Warning    = Severity("WARNING")
	Error      = Severity("ERROR")
	Critical   = Severity("CRITICAL")
	Info       = Severity("INFO")
)
