package web

type Command struct {
	Assignments []string `json:"assignments"`
	Func        ID       `json:"func"`
	Args        []string `json:"args"`
}
