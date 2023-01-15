package web

type Event struct {
	Type    string `json:"type"`
	Payload []byte `json:"payload"`
}
