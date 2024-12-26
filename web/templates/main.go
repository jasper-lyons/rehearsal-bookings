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
	bytes, err := TemplateFiles.ReadFile(templateString)
	if err != nil {
		log.Fatal(err)
	}
	t, err := template.New(templateString).Parse(string(bytes))
	if err != nil {
		log.Fatal(err)
	}
	return t.Execute(rw, context)
}
