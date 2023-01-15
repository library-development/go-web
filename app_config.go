package web

type AppConfig struct {
	StaticFiles map[string][]byte
	RootType    Type
	Schema      struct {
		User   Schema
		Org    Schema
		Public Schema
	}
	Views map[string]*View
}
