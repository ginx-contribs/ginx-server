package templates

import (
	"embed"
	"html/template"
)

//go:embed *.tmpl
var EmbedFs embed.FS

func ParseEmbedTemplates() (templates *template.Template, err error) {
	return template.ParseFS(EmbedFs, "*.tmpl")
}
