package web

type Type struct {
	IsPointer bool `json:"is_pointer"`
	IsArray   bool `json:"is_array"`
	IsMap     bool `json:"is_map"`
	BaseType  ID   `json:"base_type"`
}
