package golang

import (
	"fmt"
	"github.com/worldiety/gdoc/internal/api"
	"golang.org/x/tools/go/packages"
	"strings"
)

type loadedPackages struct {
	pkgs      map[string]*packages.Package
	fnMap     map[string]*api.Function
	structMap map[string]*api.Struct
}

func newLoadedPackages() *loadedPackages {
	return &loadedPackages{
		pkgs:  map[string]*packages.Package{},
		fnMap: map[string]*api.Function{},
	}
}

func Resolve(m *api.Module) error {
	lp := newLoadedPackages()
	var err error

	for dir, _ := range m.Packages {
		err = lp.loadPackages(dir)
		if err != nil {
			return err
		}
	}

	addFieldInformation(m, lp)

	return nil
}

func (lp *loadedPackages) loadPackages(dir string) error {
	pkgs, err := packages.Load(
		&packages.Config{Mode: packages.NeedName | packages.NeedTypes | packages.NeedModule, Tests: false}, dir)
	if err != nil {
		return fmt.Errorf("could not load packages from %s: %w", dir, err)
	}

	for _, pkg := range pkgs {
		lp.pkgs[pkg.Name] = pkg
	}

	return nil
}

// add information to the module and all it's sub-parts that the ast package does not provide, but the packages.Package does
func addFieldInformation(m *api.Module, lp *loadedPackages) {
	for path, p := range m.Packages {
		if p.Name != "app" {
			continue
		}
		p.PackageDefinition = api.NewRefID(path, p.Name)
		for id, variable := range p.Vars {
			variable.TypeDefinition = api.NewRefID(path, id)
			p.Types[variable.Name] = variable.TypeDefinition
		}
		for id, constant := range p.Consts {
			constant.TypeDefinition = api.NewRefID(path, id)
			p.Types[constant.Name] = constant.TypeDefinition
		}
		for id, function := range p.Functions {
			function.TypeDefinition = api.NewRefID(path, id)
			p.Types[function.Name] = function.TypeDefinition
		}
		for id, s := range p.Structs {
			s.TypeDefinition = api.NewRefID(path, id)
			p.Types[s.Name] = s.TypeDefinition

			for _, f := range s.Fields {
				var importPath string
				var identifier string
				if strings.Contains(f.SrcTypeDefinition, ".") {
					parts := strings.Split(f.SrcTypeDefinition, ".")
					if pack := lp.pkgs[strings.Replace(parts[0], "*", "", -1)]; pack != nil {
						importPath = pack.PkgPath
						identifier = parts[1]
					}
				} else if lp.pkgs[p.Name].Types.Scope().Lookup(f.SrcTypeDefinition) != nil {
					importPath = lp.pkgs[p.Name].PkgPath
					identifier = f.SrcTypeDefinition
				}

				f.TypeDefinition = api.NewRefID(importPath, identifier)

				// if the Field f is of a type from a package, contained in the project, it should be linked, otherwise it should not.
				linkType(f, lp)
			}
		}
	}
}

func linkType(f *api.Field, lp *loadedPackages) {
	if lp.pkgs[f.TypeDefinition.PackageName()] != nil && lp.pkgs[f.TypeDefinition.PackageName()].Types.Scope().
		Lookup(strings.Replace(f.TypeDefinition.Identifier, "*", "", -1)) != nil {
		f.Link = true
	}
}
