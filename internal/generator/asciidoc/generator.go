package asciidoc

import (
	"bytes"
	"fmt"
	"github.com/worldiety/gdoc/internal/api"
	"github.com/worldiety/gdoc/internal/parser/golang"
	"log"
	"reflect"
	"sort"
	"text/template"
)

type Generator struct {
	module golang.AModule
}

const (
	headerTemplate    = "header"
	constantsTemplate = "constants"
	variablesTemplate = "variables"
	functionsTemplate = "functions"
	moduleTemplate    = "module"
	packageTemplate   = "package"
	structTemplate    = "structs"
)

// executeTemplate uses a type switch to execute the correct template for all input items
func executeTemplate(t *template.Template, items any, dest *bytes.Buffer) error {
	switch items.(type) {
	case golang.AFunctions:
		if err := t.ExecuteTemplate(dest, functionsTemplate, items); err != nil {

			return fmt.Errorf("unable to execute %v: %w", functionsTemplate, err)
		}
	case golang.AModule:
		if err := t.ExecuteTemplate(dest, moduleTemplate, items); err != nil {

			return fmt.Errorf("unable to execute %s: %w", moduleTemplate, err)
		}
	case golang.APackage:
		if err := t.ExecuteTemplate(dest, packageTemplate, items); err != nil {

			return fmt.Errorf("unable to execute %s: %w", packageTemplate, err)
		}
	case golang.AStructs:
		if err := t.ExecuteTemplate(dest, structTemplate, items); err != nil {
			return fmt.Errorf("unable to execute %s: %w", structTemplate, err)
		}
	case golang.AVariables:
		if err := t.ExecuteTemplate(dest, variablesTemplate, items); err != nil {
			return fmt.Errorf("unable to execute %s: %w", variablesTemplate, err)
		}
	case golang.AsciiDocHeader:
		if err := t.ExecuteTemplate(dest, headerTemplate, items); err != nil {
			return fmt.Errorf("unable to execute %s: %w", headerTemplate, err)
		}
	default:
		log.Fatalf("Type %v not supported", reflect.TypeOf(items))
	}
	return nil
}

// CreateModuleTemplate takes the parsed module, adds all its information to text templates and returns the outPut buffer
func CreateModuleTemplate(module golang.AModule) (*bytes.Buffer, error) {
	var outPut bytes.Buffer

	if err := executeTemplate(Templates, golang.NewAsciiDocHeader(), &outPut); err != nil {

		return nil, fmt.Errorf("failed to execute index template: %w", err)
	}
	if err := executeTemplate(Templates, module, &outPut); err != nil {

		return nil, fmt.Errorf("failed to execute template: %w", err)
	}
	sortedPackages := sortPackages(module.Packages)
	for _, p := range sortedPackages {
		if err := executeTemplate(Templates, p, &outPut); err != nil {

			return nil, fmt.Errorf("failed to execute template: %w", err)
		}

		data := []any{golang.NewAVariables(p.Vars) /*, p.Consts*/, golang.NewAStructs(p.Structs), golang.NewAFunctions(p.Functions)}
		for _, items := range data {
			if err := executeTemplate(Templates, items, &outPut); err != nil {

				return nil, fmt.Errorf("failed to execute template: %w", err)
			}
		}
	}

	return &outPut, nil
}

func sortPackages(packages map[api.ImportPath]golang.APackage) []golang.APackage {
	return sortMapValues(packages, func(a, b golang.APackage) bool {
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
