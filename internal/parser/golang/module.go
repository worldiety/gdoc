package golang

import (
	"fmt"
	"github.com/worldiety/gdoc/internal/api"
	"go/ast"
	"go/doc"
	"os"
	"path/filepath"
	"reflect"
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
					TypeDesc:    &api.TypeDesc{},
					Stereotypes: []api.Stereotype{api.StereotypeStruct},
				}
			}
		}
	}

	if len(pkg.dpkg.Vars) > 0 {
		p.Vars = map[string]*api.Variable{}
		for _, value := range pkg.dpkg.Vars {
			for _, d := range newValue(value) {
				p.Vars[d.name] = &api.Variable{
					Name:        d.name,
					Comment:     d.doc,
					TypeDesc:    &api.TypeDesc{},
					Stereotypes: []api.Stereotype{api.StereotypeProperty}}
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
							f = append(f, newField(field, myStruct, ident))
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
			dst["__"+strconv.Itoa(fnum)] = &api.Field{
				Comment: in.Comment,
				TypeDesc: &api.TypeDesc{
					SrcTypeDefinition: in.TypeDesc.SrcTypeDefinition,
					Pointer:           in.TypeDesc.Pointer,
				},
				Stereotypes: st,
			}
			continue
		}

		for _, name := range field.Names {
			c++
			in := newField(field, nil, name)
			myName := in.Name
			if myName == "" {
				myName = "__" + strconv.Itoa(c)
			}
			dst[name.Name] = &api.Field{
				Name:    myName,
				Comment: in.Comment,
				TypeDesc: &api.TypeDesc{
					SrcTypeDefinition: in.TypeDesc.SrcTypeDefinition,
					Pointer:           in.TypeDesc.Pointer,
				},
				Stereotypes: st,
			}
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
	if isMapField(f) {
		mapType := f.Type.(*ast.MapType)
		m.KeyType = &api.TypeDesc{
			SrcTypeDefinition: ast2str(mapType.Key),
			Pointer:           isPointerType(mapType.Key),
		}
		m.ValueType = &api.TypeDesc{
			SrcTypeDefinition: ast2str(mapType.Value),
			Pointer:           isPointerType(mapType.Value),
		}
	}

	n := &api.Field{
		Name:    name,
		Comment: f.Doc.Text(),
		TypeDesc: &api.TypeDesc{
			SrcTypeDefinition: ast2str(f.Type),
			Pointer:           isPointerField(f),
			MapType:           m,
		},
		ParentStruct: s,
		Stereotypes:  []api.Stereotype{api.StereotypeProperty},
	}

	return n
}

// isMapField checks if the given *ast.Field is a map or not
func isMapField(field *ast.Field) bool {
	if field == nil {
		return false
	}
	if field.Type == nil {
		return false
	}
	_, isMap := field.Type.(*ast.MapType)
	return isMap
}

func isPointerField(field *ast.Field) bool {
	_, isPointer := field.Type.(*ast.StarExpr)
	return isPointer
}

// isPointerType checks if the given ast.Expr is a pointer type
func isPointerType(expr ast.Expr) bool {
	if expr == nil {
		return false
	}
	switch t := expr.(type) {
	case *ast.StarExpr:
		return true
	case *ast.Ident:
		return reflect.TypeOf(t.Obj).Kind() == reflect.Ptr
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
		return "..."
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
