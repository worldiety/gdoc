package golang

import (
	"fmt"
	"golang.org/x/tools/go/packages"
)

func Resolve(dir string) ([]*packages.Package, error) {
	pkgs, err := packages.Load(&packages.Config{
		Mode:  packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles | packages.NeedImports | packages.NeedDeps | packages.NeedExportFile | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedTypesSizes | packages.NeedModule | packages.NeedEmbedFiles | packages.NeedEmbedPatterns,
		Tests: true,
	}, dir)
	if err != nil {
		return nil, err
	}
	for _, p := range pkgs {
		fmt.Println(p.Types)
		fmt.Println()
	}

	return pkgs, nil
}
