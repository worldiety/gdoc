package asciidoc

import (
	"embed"
	"fmt"
	"text/template"
)

//go:embed *.tmpl
var templateFiles embed.FS
var Templates *template.Template

func init() {
	tpl, err := template.ParseFS(templateFiles, "*.tmpl")
	if err != nil {
		panic(fmt.Errorf("cannot parse embedded templates: %w", err))
	}

	Templates = tpl
}
