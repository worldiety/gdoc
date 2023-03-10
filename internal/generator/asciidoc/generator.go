package asciidoc

import (
	"bytes"
	"fmt"
	"godocgenerator/internal/api"
	"godocgenerator/internal/parser/golang"
	"log"
	"reflect"
	"sort"
	"text/template"
)

const (
	constantsTemplate = "constants"
	variablesTemplate = "variables"
	functionsTemplate = "functions"
	moduleTemplate    = "module"
	packageTemplate   = "package"
)

// ExecuteTemplate uses a type switch to execute the correct template for all input items
func ExecuteTemplate(t *template.Template, items any, dest *bytes.Buffer) error {
	switch items.(type) {
	case map[string]*api.Constant:
		if err := t.ExecuteTemplate(dest, constantsTemplate, items); err != nil {

			return fmt.Errorf("unable to exec %v: %w", constantsTemplate, err)
		}
	case map[string]*api.Variable:
		if err := t.ExecuteTemplate(dest, variablesTemplate, items); err != nil {

			return fmt.Errorf("unable to exec %v: %w", variablesTemplate, err)
		}
	case map[string]*api.Function:
		if err := t.ExecuteTemplate(dest, functionsTemplate, items); err != nil {

			return fmt.Errorf("unable to execute %v: %w", functionsTemplate, err)
		}
	case *api.Module:
		if err := t.ExecuteTemplate(dest, moduleTemplate, items); err != nil {

			return fmt.Errorf("unable to execute %s: %w", moduleTemplate, err)
		}
	case *api.Package:
		if err := t.ExecuteTemplate(dest, packageTemplate, items); err != nil {

			return fmt.Errorf("unable to execute %s: %w", packageTemplate, err)
		}
	default:
		log.Fatalf("Type %v not supported", reflect.TypeOf(items))
	}
	return nil
}

var TemplatePattern = "assets/templates/*.tmpl"

// CreateModuleTemplate takes the parsed module, adds all its information to text templates and returns the outPut buffer
func CreateModuleTemplate(module *api.Module) (*bytes.Buffer, error) {
	mainTemplate, err := template.New("").Funcs(template.FuncMap{"lastKey": golang.LastKey}).ParseGlob(TemplatePattern)
	if err != nil {

		return nil, fmt.Errorf("failed to glob parse template files: %w", err)
	}

	var outPut bytes.Buffer
	if err = ExecuteTemplate(mainTemplate, module, &outPut); err != nil {

		return nil, fmt.Errorf("failed to execute template: %w", err)
	}
	sortedPackages := sortPackages(module.Packages)
	for _, p := range sortedPackages {
		if err = ExecuteTemplate(mainTemplate, p, &outPut); err != nil {

			return nil, fmt.Errorf("failed to execute template: %w", err)
		}
		data := []any{p.Consts, p.Vars, p.Functions}
		for _, items := range data {
			if err = ExecuteTemplate(mainTemplate, items, &outPut); err != nil {

				return nil, fmt.Errorf("failed to execute template: %w", err)
			}
		}
	}

	return &outPut, nil
}

func sortPackages(packages map[api.ImportPath]*api.Package) []*api.Package {
	return sortMapValues(packages, func(a, b *api.Package) bool {
		return a.Name < b.Name
	})
}

func sortMapValues[K comparable, V any](m map[K]V, less func(a, b V) bool) []V {
	tmp := make([]V, 0, len(m))
	for _, v := range m {
		tmp = append(tmp, v)
	}

	sort.Slice(tmp, func(i, j int) bool {
		return less(tmp[i], tmp[j])
	})

	return tmp
}
