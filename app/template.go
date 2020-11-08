package app

import (
	"html/template"
	"path"
)

const baseTemplateName = "views/base.html"

func Page(filename string) *template.Template {
	return template.Must(template.ParseFiles(path.Join("views", filename), baseTemplateName))
}
