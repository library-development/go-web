package web

type ResourceList interface {
	Search(str string) SearchResults
	Filter(func(Resource) bool) ResourceList
}
