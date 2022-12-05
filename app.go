package web

type App struct {
	Models map[string]Model `json:"models"`
}

type Model struct {
	Fields  []Field    `json:"fields"`
	Methods []Function `json:"methods"`
}

type Function struct {
	Name    string    `json:"name"`
	Inputs  []Field   `json:"inputs"`
	Outputs []Field   `json:"outputs"`
	Impl    []Command `json:"impl"`
}

type Command struct {
	Assignments []string `json:"assignments"`
	Func        ID       `json:"func"`
	Args        []string `json:"args"`
}

type Field struct {
	Type Type `json:"type"`
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
