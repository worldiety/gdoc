package golang

import (
	"fmt"
	"github.com/worldiety/gdoc/internal/api"
)

const (
	readmeTitle        = "**__Readme__**"
	moduleTitlePrefix  = "Module"
	packageTitlePrefix = "Package"
	toc                = ":toc:"
)

type ImportPath = string
type AModule struct {
	Readme   string
	Name     string
	Packages map[ImportPath]APackage
}

func NewAModule(module api.Module) AModule {
	return AModule{
		Readme:   module.Readme,
		Name:     module.Name,
		Packages: NewAPackages(module.Packages),
	}
}

func (m AModule) String() string {
	return fmt.Sprintf("%s\n%s%s", m.title(), toc, m.readme())
}

type APackage struct {
	api.Package
}

func NewAPackage(packageVal api.Package) APackage {
	return APackage{Package: packageVal}
}

func NewAPackages(packagesVal map[ImportPath]*api.Package) map[string]APackage {
	pkgs := map[string]APackage{}
	for importPath, p := range packagesVal {
		pkgs[importPath] = NewAPackage(*p)
	}
	return pkgs
}

func (p APackage) String() string {
	return fmt.Sprintf("%s%s", p.title(), p.readme())
}

type ARefId struct {
	api.RefId
}

func NewARefId(refId api.RefId) ARefId {
	return ARefId{RefId: refId}
}

func (r ARefId) String() string {
	return fmt.Sprintf("ARefId{RefId: %v}", r.RefId)
}

type AStruct struct {
	api.Struct
}

func NewAStruct(structVal api.Struct) AStruct {
	return AStruct{Struct: structVal}
}

func NewAStructs(domainStructs map[string]*api.Struct) map[string]AStruct {
	aStructs := map[string]AStruct{}
	for _, s := range domainStructs {
		aStructs[s.Name] = NewAStruct(*s)
	}
	return aStructs
}

func (s AStruct) String() string {
	var commentString string
	if s.Comment != "" {
		commentString = s.asciidocFormattedComment()
	}
	var fieldsString string
	for _, f := range s.AFields() {
		fieldsString += f.String()
	}

	return codeBlock(fmt.Sprintf("%s%s%s%s", commentString, s.asciidocFormattedSigOpen(), fieldsString, s.asciidocFormattedSigClose()))
}

type AFunction struct {
	api.Function
}

func NewAFunction(functionVal api.Function) AFunction {
	return AFunction{Function: functionVal}
}

func NewAFunctions(funcs map[string]*api.Function) map[string]AFunction {
	var aFunctions map[string]AFunction
	for _, fn := range funcs {
		aFunctions[fn.Name] = NewAFunction(*fn)
	}
	return aFunctions
}

func (fn AFunction) String() string {
	return fmt.Sprintf("AFunction{Function: %v}", fn.Function)
}

type AField struct {
	api.Field
}

func NewAField(fieldVal api.Field) AField {
	return AField{Field: fieldVal}
}

func (s AStruct) AFields() []AField {
	aFields := make([]AField, 0)
	for _, f := range s.Fields {
		aFields = append(aFields, NewAField(*f))
	}
	return aFields
}

func (f AField) String() string {

	return fmt.Sprintf("%s%s%s%s",
		f.asciidocFormattedComment(),
		f.asciidocFormattedName(),
		f.asciidocWhiteSpaceBetween(),
		f.asciidocFormattedType(),
	)
}

func (f AField) TypeDescription() ATypeDesc {
	return NewATypeDesc(*f.TypeDesc)
}

type AMapType struct {
	api.MapType
}

func NewAMapType(mapTypeVal api.MapType) AMapType {
	return AMapType{MapType: mapTypeVal}
}

func (m AMapType) String() string {
	return fmt.Sprintf("AMapType{MapType: %v}", m.MapType)
}

type ATypeDesc struct {
	api.TypeDesc
}

func NewATypeDesc(typeDescVal api.TypeDesc) ATypeDesc {
	return ATypeDesc{TypeDesc: typeDescVal}
}

func (td ATypeDesc) RefId() ARefId {
	return NewARefId(td.TypeDefinition)
}

func (td ATypeDesc) String() string {
	var s string
	if td.Pointer {
		s += "*"
	}

	s += td.RefId().String()

	return s
}
