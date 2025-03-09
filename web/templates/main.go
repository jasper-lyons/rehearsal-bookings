package templates

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

//go:embed *.html.tmpl
var TemplateFiles embed.FS

func Render(rw http.ResponseWriter, templateString string, context any) error {
	t, err := template.New(filepath.Base(templateString)).ParseFS(TemplateFiles, templateString, "layout.html.tmpl")
	if err != nil {
		log.Fatal(err)
	}
	return t.Execute(rw, context)
}
