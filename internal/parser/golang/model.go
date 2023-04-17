package golang

import (
	"fmt"
	"github.com/worldiety/gdoc/internal/api"
	"strings"
)

const (
	readmeTitle          = "Readme"
	moduleTitlePrefix    = "Module"
	packageTitlePrefix   = "Package"
	structsTitlePrefix   = "Structs"
	funcsTitlePrefix     = "Functions"
	variablesTitlePrefix = "Variables"
	funcTitlePrefix      = "func"
	varPrefix            = "var"
	toc                  = ":toc:"
	filteredFieldsNotice = "// contains filtered or unexported fields"
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
	return fmt.Sprintf("%s%s%s%s", m.title(), simpleLinebreak, toc, m.readme())
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
func (id APackageRefID) String() string {
	return fmt.Sprintf("%s", enclosingDoubleBrackets(angle, fmt.Sprintf("%s,%s%s", id.Identifier, ws, id.Identifier)))
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

func (p APackage) String() string {
	return fmt.Sprintf("%s%s", p.title(), p.readme())
}

type AComment string

func (ac AComment) String() string {
	if ac != "" {
		return formattedComment(string(ac), false)
	}
	return ""
}

type AFunctionComment string

func (afc AFunctionComment) String() string {
	if afc != "" {
		return formattedComment(string(afc), true)
	}
	return ""
}

type ARefId struct {
	api.RefId
}

func NewARefId(refId api.RefId) ARefId {
	return ARefId{RefId: refId}
}

func (r ARefId) String() string {
	return fmt.Sprintf("%s", enclosingDoubleBrackets(angle, fmt.Sprintf("%s,%s%s", r.ID(), ws, r.Identifier)))
}

func (r ARefId) AnchorID() string {
	return enclosingDoubleBrackets(square, fmt.Sprintf("%s", r.ID()))
}

type AStruct struct {
	api.Struct
}

func NewAStruct(structVal api.Struct) AStruct {
	return AStruct{Struct: structVal}
}

type AStructs map[string]AStruct

func (as AStructs) title() string {
	return title(structsTitlePrefix, "", "", 3)
}

func (as AStructs) String() string {
	return as.title()
}

func NewAStructs(domainStructs map[string]*api.Struct) AStructs {
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

	if fieldsString == "" {
		fieldsString = fmt.Sprintf("%s%s%s", enclosingBrackets(square, info),
			enclose(formatDelimiter, filteredFieldsNotice), preservedLinebreak)
	}

	return codeBlock(fmt.Sprintf("%s%s%s%s", commentString, s.asciidocFormattedSigOpen(), fieldsString, s.asciidocFormattedSigClose()))
}

type AFunction struct {
	api.Function
}

type AFunctions map[string]AFunction

func (af AFunctions) String() string {
	return af.title()
}
func (af AFunctions) title() string {
	return title(funcsTitlePrefix, "", "", 3)
}
func NewAFunction(functionVal api.Function) AFunction {
	return AFunction{Function: functionVal}
}

func NewAFunctions(funcs map[string]*api.Function) AFunctions {
	aFunctions := map[string]AFunction{}
	for _, fn := range funcs {
		aFunctions[fn.Name] = NewAFunction(*fn)
	}
	return aFunctions
}

func (fn AFunction) comment() AFunctionComment {
	return AFunctionComment(fn.Comment)
}
func (fn AFunction) String() string {
	return fmt.Sprintf("%s%s%s%s%s", fn.title(), preservedLinebreak,
		codeBlock(fn.asciidocFormattedSignature()), simpleLinebreak, fn.comment().String())
}

func (fn AFunction) RefID() ARefId {
	return NewARefId(fn.TypeDefinition)
}

func (fn AFunction) title() string {
	return fmt.Sprintf("%s", bold(fmt.Sprintf("%s%s%s%s%s", enclosingBrackets(square, keyword),
		enclose(formatDelimiter, funcTitlePrefix), ws, fn.RefID().AnchorID(), fn.RefID().Identifier)))
}

type AVariable struct {
	api.Variable
}

func NewAVariable(v api.Variable) AVariable {
	return AVariable{Variable: v}
}

func (v AVariable) String() string {
	return fmt.Sprintf("%s%s%s", v.Comment, simpleLinebreak, codeBlock(fmt.Sprintf("%s%s%s%s%s",
		builtinFormat(varPrefix), ws, nameFormat(v.Name), ws, v.asciidocFormattedType())))
}

func (v AVariable) StringRaw() string {
	return fmt.Sprintf("%s%s%s%s%s",
		typeFormat(varPrefix), ws, nameFormat(v.Name), ws, v.asciidocFormattedType())
}

type AVariables map[string]AVariable

func NewAVariables(vars map[string]*api.Variable) AVariables {
	nv := map[string]AVariable{}
	for name, v := range vars {
		nv[name] = NewAVariable(*v)
	}
	return nv
}

func (v AVariables) String() string {
	var s string
	varMap := make(map[string]string)

	for _, current := range v {
		if current.Comment == "" {
			varMap["noComment"] += fmt.Sprintf("%s%s", current.StringRaw(), simpleLinebreak)
		} else {
			varMap["commented"] += fmt.Sprintf("%s%s", current.String(), simpleLinebreak)
		}
	}
	s = fmt.Sprintf("%s%s%s", codeBlock(varMap["noComment"]), simpleLinebreak, varMap["commented"])
	s = strings.Trim(s, "\n")
	return fmt.Sprintf("%s%s%s", v.title(), simpleLinebreak, s)
}

func (v AVariables) title() string {
	return title(variablesTitlePrefix, "", "", 3)
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

func (f AField) String() string {
	var whiteSpace string
	if f.Name != "" {
		whiteSpace = f.asciidocWhiteSpaceBetween()
	}
	s := fmt.Sprintf("%s%s%s%s",
		f.comment().String(),
		f.name().String(),
		whiteSpace,
		f.asciidocFormattedType(),
	)
	return s
}

func (f AField) TypeDescription() ATypeDesc {
	return NewATypeDesc(*f.TypeDesc)
}

func (f AField) comment() AComment {
	return AComment(f.Comment)
}

type AFieldName string

func (afn AFieldName) String() string {
	if afn == "" {
		return ""
	}
	return fmt.Sprintf("%s%s", enclosingBrackets(square, t3xt), enclose(formatDelimiter, string(afn)))
}
func (f AField) name() AFieldName {
	return AFieldName(f.Name)
}

type AMapType struct {
	api.MapType
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
		td.TypeDefinition.ID(), ws, enclosingBrackets(square, typ3), enclose(formatDelimiter, td.Identifier()))))
}

func (td ATypeDesc) externalCustomTypeLink() string {
	// custom type from external package from this project
	return fmt.Sprintf("%s%s%s%s",
		// remove the asterisk to find the linked id, it's still displayed in the doc
		td.Prefix(), enclosingDoubleBrackets(angle, fmt.Sprintf("%s,%s%s%s", td.PkgName(), ws,
			enclosingBrackets(square, typ3), enclose(formatDelimiter, td.PkgName()))), dot,
		enclosingDoubleBrackets(angle, fmt.Sprintf("%s,%s%s%s", td.TypeDefinition.ID(), ws,
			enclosingBrackets(square, typ3), enclose(formatDelimiter, td.Identifier()))))
}

func (td ATypeDesc) externalNonCustomTypeLink() string {
	return fmt.Sprintf("%s%s%s%s%s", enclosingBrackets(square, typ3), enclose(formatDelimiter, td.PkgName()), dot, enclosingBrackets(square, typ3), enclose(formatDelimiter, td.Identifier()))
}

func (td ATypeDesc) builtInTypeLink() string {
	return builtinFormat(td.SrcTypeDefinition)
}
