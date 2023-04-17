package golang

import (
	"fmt"
	"github.com/worldiety/gdoc/internal/api"
	"go/ast"
	"go/doc"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func newModule(dir string, modname string, pkgs map[string]Package) (*api.Module, error) {
	m := &api.Module{
		Name:   modname,
		Readme: tryLoadReadme(dir),
	}

	if len(pkgs) > 0 {
		m.Packages = map[api.ImportPath]*api.Package{}
		for _, p := range pkgs {
			np := newPackage(p)
			np.Readme = tryLoadReadme(p.dir)
			m.Packages[p.dpkg.ImportPath] = np
		}
	}

	return m, nil
}

func tryLoadReadme(dir string) string {
	files, _ := os.ReadDir(dir)
	for _, file := range files {
		if file.Type().IsRegular() && strings.ToLower(file.Name()) == "readme.md" {
			buf, _ := os.ReadFile(filepath.Join(dir, file.Name()))
			if len(buf) != 0 {
				return string(buf)
			}
		}
	}

	return ""
}

func newPackage(pkg Package) *api.Package {
	var tmpImports api.Imports

	for _, s := range pkg.dpkg.Imports {
		tmpImports = append(tmpImports, api.Import(s))
	}
	tmpImports.Sort()

	p := &api.Package{
		Doc:     pkg.dpkg.Doc,
		Name:    pkg.dpkg.Name,
		Imports: tmpImports,
		Types:   map[string]api.RefId{},
	}

	if pkg.dpkg.Name == "main" {
		p.Stereotypes = append(p.Stereotypes, api.StereotypeExecutable)
	}

	if len(pkg.dpkg.Funcs) > 0 {
		p.Functions = map[string]*api.Function{}
		for _, f := range pkg.dpkg.Funcs {
			if !f.Decl.Name.IsExported() {
				continue
			}

			p.Functions[f.Name] = newFunc(f)
		}

	}

	if len(pkg.dpkg.Consts) > 0 {
		p.Consts = map[string]*api.Constant{}
		for _, value := range pkg.dpkg.Consts {
			for _, d := range newValue(value) {
				p.Consts[d.name] = &api.Constant{
					Comment:     d.doc,
					Name:        d.name,
					TypeDesc:    &api.TypeDesc{Linebreak: true},
					Stereotypes: []api.Stereotype{api.StereotypeStruct},
				}
			}
		}
	}

	if len(pkg.dpkg.Vars) > 0 {
		p.Vars = map[string]*api.Variable{}
		for _, value := range pkg.dpkg.Vars {
			for _, spec := range value.Decl.Specs {
				switch t := spec.(type) {
				case *ast.ValueSpec:
					for _, ident := range t.Names {
						if isExported(ident.Name) {
							p.Vars[ident.Name] =
								api.NewVariable(ident.Name, t.Comment.Text(),
									&api.TypeDesc{
										SrcTypeDefinition: ast2str(t.Type),
										Pointer:           isPointerType(t.Type),
										Linebreak:         true,
									})
						}
					}
				}
			}
		}
	}

	if len(pkg.dpkg.Types) > 0 {
		p.Structs = map[string]*api.Struct{}
		for _, value := range pkg.dpkg.Types {
			if isExported(value.Name) {
				p.Structs[value.Name] = newStruct(value)
			}
		}

	}

	return p
}

func newStruct(value *doc.Type) *api.Struct {
	var f []*api.Field
	myStruct := &api.Struct{
		Comment: strings.Trim(value.Doc, "\n"),
		Name:    value.Name,
	}

	for _, spec := range value.Decl.Specs {
		switch s := spec.(type) {
		case *ast.TypeSpec:
			if structType, ok := s.Type.(*ast.StructType); ok {

				for _, field := range structType.Fields.List {
					for _, ident := range field.Names {
						if isExported(ident.Name) {
							field := newField(field, myStruct, ident)
							field.TypeDesc.Linebreak = true
							f = append(f, field)
						}
					}
				}
			}
		}
	}

	myStruct.Fields = f

	return myStruct
}

func newFunc(docFunc *doc.Func) *api.Function {
	f := &api.Function{
		Name:    docFunc.Name,
		Comment: docFunc.Doc,
	}

	fn := docFunc.Decl.Type

	inArgs := fn.Params.List
	if len(inArgs) > 0 {
		f.Parameters = map[string]*api.Field{}
		insertParams(f.Parameters, inArgs, api.StereotypeParameter, api.StereotypeParameterIn)
	}

	if fn.Results != nil {
		outArgs := fn.Results.List
		if len(outArgs) > 0 {
			f.Results = map[string]*api.Field{}
			insertParams(f.Results, outArgs, api.StereotypeParameter, api.StereotypeParameterOut, api.StereotypeParameterResult)
		}
	}

	return f
}
func insertParams(dst map[string]*api.Field, src []*ast.Field, st ...api.Stereotype) {
	c := 0
	for fnum, field := range src {
		if len(field.Names) == 0 {
			in := newField(field, nil, nil)
			dst["__"+strconv.Itoa(fnum)] = api.NewField("", in.Comment,
				api.NewTypeDesc(
					api.RefId{}, in.TypeDesc.SrcTypeDefinition, in.TypeDesc.Pointer, nil), nil)
			dst["__"+strconv.Itoa(fnum)].Stereotypes = st
			continue
		}

		for _, name := range field.Names {
			c++
			in := newField(field, nil, name)
			myName := in.Name
			if myName == "" {
				myName = "__" + strconv.Itoa(c)
			}
			dst[name.Name] =
				api.NewField(myName, in.Comment,
					api.NewTypeDesc(
						api.RefId{}, in.TypeDesc.SrcTypeDefinition, in.TypeDesc.Pointer, nil), nil)
			dst[name.Name].Stereotypes = st
		}
	}
}

func newField(f *ast.Field, s *api.Struct, id *ast.Ident) *api.Field {
	var name string
	if id != nil {
		name = id.Name
		if s != nil && len([]rune(name)) > s.WhiteSpaceInFields {
			s.WhiteSpaceInFields = len([]rune(name))
		}
	}

	m := &api.MapType{}
	if ok, mapType := isMapField(f); ok {
		m.KeyType = api.NewTypeDesc(api.RefId{}, ast2str(mapType.Key), isPointerType(mapType.Key), m)
		m.ValueType = api.NewTypeDesc(api.RefId{}, ast2str(mapType.Value), isPointerType(mapType.Value), m)
	}

	n := api.NewField(name, f.Doc.Text(), api.NewTypeDesc(api.RefId{}, ast2str(f.Type), isPointerType(f.Type), m), s)
	n.Stereotypes = []api.Stereotype{api.StereotypeProperty}

	if ok, arrayType := isArrayField(f); ok {
		n.TypeDesc.Pointer = isPointerType(arrayType.Elt)
	}

	return n
}

// isMapField checks if the given *ast.Field is a map or not
func isMapField(f *ast.Field) (bool, *ast.MapType) {
	if f == nil {
		return false, nil
	}
	if f.Type == nil {
		return false, nil
	}
	mapType, isMap := f.Type.(*ast.MapType)
	return isMap, mapType
}

func isArrayField(f *ast.Field) (bool, *ast.ArrayType) {
	if f == nil {
		return false, nil
	}
	if f.Type == nil {
		return false, nil
	}
	arrayType, isArray := f.Type.(*ast.ArrayType)
	return isArray, arrayType
}

// isPointerType checks if the given ast.Expr is a pointer type
func isPointerType(expr ast.Expr) bool {
	if expr == nil {
		return false
	}
	switch expr.(type) {
	case *ast.StarExpr:
		return true
	}
	return false
}

type docValue struct {
	doc  string
	name string
}

func newValue(value *doc.Value) []docValue {
	var res []docValue
	groupDoc := value.Doc
	for _, spec := range value.Decl.Specs {
		switch t := spec.(type) {
		case *ast.ValueSpec:
			actualDoc := t.Doc.Text()
			for _, name := range t.Names {
				if !name.IsExported() {
					continue
				}
				res = append(res, docValue{
					doc:  strings.TrimSpace(groupDoc + "\n" + actualDoc),
					name: name.Name,
				})
			}

		}
	}

	return res
}

func ast2str(n ast.Node) string {
	switch t := n.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		return ast2str(t.X) + "." + t.Sel.Name
	case *ast.StarExpr:
		return "*" + ast2str(t.X)
	case *ast.MapType:
		return "map[" + ast2str(t.Key) + "]" + ast2str(t.Value)
	case *ast.IndexExpr:
		return ast2str(t.X) + "[" + ast2str(t.Index) + "]"
	case *ast.ChanType:
		s := "chan"
		switch t.Dir {
		case ast.SEND:
			s += "<-"
		case ast.RECV:
			s += "->"
		}

		s += ast2str(t.Value)

		return s
	case *ast.ArrayType:
		s := "["
		if t.Len != nil {
			s += ast2str(t.Len)
		}
		s += "]"
		s += ast2str(t.Elt)
		return s
	case *ast.TypeSpec:
		return ast2str(t.Type)
	case *ast.InterfaceType:
		return "interface"
	case *ast.StructType:
		return "struct"
	case *ast.Ellipsis:
		return "..." + ast2str(t.Elt)
	case *ast.FuncType:
		s := "func"
		if t.TypeParams != nil && len(t.TypeParams.List) > 0 {
			s += "["
			s += ast2str(t.TypeParams)
			s += "]"
		}
		s += "("
		if t.Params != nil {
			s += ast2str(t.Params)
		}
		s += ")"
		if t.Results != nil {
			switch len(t.Results.List) {
			case 0:
			case 1:
				s += " "
				s += ast2str(t.Results)
			default:
				s += " ("
				s += ast2str(t.Results)
				s += ")"
			}
		}
		return s
	case *ast.FieldList:
		s := ""
		for _, field := range t.List {
			for _, name := range field.Names {
				s += name.Name + " ,"
			}
			s = strings.TrimSuffix(s, ",")
			s += ast2str(field.Type)
		}

		s = strings.TrimSuffix(s, " ,")
		return s
	default:
		panic(fmt.Errorf("implement me %T", t))
	}
}

func isExported(s string) bool {
	if s[0] >= 'A' && s[0] <= 'Z' {
		return true
	}
	return false
}
