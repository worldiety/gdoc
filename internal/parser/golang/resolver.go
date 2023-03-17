package golang

import (
	"fmt"
	"github.com/worldiety/gdoc/internal/api"
	"go/types"
	"golang.org/x/tools/go/packages"
	"reflect"
	"strings"
)

type loadedPackages struct {
	pkgs  map[string]*packages.Package
	fnMap map[string]*api.Function
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

	fnSignatures(m, lp)
	addRefIDs(m)

	return nil
}

func (lp *loadedPackages) loadPackages(dir string) error {
	pkgs, err := packages.Load(&packages.Config{
		Mode:  packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles | packages.NeedImports | packages.NeedDeps | packages.NeedExportFile | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedTypesSizes | packages.NeedModule | packages.NeedEmbedFiles | packages.NeedEmbedPatterns,
		Tests: false,
	}, dir)
	if err != nil {
		return fmt.Errorf("could not load packages from %s: %w", dir, err)
	}

	for _, pkg := range pkgs {
		lp.pkgs[pkg.Name] = pkg
	}

	return nil
}

func fnSignatures(m *api.Module, lp *loadedPackages) {
	for _, p := range m.Packages {
		for _, fn := range p.Functions {
			lp.fnMap[fn.Name] = fn
		}
	}

	for _, pkg := range lp.pkgs {
		for n, function := range lp.fnMap {
			if fn := pkg.Types.Scope().Lookup(n); fn != nil {
				function.Signature = addNameToSignature(fn)
			}
		}
	}
}

func addRefIDs(m *api.Module) {
	for path, p := range m.Packages {
		p.ID = api.RefId{
			ImportPath: path,
			Identifier: p.Name,
		}
		for id, variable := range p.Vars {
			variable.TypeDefinition = api.RefId{
				ImportPath: path,
				Identifier: id,
			}
			p.Types[variable.TypeDefinition.ID()] = variable.TypeDefinition
		}
		for id, constant := range p.Consts {
			constant.TypeDefinition = api.RefId{
				ImportPath: path,
				Identifier: id,
			}
			p.Types[constant.TypeDefinition.ID()] = constant.TypeDefinition
		}
		for id, function := range p.Functions {
			function.TypeDefinition = api.RefId{
				ImportPath: path,
				Identifier: id,
			}
			p.Types[function.TypeDefinition.ID()] = function.TypeDefinition
		}
	}
}

func addNameToSignature(fn types.Object) string {
	sigString := strings.Replace(fn.Type().String(), fn.Pkg().Path(), fn.Pkg().Name(), -1)
	params := fn.Type().(*types.Signature).Params()
	if v, ok := fn.Type().(*types.Signature); ok {
		v.Results()
	} else {
		panic(fmt.Errorf("implement me: %v", reflect.TypeOf(fn.Type())))
	}
	results := fn.Type().(*types.Signature).Results()
	sigString = replaceParameterString(params, sigString)
	sigString = replaceParameterString(results, sigString)

	return fmt.Sprintf("%s%s%s", sigString[:4], " "+fn.Name(), sigString[4:])
}

func replaceParameterString(params *types.Tuple, sigString string) string {
	for i := 0; i < params.Len(); i++ {
		p := params.At(i)
		origin := p.Origin().Type().String()
		replacement := replacementString(origin)
		replacement = addCrossRef(replacement)
		sigString = strings.Replace(sigString, origin, replacement, -1)
	}

	return sigString
}

func replacementString(origin string) string {
	if !strings.Contains(origin, "/") {
		return origin
	}
	lastSlashIndex := strings.LastIndex(origin, "/")
	var firstReplacementIndex int
	if strings.Contains(origin, "*") {
		firstReplacementIndex = strings.LastIndex(origin, "*") + 1
	}

	return strings.Replace(origin, origin[firstReplacementIndex:lastSlashIndex+1], "", -1)
}

func addCrossRef(s string) string {
	if strings.Contains(s, ".") {
		var firstIndex int
		if strings.Contains(s, "*") {
			firstIndex = strings.LastIndex(s, "*") + 1
		}
		lastDotIndex := strings.LastIndex(s, ".")
		pName := s[firstIndex:lastDotIndex]
		s = strings.Replace(s, pName, fmt.Sprintf("<<%s>>", pName), -1)
	}
	return s
}
