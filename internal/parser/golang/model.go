package golang

import (
	"fmt"
	"github.com/worldiety/gdoc/internal/api"
	"golang.org/x/exp/slices"
	"strings"
)

const (
	readmeTitle          = "Readme"
	moduleTitlePrefix    = "Module"
	packageTitlePrefix   = "Package"
	structsTitlePrefix   = "Structs"
	structTitlePrefix    = "Struct"
	funcsTitlePrefix     = "Functions"
	variablesTitlePrefix = "Variables"
	constantsTitlePrefix = "Consts"
	funcTitlePrefix      = "func"
	varPrefix            = "var"
	toc                  = ":toc:"
	docInfo              = ":docinfo: shared"
	theme                = ":pdf-theme: my-theme.yml"
	filteredFieldsNotice = "// contains filtered or unexported fields"
)

type ImportPath = string

type AsciiDocHeader struct {
	Attributes []string
}

func NewAsciiDocHeader() AsciiDocHeader {
	s := []string{docInfo, toc, theme}
	return AsciiDocHeader{Attributes: s}
}

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

type APackageRefID struct {
	api.RefId
}

func NewAPackageRefID(id api.RefId) APackageRefID {
	return APackageRefID{RefId: id}
}
func (p APackage) RefID() APackageRefID {
	return NewAPackageRefID(p.PackageDefinition)
}

// APackage is a decorator struct for the api.Package struct
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

type ADoc struct {
	AComment
}

type AComment struct {
	Raw   string
	Lines []string
}

func NewADoc(s string) ADoc {
	return ADoc{NewAComment(s)}
}

func NewAComment(s string) AComment {
	var lines []string
	for _, line := range strings.Split(s, simpleLinebreak) {
		lines = append(lines, line)
	}
	return AComment{
		Raw:   s,
		Lines: lines,
	}
}

type AFunctionComment string

type ARefId struct {
	api.RefId
}

func NewARefId(refId api.RefId) ARefId {
	return ARefId{RefId: refId}
}

func (r ARefId) AnchorID() string {
	return enclosingDoubleBrackets(square, r.ID())
}

type AStruct struct {
	api.Struct
}

func (s AStruct) title() string {
	return bold(keywordFormat(structTitlePrefix), ws, nameFormat(s.Name))
}

func NewAStruct(structVal api.Struct) AStruct {
	return AStruct{Struct: structVal}
}

func (s AStruct) methods() AMethods {
	return NewAMethods(s.Methods)
}

type AStructs map[string]AStruct

func (as AStructs) title() string {
	return title(structsTitlePrefix, "", "", 3)
}

func NewAStructs(domainStructs map[string]*api.Struct) AStructs {
	aStructs := map[string]AStruct{}
	for _, s := range domainStructs {
		aStructs[s.Name] = NewAStruct(*s)
	}
	return aStructs
}

func (s AStruct) comment() AComment {
	return NewAComment(s.Comment)
}

func (s AStruct) generics() AGenerics {
	return NewAGenerics(s.Generics)
}

type AGenerics []AField

func NewAGenerics(generics api.Generics) AGenerics {
	var result []AField
	for _, field := range generics {
		result = append(result, NewAField(*field))
	}
	return result
}

type AFunction struct {
	api.Function
}

func (af AFunctions) title() string {
	return title(funcsTitlePrefix, "", "", 3)
}
func NewAFunction(functionVal api.Function) AFunction {
	return AFunction{Function: functionVal}
}

func (fn AFunction) comment() AFunctionComment {
	return AFunctionComment(fn.Comment)
}

func (fn AFunction) RefID() ARefId {
	return NewARefId(fn.TypeDefinition)
}

func (fn AFunction) name() string {
	return name(fn, nil)
}

func (m AMethod) function() AFunction {
	return NewAFunction(*m.Function)
}

func (m AMethod) recv() *ARecv {
	r := NewARecv(m.Recv)
	return &r
}
func (m AMethod) name() string {
	return name(m.function(), m.recv())
}

func (r ARecv) String() string {
	return enclosingBrackets(round, fmt.Sprintf("%s%s%s", variableFormat(r.Name), ws, typeFormat(r.TypeString)))
}

func name(fn AFunction, recv *ARecv) string {
	if recv == nil {
		return fmt.Sprintf("%s%s%s%s%s", enclosingBrackets(square, keyword),
			enclose(hash, funcTitlePrefix), ws, fn.RefID().AnchorID(), nameFormat(fn.RefID().Identifier))
	} else {
		return fmt.Sprintf("%s%s%s%s%s%s", enclosingBrackets(square, keyword),
			enclose(hash, funcTitlePrefix), ws, recv.String(), ws, nameFormat(fn.Name))
	}
}

type AFunctions map[string]AFunction

func NewAFunctions(funcs map[string]*api.Function) AFunctions {
	aFunctions := map[string]AFunction{}
	for _, fn := range funcs {
		aFunctions[fn.Name] = NewAFunction(*fn)
	}
	return aFunctions
}

type AMethod struct {
	api.Method
}

func NewAMethod(methodVal api.Method) AMethod {
	return AMethod{methodVal}
}

type ARecv struct {
	*api.Recv
}

func NewARecv(s *api.Recv) ARecv {
	return ARecv{s}
}

type AMethods map[string]AMethod

func NewAMethods(methods []*api.Method) AMethods {
	aMethods := map[string]AMethod{}
	for _, m := range methods {
		aMethods[m.Name] = NewAMethod(*m)
	}
	return aMethods
}

