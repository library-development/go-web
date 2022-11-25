package web

type API struct {
	Queries  map[string]*Query
	Commands map[string]*Command
}

type Query struct {
	InputTypes []*Field
	OutputType *Field
	Execute    *Function
}

type Command struct {
	InputTypes []*Field
	OutputType *Field
	Execute    *Function
}

type Function func(inputs map[string]string) (any, error)

type Field struct {
	Name     string
	Type     string
	Required bool
}
