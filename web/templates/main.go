package templates

import (
	"embed"
	"net/http"
	"log"
	"html/template"
)

//go:embed *
var TemplateFiles embed.FS

func Render(rw http.ResponseWriter, templateString string, context any) error {
	t, err := template.New(templateString).ParseFS(TemplateFiles, templateString, "layout.html.tmpl")
	if err != nil {
		log.Fatal(err)
	}
	return t.Execute(rw, context)
}
