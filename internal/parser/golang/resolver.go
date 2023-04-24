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

	addTypeInformation(m, lp)
	addCommentLinks(m)

	return nil
}

func addCommentLinks(m *api.Module) {
	for _, p := range m.Packages {
		for _, function := range p.Functions {
			function.Comment = handleComment(function.Comment, p, m)
		}
		for _, s := range p.Structs {
			s.Comment = handleComment(s.Comment, p, m)
		}
		for _, v := range p.Vars {
			v.Comment = handleComment(v.Comment, p, m)
			v.Doc = handleComment(v.Doc, p, m)
		}
		for _, consts := range p.Consts {
			for _, c := range consts.Content {
				c.Comment = handleComment(c.Comment, p, m)
			}
		}
	}
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
func addTypeInformation(m *api.Module, lp *loadedPackages) {
	for path, p := range m.Packages {
		p.PackageDefinition = api.NewRefID(path, p.Name)
		addVariableInfo(p, lp, path)
		addConstantInfo(p, path)
		addFunctionInfo(p, lp, path)
		addStructInfo(p, lp, path)
	}
}
func addVariableInfo(p *api.Package, lp *loadedPackages, path string) {
	for id, v := range p.Vars {
		v.TypeDesc.TypeDefinition = api.NewRefID(path, id)
		p.Types[v.Name] = v.TypeDesc.TypeDefinition
		typeDescInfo(v.Name, v.TypeDesc, lp)
	}
}
func addConstantInfo(p *api.Package, path string) {
	for _, block := range p.Consts {
		for _, constant := range block.Content {
			constant.RefId = api.NewRefID(path, p.PackageDefinition.ImportPath)
			p.Types[constant.RefId.Identifier] = constant.RefId
		}
	}
}
func addFunctionInfo(p *api.Package, lp *loadedPackages, path string) {

	for id, function := range p.Functions {
		function.TypeDefinition = api.NewRefID(path, id)
		p.Types[function.Name] = function.TypeDefinition
		handleFields(function.Parameters, p, lp)
		handleFields(function.Results, p, lp)
	}
}
func addStructInfo(p *api.Package, lp *loadedPackages, path string) {
	for id, s := range p.Structs {
		s.TypeDefinition = api.NewRefID(path, id)
		p.Types[s.Name] = s.TypeDefinition

		for _, f := range s.Fields {
			handleField(f, p, lp)
		}
		for _, g := range s.Generics {
			handleField(g, p, lp)
		}
		for _, m := range s.Methods {
			handleMethod(m, p, lp)
		}
	}
}

func handleMethod(m *api.Method, p *api.Package, lp *loadedPackages) {
	typeDescInfo(p.Name, m.Recv.TypeDesc, lp)
}

func handleComment(comment string, p *api.Package, m *api.Module) string {
	if comment == "" {
		return comment
	}
	replacementMap := map[string]string{}
	// split word
	for _, s := range strings.Split(comment, ws) {
		// check for type from external package
		if externalPkg(s) {
			parts := strings.Split(s, dot)
			for path, extPkg := range m.Packages {
				// check import paths for ext package name
				if strings.HasSuffix(path, parts[0]) {
					// add replacement string for pkg name to map
					pkgReplacement := NewAPackageRefID(extPkg.PackageDefinition).String()
					var typeReplacement string
					if t, ok := extPkg.Types[parts[1]]; ok {
						// add replacement string for external type to map
						typeReplacement = NewARefId(t).String()
					}
					replacementMap[s] = fmt.Sprintf("%s%s%s", pkgReplacement, dot, typeReplacement)
				}
			}
		} else if t, ok := p.Types[s]; ok {
			// if from current package
			replacementMap[t.Identifier] = NewARefId(t).String()
		} else if enclosedInSquareBrackets(s) {
			if t, ok := p.Types[removeEnclosingSquaredBrackets(s)]; ok {
				replacementMap[s] = NewARefId(t).String()
			}
		}
	}

	for sToReplace, replacement := range replacementMap {
		// create asciidoc formatted codeBlockComment
		comment = strings.Replace(comment, sToReplace, replacement, 1)
	}
	return comment
}

func handleField(f *api.Field, p *api.Package, lp *loadedPackages) {
	if f.TypeDesc.Map() {
		handleMapType(f, p.Name, lp)
	} else {
		typeDescInfo(p.Name, f.TypeDesc, lp)
	}
}

func handleFields(parameters map[string]*api.Field, currentPackage *api.Package, lp *loadedPackages) {
	for _, p := range parameters {
		handleField(p, currentPackage, lp)
	}
}

func handleMapType(f *api.Field, pName string, lp *loadedPackages) {
	keyTypeDef, valueTypeDef := f.TypeDesc.MapSrcDefs()
	f.TypeDesc.MapType = &api.MapType{}
	f.TypeDesc.MapType.KeyType = &api.TypeDesc{SrcTypeDefinition: keyTypeDef,
		Pointer:   strings.Contains(keyTypeDef, "*"),
		Linebreak: false,
	}
	f.TypeDesc.MapType.ValueType = &api.TypeDesc{SrcTypeDefinition: valueTypeDef,
		Pointer:   strings.Contains(valueTypeDef, "*"),
		Linebreak: true,
	}
	typeDescInfo(pName, f.TypeDesc.MapType.KeyType, lp)
	typeDescInfo(pName, f.TypeDesc.MapType.ValueType, lp)
}

func typeDescInfo(pName string, td *api.TypeDesc, lp *loadedPackages) {
	var importPath string
	var origin api.TypeOrigin

	// from current package or built-in
	if currentPackageOrBuiltIn(td) {
		if p := lp.pkgs[pName]; p != nil {
			if t := p.Types.Scope().Lookup(td.Identifier()); t != nil {
				importPath = p.PkgPath
				origin = api.LocalCustom
			}
		}
	} else {
		// from other package
		if ok, p := externalCustomPkg(*td, *lp); ok {
			importPath = p.PkgPath
			origin = api.ExternalCustom
		} else {
			origin = api.ExternalNonCustom
		}
	}
	td.TypeDefinition = api.NewRefID(importPath, td.Identifier())
	td.TypeOrigin = origin
}

func externalPkg(s string) bool {
	if strings.Contains(s, dot) {
		return true
	}
	return false
}

func currentPackageOrBuiltIn(td *api.TypeDesc) bool {
	if include, _ := td.IncludesPkg(); !include {
		return true
	}
	return false
}

func externalCustomPkg(td api.TypeDesc, lp loadedPackages) (bool, *packages.Package) {
	if p := lp.pkgs[td.PkgName()]; p != nil {
		return true, p
	}
	return false, nil
}

func enclosedInSquareBrackets(s string) bool {
	if strings.HasPrefix(s, "[") && strings.HasSuffix(s, "]") {
		return true
	}
	return false
}
func removeEnclosingSquaredBrackets(s string) string {
	var result string
	if after, ok := strings.CutPrefix(s, "["); ok {
		result, _ = strings.CutSuffix(after, "]")
	}
	return result

}
