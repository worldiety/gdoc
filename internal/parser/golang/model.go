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
	funcsTitlePrefix     = "Functions"
	variablesTitlePrefix = "Variables"
	funcTitlePrefix      = "func"
	varPrefix            = "var"
	toc                  = ":toc:"
	docInfo              = ":docinfo: shared"
	filteredFieldsNotice = "// contains filtered or unexported fields"
)

type ImportPath = string

type AsciiDocHeader struct {
	Attributes []string
}

func NewAsciiDocHeader() AsciiDocHeader {
	s := []string{docInfo, toc}
	return AsciiDocHeader{Attributes: s}
}

func (h AsciiDocHeader) String() string {
	var s string
	for _, a := range h.Attributes {
		s += fmt.Sprintf("%s%s", a, simpleLinebreak)
	}
	s = strings.Trim(s, simpleLinebreak)
	return s
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

func (m AModule) String() string {
	return fmt.Sprintf("%s%s%s", m.title(), simpleLinebreak, m.readme())
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

func (ac AComment) String() string {
	var original, current string
	var originalList, tmpList []string
	var inList, inIndentedBlock bool

	result := ac.Raw
	for i := 0; i <= len(ac.Lines); i++ {
		if i < len(ac.Lines) {
			current = ac.Lines[i]
			if current == "" {
				continue
			}
		} else {
			current = ""
		}
		original = current
		if startsWithEitherPrefix(current, ws, tab) {
			if strings.HasPrefix(strings.Trim(current, ws), hyphen) {
				// is list -> change godoc list marker (-) for asciidoc list marker (*)
				originalList = append(originalList, current)
				current = strings.Replace(current, hyphen, asterisk, 1)
				tmpList = append(tmpList, current)
				inList = true
				continue
			} else {
				var count int
				if !inIndentedBlock {
					originalList = []string{}
					tmpList = []string{}
				}
				originalList = append(originalList, original)
				for _, r := range current {
					if string(r) == tab {
						count += 4
						continue
					}
					if string(r) == ws {
						count++
						continue
					}
					current = trimAllPrefixWSAndTabs(current)
					for i := 0; i < count; i++ {
						current = nbsp + current
					}
					current += plusSuffix
					tmpList = append(tmpList, current)
					inIndentedBlock = true
					break
				}
			}
		} else if inList || inIndentedBlock {
			// format full list
			tmp := formatBlock(tmpList...)
			if inList {
				tmp = simpleLinebreak + tmp
			}
			original = ""
			for _, s := range originalList {
				original += s + simpleLinebreak
			}

			if inIndentedBlock {
				tmp = mono + enclose(hash, trimAllSuffixLinebreaks(tmp)) + plusSuffix
			}
			if i == len(ac.Lines) {
				original = trimAllSuffixLinebreaks(original)
			}
			result = strings.Replace(result, original, tmp, 1)
			inList = false
			inIndentedBlock = false
			i = i - 1
		} else if strings.HasPrefix(current, hash) {
			// Caption
			current = formatCaption(current)
			result = strings.Replace(result, original, current, 1)
		}
	}

	return result
}

type AFunctionComment string

func (afc AFunctionComment) String() string {
	if afc != "" {
		return formattedComment(string(afc))
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

func (s AStruct) comment() AComment {
	return NewAComment(s.Comment)
}

func (s AStruct) String() string {
	var commentString string
	if s.Comment != "" {
		commentString = s.comment().String()
	}
	var fieldsString string
	for _, f := range s.AFields() {
		fieldsString += f.String()
	}

	if fieldsString == "" {
		fieldsString = fmt.Sprintf("%s%s%s", enclosingBrackets(square, info),
			enclose(hash, indent(filteredFieldsNotice, 2)), preservedLinebreak)
	}

	return fmt.Sprintf("%s%s", codeBlock(fmt.Sprintf("%s%s%s%s", s.asciidocFormattedSigOpen(),
		fieldsString, s.asciidocFormattedSigClose(), preservedLinebreak)), commentString)
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
		enclose(hash, funcTitlePrefix), ws, fn.RefID().AnchorID(), fn.RefID().Identifier)))
}

type AVariable struct {
	api.Variable
}

func NewAVariable(v api.Variable) AVariable {
	return AVariable{Variable: v}
}

func (v AVariable) String() string {
	var docString string
	if v.Doc != "" {
		docString = NewADoc(v.Doc).String() + preservedLinebreak
	}
	return codeBlock(fmt.Sprintf("%s%s%s%s%s%s%s%s%s%s",
		docString, builtinFormat(varPrefix), ws, nameFormat(v.Name), ws, trimAllSuffixLinebreaks(v.asciidocFormattedType()), ws, passThrough(commentPrefix), ws, v.Comment))
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

func (v AVariables) sort() []AVariable {
	return SortMapValues(v, func(a, b AVariable) bool {
		return a.Name < b.Name
	})
}

type commentStatus int

const (
	uncommented commentStatus = iota
	commented
)

func (v AVariables) String() string {
	var s string
	varMap := make(map[commentStatus]string)

	for _, current := range v.sort() {
		if current.Comment == "" && current.Doc == "" {
			varMap[uncommented] += fmt.Sprintf("%s", current.StringRaw())
		} else {
			varMap[commented] += current.String()
		}
	}

	var noCommentVars, commentedVars string
	if varMap[uncommented] != "" {
		noCommentVars = codeBlock(varMap[uncommented])
	}
	if varMap[commented] != "" {
		commentedVars = varMap[commented]
	}
	s = fmt.Sprintf("%s%s%s", noCommentVars, simpleLinebreak, commentedVars)
	s = strings.Trim(s, simpleLinebreak)
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
	var comment, doc string
	if f.Comment != "" {
		comment = fmt.Sprintf("%s%s%s%s", ws, commentPrefix, ws, f.comment().String())
	}
	if f.Doc != "" {
		if slices.Contains(f.Stereotypes, api.StereotypeProperty) {
			doc = indent(fmt.Sprintf("%s%s%s%s", commentPrefix, ws, f.doc().String(), preservedLinebreak), 2)
		}
	}

	var nameString string
	if slices.Contains(f.Stereotypes, api.StereotypeProperty) {
		nameString = indent(f.name().String(), 2)
	} else {
		nameString = f.name().String()
	}
	s := fmt.Sprintf("%s%s%s%s%s",
		doc,
		nameString,
		whiteSpace,
		trimAllSuffixLinebreaks(f.asciidocFormattedType()),
		comment,
	)

	if slices.Contains(f.Stereotypes, api.StereotypeProperty) {
		s += preservedLinebreak
	}
	return s
}

func (f AField) TypeDescription() ATypeDesc {
	return NewATypeDesc(*f.TypeDesc)
}

func (f AField) comment() AComment {
	return NewAComment(f.Comment)
}

func (f AField) doc() ADoc {
	return NewADoc(f.Doc)
}

type AFieldName string

func (afn AFieldName) String() string {
	if afn == "" {
		return ""
	}
	return fmt.Sprintf("%s%s", enclosingBrackets(square, variable), enclose(hash, string(afn)))
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
