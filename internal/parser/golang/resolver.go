package golang

import (
	"fmt"
	"github.com/worldiety/gdoc/internal/api"
	"golang.org/x/tools/go/packages"
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
		p.PackageDefinition = api.NewRefID(path, p.Name)
		for id, variable := range p.Vars {
			variable.TypeDesc.TypeDefinition = api.NewRefID(path, id)
			p.Types[variable.Name] = variable.TypeDesc.TypeDefinition
		}
		for id, constant := range p.Consts {
			constant.TypeDesc.TypeDefinition = api.NewRefID(path, id)
			p.Types[constant.Name] = constant.TypeDesc.TypeDefinition
		}
		for id, function := range p.Functions {
			function.TypeDefinition = api.NewRefID(path, id)
			p.Types[function.Name] = function.TypeDefinition
			handleFields(function.Parameters, p, lp)
			handleFields(function.Results, p, lp)
		}
		for id, s := range p.Structs {
			s.TypeDefinition = api.NewRefID(path, id)
			p.Types[s.Name] = s.TypeDefinition

			for _, f := range s.Fields {
				handleField(f, p, lp)
			}
		}
	}
}

func handleField(f *api.Field, p *api.Package, lp *loadedPackages) {
	handleType(f, p, lp)
}

func handleType(f *api.Field, p *api.Package, lp *loadedPackages) {

	if f.TypeDesc.Map() {
		handleMapType(f, p.Name, lp)
	} else {
		typeDescInfo(p.Name, f.TypeDesc, lp)
	}
}

func handleBasicType(f *api.Field, pName string, lp *loadedPackages) {
	typeDescInfo(pName, f.TypeDesc, lp)
}

func handleFields(parameters map[string]*api.Field, currentPackage *api.Package, lp *loadedPackages) {
	for _, p := range parameters {
		handleField(p, currentPackage, lp)
	}
}

func handleMapType(f *api.Field, pName string, lp *loadedPackages) {
	keyTypeDef, valueTypeDef := f.TypeDesc.MapSrcDefs()
	f.MapType = &api.MapType{}
	f.MapType.ValueType = &api.TypeDesc{SrcTypeDefinition: valueTypeDef}
	f.MapType.KeyType = &api.TypeDesc{SrcTypeDefinition: keyTypeDef}
	typeDescInfo(pName, f.MapType.KeyType, lp)
	typeDescInfo(pName, f.MapType.ValueType, lp)
}

func typeDescInfo(pName string, td *api.TypeDesc, lp *loadedPackages) {
	var importPath string
	var link bool
	var origin api.TypeOrigin

	// from current package or built-in
	if include, _ := td.IncludesPkg(); !include {
		if p := lp.pkgs[pName]; p != nil {
			if t := p.Types.Scope().Lookup(td.Identifier()); t != nil {
				importPath = p.PkgPath
				link = true
				origin = api.LocalCustom
			}
		}
	} else {
		// from other package
		if p := lp.pkgs[td.PkgName()]; p != nil {
			importPath = p.PkgPath
			link = true
			origin = api.ExternalCustom
		} else {
			origin = api.ExternalNonCustom
		}
	}
	td.TypeDefinition = api.NewRefID(importPath, td.Identifier())
	td.Link = link
	td.TypeOrigin = origin
}
