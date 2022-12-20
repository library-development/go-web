package web

type ExecRequest struct {
	Token  string         `json:"token"`
	Pkg    string         `json:"pkg"`
	Func   string         `json:"func"`
	Inputs map[string]any `json:"inputs"`
}
