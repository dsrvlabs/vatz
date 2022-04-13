package message

type ReqMsg struct {
	FuncName     string `json:"func_name"`
	State        string `json:"state"`
	Msg          string `json:"msg"`
	Severity     string `json:"severity"`
	ResourceType string `json:"resource_type"`
}
