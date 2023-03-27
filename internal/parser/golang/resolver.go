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
				if strings.Contains(f.SrcTypeDefinition, "map[") {
					handleMapType(f, p.Name, lp)
				} else {
					importPath, identifier, f.Link = typeDescInfo(p.Name, f.SrcTypeDefinition, lp)
				}

				f.TypeDefinition = api.NewRefID(importPath, identifier)

				//// if the Field f is of a type from a package, contained in the project, it should be linked, otherwise it should not.
				//linkType(f, lp)
			}
		}
	}
}

func handleMapType(f *api.Field, pName string, lp *loadedPackages) {
	var keyTypeDef, valueTypeDef string
	tmp := strings.Replace(f.SrcTypeDefinition, "map[", "", 1)
	tmpArr := strings.Split(tmp, "]")
	keyTypeDef = tmpArr[0]
	valueTypeDef = tmpArr[1]

	var importPath, identifier string
	var link bool
	importPath, identifier, link = typeDescInfo(pName, keyTypeDef, lp)
	f.MapType = &api.MapType{}
	f.MapType.KeyType = &api.TypeDesc{
		TypeDefinition:    api.NewRefID(importPath, identifier),
		SrcTypeDefinition: keyTypeDef,
		Link:              link,
	}

	importPath, identifier, link = typeDescInfo(pName, valueTypeDef, lp)
	f.MapType.ValueType = &api.TypeDesc{
		TypeDefinition:    api.NewRefID(importPath, identifier),
		SrcTypeDefinition: valueTypeDef,
		Link:              link,
	}
}

//func linkType(f *api.Field, lp *loadedPackages) {
//
//	if strings.Contains(f.SrcTypeDefinition, "map[") {
//		f.Link = handleMapLinks(f, lp)
//	}
//	if strings.Contains(f.SrcTypeDefinition, "[]") {
//		f.Link = handleArrayLink(f, lp)
//	}
//	if lp.pkgs[f.TypeDefinition.PackageName()] == nil {
//		f.Link = false
//	} else {
//		f.Link = handleField(f.TypeDefinition.PackageName(), f.TypeDefinition.Identifier, *lp)
//	}
//}
//
//func handleField(pName, fId string, lp loadedPackages) api.Link {
//	if lp.pkgs[pName].Types.Scope().
//		Lookup(withoutAsterix(fId)) != nil {
//		return api.FieldType
//	}
//	return api.None
//}
//
//func handleArrayLink(f *api.Field, lp *loadedPackages) api.Link {
//	var link api.Link
//	return link
//}
//
//func handleMapLinks(f *api.Field, lp *loadedPackages) api.Link {
//	var keyType, valueType string
//	pName := f.TypeDefinition.PackageName()
//	tmp := strings.Replace(withoutAsterix(f.SrcTypeDefinition), "map[", "", 1)
//	tmpArr := strings.Split(tmp, "]")
//	keyType = tmpArr[0]
//	valueType = tmpArr[1]
//
//	if handleField(pName, keyType, *lp) == api.None && handleField(pName, valueType, *lp) == api.None {
//		return api.None
//	}
//	if handleField(pName, keyType, *lp) == api.FieldType && handleField(pName, valueType, *lp) == api.None {
//		return api.MapKey
//	}
//	if handleField(pName, keyType, *lp) == api.None && handleField(pName, valueType, *lp) == api.FieldType {
//		return api.MapValue
//	}
//	if handleField(pName, keyType, *lp) == api.FieldType && handleField(pName, valueType, *lp) == api.FieldType {
//		return api.MapAll
//	}
//	return api.None
//}

func withoutAsterix(s string) string {
	return strings.Replace(s, "*", "", -1)
}

func typeDescInfo(pName, srcDef string, lp *loadedPackages) (string, string, bool) {
	var importPath, identifier string
	var link bool
	// from current package
	if !strings.Contains(srcDef, ".") {
		if p := lp.pkgs[pName]; p != nil {
			if p.Types.Scope().Lookup(withoutAsterix(srcDef)) != nil {
				importPath = p.PkgPath
				link = true
			}
		}
		identifier = withoutAsterix(srcDef)
		return importPath, identifier, link
	}

	// from other package
	parts := strings.Split(srcDef, ".")
	if p := lp.pkgs[strings.Replace(parts[0], "*", "", -1)]; p != nil {
		importPath = p.PkgPath
		identifier = parts[1]
		link = true
	}
	return importPath, identifier, link
}

func array(f *api.Field, lp *loadedPackages) (string, string) {
	var importPath, identifier string
	return importPath, identifier
}
