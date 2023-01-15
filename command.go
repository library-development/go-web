package web

type Command struct {
	Assignments []string `json:"assignments"`
	Func        string   `json:"func"`
	Args        []string `json:"args"`
}
