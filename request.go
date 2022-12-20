package web

type Request struct {
	Timestamp int64               `json:"timestamp"`
	FromIP    string              `json:"from_ip"`
	Method    string              `json:"method"`
	Host      string              `json:"host"`
	Path      string              `json:"path"`
	Query     map[string][]string `json:"query"`
	Headers   map[string][]string `json:"headers"`
	Body      []byte              `json:"body"`
}
