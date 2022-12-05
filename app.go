package web

type App struct {
	Models []Model `json:"models"`
}

type Model struct {
	Name    string     `json:"name"`
	Fields  []Field    `json:"fields"`
	Methods []Function `json:"methods"`
}

type Function struct {
	Name     string    `json:"name"`
	Inputs   []Field   `json:"inputs"`
	Outputs  []Field   `json:"outputs"`
	Validate []Command `json:"validate"`
	Execute  []Command `json:"execute"`
}

type Command struct {
	Assignments []string `json:"assignments"`
	Func        ID       `json:"func"`
	Args        []string `json:"args"`
}

type Field struct {
	Name string `json:"name"`
	Type Type   `json:"type"`
}

type Type struct {
	IsPointer bool `json:"is_pointer"`
	IsArray   bool `json:"is_array"`
	IsMap     bool `json:"is_map"`
	BaseType  ID   `json:"base_type"`
}

type ID struct {
	Path string `json:"path"`
}
