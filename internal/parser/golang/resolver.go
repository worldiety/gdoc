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

	addRefIDs(m, lp)

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

func addRefIDs(m *api.Module, lp *loadedPackages) {
	for path, p := range m.Packages {
		p.PackageDefinition = api.RefId{
			ImportPath: path,
			Identifier: p.Name,
		}
		for id, variable := range p.Vars {
			variable.TypeDefinition = api.RefId{
				ImportPath: path,
				Identifier: id,
			}
			p.Types[variable.Name] = variable.TypeDefinition
		}
		for id, constant := range p.Consts {
			constant.TypeDefinition = api.RefId{
				ImportPath: path,
				Identifier: id,
			}
			p.Types[constant.Name] = constant.TypeDefinition
		}
		for id, function := range p.Functions {
			function.TypeDefinition = api.RefId{
				ImportPath: path,
				Identifier: id,
			}
			p.Types[function.Name] = function.TypeDefinition
		}
		for id, s := range p.Structs {
			s.TypeDefinition = api.RefId{
				ImportPath: path,
				Identifier: id,
			}
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

				f.TypeDefinition = api.RefId{
					ImportPath: importPath,
					Identifier: identifier,
				}

				if lp.pkgs[f.TypeDefinition.PackageName()] != nil && lp.pkgs[f.TypeDefinition.PackageName()].Types.Scope().
					Lookup(strings.Replace(identifier, "*", "", -1)) != nil {
					f.Link = true
				}
			}
		}
	}
}
