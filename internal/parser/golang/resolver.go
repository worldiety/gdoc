package golang

import (
	"fmt"
	"go/types"
	"godocgenerator/internal/api"
	"golang.org/x/tools/go/packages"
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
	packageTypes(m, lp)

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

func packageTypes(m *api.Module, lp *loadedPackages) {
	for _, pkg := range lp.pkgs {
		for i, tn := range pkg.Types.Scope().Names() {
			if m.Packages[pkg.PkgPath].Types == nil {
				m.Packages[pkg.PkgPath].Types = make(map[int]api.BaseType, 0)
			}
			if pkg.Types.Scope().Lookup(tn).Exported() {
				m.Packages[pkg.PkgPath].Types[i] = api.BaseType(tn)
			}
		}
	}
}

func addNameToSignature(fn types.Object) string {
	sigString := strings.Replace(fn.Type().String(), fn.Pkg().Path(), fn.Pkg().Name(), -1)
	params := fn.Type().(*types.Signature).Params()
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