type AVariable struct {
	api.Variable
}

func (v AVariable) AnchorID() string {
	return enclosingDoubleBrackets(square, v.TypeDesc.TypeDefinition.ID())
}

func (v AVariable) name() AFieldName {
	return AFieldName(v.Name)
}

func NewAVariable(v api.Variable) AVariable {
	return AVariable{Variable: v}
}
func (v AVariable) StringRaw() string {
	return fmt.Sprintf("%s%s%s%s%s",
		typeFormat(varPrefix), ws, variableFormat(v.Name), ws, v.asciidocFormattedType())
}

type AVariables map[string]AVariable

func NewAVariables(vars map[string]*api.Variable) AVariables {
	nv := map[string]AVariable{}
	for name, v := range vars {
		nv[name] = NewAVariable(*v)
	}
	return nv
}

func (v AVariables) sort() []AVariable {
	return SortMapValues(v, func(a, b AVariable) bool {
		return a.Name < b.Name
	})
}

func (consts AConstBlock) sort() AConstBlock {
	slices.SortFunc(consts.Content, func(a, b api.Constant) bool {
		return a.RefId.Identifier < b.RefId.Identifier
	})
	return consts
}

func (consts AConstBlock) consts() []AConst {
	res := make([]AConst, 0)
	for _, constant := range consts.Content {
		res = append(res, NewAConst(constant))
	}
	return res
}

type AConst struct {
	api.Constant
}

func (c AConst) name() AFieldName {
	return AFieldName(c.RefId.Identifier)
}

func NewAConst(c api.Constant) AConst {
	return AConst{c}
}

type AConstBlock struct {
	api.ConstantBlock
}

func NewAConstBlock(consts api.ConstantBlock) AConstBlock {
	return AConstBlock{consts}
}

type AConstBlockList []AConstBlock

func NewAConstBlockList(blocks []api.ConstantBlock) AConstBlockList {
	res := make([]AConstBlock, 0)
	for _, block := range blocks {
		res = append(res, NewAConstBlock(block))
	}
	return res
}

func (ms AMethods) sort() []AMethod {
	return SortMapValues(ms, func(a, b AMethod) bool {
		return a.Name < b.Name
	})
}

type commentStatus int

const (
	uncommented commentStatus = iota
	commented
)

func (AVariables) title() string {
	return title(variablesTitlePrefix, "", "", 3)
}

func (AConstBlockList) title() string {
	return title(constantsTitlePrefix, "", "", 3)
}

type AField struct {
	api.Field
}

type AFields map[string]AField

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

func (f AField) typeDescription() ATypeDesc {
	return NewATypeDesc(*f.TypeDesc)
}

func (f AField) comment() AComment {
	return NewAComment(f.Comment)
}

func (f AField) doc() ADoc {
	return NewADoc(f.Doc)
}

type AFieldName string

func (f AField) name() AFieldName {
	return AFieldName(f.Name)
}

type AMapType struct {
	api.MapType
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

func (td ATypeDesc) Prefix() string {
	var s string
	if td.Array() {
		var prefixEnd int
		if td.lastAsteriskIndex() > td.lastClosedArrayBracketIndex() {
			prefixEnd = td.lastAsteriskIndex()
		} else {
			prefixEnd = td.lastClosedArrayBracketIndex()
		}
		s = td.SrcTypeDefinition[:prefixEnd+1]
	} else if td.Pointer {
		s = td.SrcTypeDefinition[:td.lastAsteriskIndex()+1]
	}
	escapeChar := "*"
	s = strings.Replace(s, escapeChar, passThrough(escapeChar), -1)

	return s
}

func (td ATypeDesc) lastAsteriskIndex() int {
	return strings.LastIndex(td.SrcTypeDefinition, "*")
}

func (td ATypeDesc) lastClosedArrayBracketIndex() int {
	return strings.LastIndex(td.SrcTypeDefinition, "]")
}

func (td ATypeDesc) localCustomTypeLink() string {
	return fmt.Sprintf("%s%s", td.Prefix(), enclosingDoubleBrackets(angle, fmt.Sprintf("%s,%s%s%s",
		td.TypeDefinition.ID(), ws, enclosingBrackets(square, typ3), enclose(hash, td.Identifier()))))
}

func (td ATypeDesc) externalCustomTypeLink() string {
	// custom type from external package from this project
	return fmt.Sprintf("%s%s%s%s",
		// remove the asterisk to find the linked id, it's still displayed in the doc
		td.Prefix(), enclosingDoubleBrackets(angle, fmt.Sprintf("%s,%s%s%s", td.PkgName(), ws,
			enclosingBrackets(square, typ3), enclose(hash, td.PkgName()))), dot,
		enclosingDoubleBrackets(angle, fmt.Sprintf("%s,%s%s%s", td.TypeDefinition.ID(), ws,
			enclosingBrackets(square, typ3), enclose(hash, td.Identifier()))))
}

func (td ATypeDesc) externalNonCustomTypeLink() string {
	return fmt.Sprintf("%s%s%s%s%s", enclosingBrackets(square, typ3), enclose(hash, td.PkgName()), dot, enclosingBrackets(square, typ3), enclose(hash, td.Identifier()))
}

func (td ATypeDesc) builtInTypeLink() string {
	return builtinFormat(td.SrcTypeDefinition)
}
