package templates

import (
	"embed"
	"net/http"
	"log"
	"html/template"
	"path/filepath"
)

//go:embed *.html.tmpl admin/*.html.tmpl
var TemplateFiles embed.FS

func Render(rw http.ResponseWriter, templateString string, context any) error {
	t, err := template.New(filepath.Base(templateString)).ParseFS(TemplateFiles, templateString, "layout.html.tmpl")
	if err != nil {
		log.Fatal(err)
	}
	return t.Execute(rw, context)
}
