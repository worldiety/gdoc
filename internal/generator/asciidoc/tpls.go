package asciidoc

import (
	"embed"
	"fmt"
	"github.com/worldiety/gdoc/internal/parser/golang"
	"text/template"
)

//go:embed templates
var templateFiles embed.FS

// Templates is globally available and binds all template files in the asciidoc/templates directory.
// All available templates have to be called by name, to use them.
var Templates *template.Template // test comment
var TestVar *golang.AVariable

func init() {
	tpl, err := template.ParseFS(templateFiles, "templates/*.tmpl")
	if err != nil {
		panic(fmt.Errorf("cannot parse embedded templates: %w", err))
	}

	Templates = tpl
}
