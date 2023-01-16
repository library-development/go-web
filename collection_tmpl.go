package web

import "html/template"

var collectionTmpl = template.Must(template.New("collection").Parse(collectionHTML))
