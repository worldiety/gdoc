package golang

import (
	"fmt"
	"go/types"
	"godocgenerator/internal/api"
	"golang.org/x/tools/go/packages"
)

var loadedPackages = map[string]*packages.Package{}
var fnMap = map[string]*api.Function{}

func Resolve(m *api.Module) error {
	for dir, _ := range m.Packages {
		err := loadPackages(dir)
		if err != nil {
			return err
		}
	}

	fnSignatures(m)

	return nil
}

func loadPackages(dir string) error {
	pkgs, err := packages.Load(&packages.Config{
		Mode:  packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles | packages.NeedImports | packages.NeedDeps | packages.NeedExportFile | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedTypesSizes | packages.NeedModule | packages.NeedEmbedFiles | packages.NeedEmbedPatterns,
		Tests: true,
	}, dir)
	if err != nil {
		return fmt.Errorf("could not load packages from %s: %w", dir, err)
	}
	for _, pkg := range pkgs {
		loadedPackages[pkg.Name] = pkg
	}
	return nil
}

func fnSignatures(m *api.Module) {
	for _, p := range m.Packages {
		for _, fn := range p.Functions {
			fnMap[fn.Name] = fn
		}
	}

	for _, lp := range loadedPackages {
		for n, function := range fnMap {
			if fn := lp.Types.Scope().Lookup(n); fn != nil {
				function.Signature = addNameToSignature(fn)
			}
		}
	}
}

func addNameToSignature(fn types.Object) string {
	sigString := fn.Type().String()
	name := fn.Name()
	return sigString[:4] + " " + name + sigString[4:]
}
