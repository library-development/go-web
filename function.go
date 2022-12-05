package web

type Function struct {
	Name     string    `json:"name"`
	Inputs   []Field   `json:"inputs"`
	Outputs  []Field   `json:"outputs"`
	Validate []Command `json:"validate"`
	Execute  []Command `json:"execute"`
}
