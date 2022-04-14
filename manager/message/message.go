package message

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
	Failure  = State("FAILURE")
	Success  = State("SUCCESS")
	Critical = Severity("CRITICAL")
	Info     = Severity("INFO")
)
