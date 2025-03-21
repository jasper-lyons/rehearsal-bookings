package templates

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"os"
)

//go:embed *.html.tmpl
var TemplateFiles embed.FS

func Render(rw http.ResponseWriter, templateString string, context any) error {
	t, err := template.
		New(filepath.Base(templateString)).
		Funcs(template.FuncMap {
			"Getenv": func (name string) string { return os.Getenv(name) },
		}).
		ParseFS(TemplateFiles, templateString, "layout.html.tmpl")
	if err != nil {
		log.Fatal(err)
	}
	return t.Execute(rw, context)
}
