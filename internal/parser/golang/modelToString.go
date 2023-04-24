package golang

import (
	"fmt"
	"github.com/worldiety/gdoc/internal/api"
	"golang.org/x/exp/slices"
	"strings"
)

func (h AsciiDocHeader) String() string {
	var s string
	for _, a := range h.Attributes {
		s += fmt.Sprintf("%s%s", a, simpleLinebreak)
	}
	s = strings.Trim(s, simpleLinebreak)
	return s
}

func (m AModule) String() string {
	return fmt.Sprintf("%s%s%s", m.title(), simpleLinebreak, m.readme())
}

func (id APackageRefID) String() string {
	return enclosingDoubleBrackets(angle, fmt.Sprintf("%s,%s%s", id.Identifier, ws, id.Identifier))
}

func (p APackage) String() string {
	return fmt.Sprintf("%s%s", p.title(), p.readme())
}

func (fn AFunction) String() string {
	return fmt.Sprintf("%s%s%s%s%s", bold(fn.name()), preservedLinebreak,
		codeBlock(fn.asciidocFormattedSignature()), simpleLinebreak, fn.comment().String())
}

func (af AFunctions) String() string {
	return af.title()
}

func (m AMethod) String() string {

	return fmt.Sprintf("%s%s%s%s%s", bold(m.name()), preservedLinebreak,
		codeBlock(m.asciidocFormattedSignature()), simpleLinebreak, NewAFunction(*m.Function).comment().String())
}

func (ms AMethods) String() string {
	var methodString string
	for _, m := range ms.sort() {
		methodString += m.String()
	}
	return methodString
}

func (v AVariable) String() string {
	var docString string
	if v.Doc != "" {
		docString = NewADoc(v.Doc).String() + preservedLinebreak
	}
	return codeBlock(fmt.Sprintf("%s%s%s%s%s%s%s%s%s%s%s",
		docString, builtinFormat(varPrefix), ws, v.AnchorID(), v.name().String(), ws, trimAllSuffixLinebreaks(v.asciidocFormattedType()), ws, passThrough(commentPrefix), ws, v.Comment))
}

func (v AVariables) String() string {
	var s string
	varMap := make(map[commentStatus]string)

	for _, current := range v.sort() {
		if current.Comment == "" && current.Doc == "" {
			varMap[uncommented] += current.StringRaw()
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
func (name AFieldName) String() string {
	if name == "" {
		return ""
	}
	return fmt.Sprintf("%s%s", enclosingBrackets(square, variable), enclose(hash, string(name)))
}
func (m AMapType) String() string {
	return fmt.Sprintf("AMapType{MapType: %v}", m.MapType)
}
func (r ARefId) String() string {
	return enclosingDoubleBrackets(angle, fmt.Sprintf("%s,%s%s", r.ID(), ws, r.Identifier))
}

func (as AStructs) String() string {
	return as.title()
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

	return fmt.Sprintf("%s%s%s%s%s%s%s%s%s", s.title(), preservedLinebreak, codeBlock(fmt.Sprintf("%s%s%s%s", s.asciidocFormattedSigOpen(),
		fieldsString, s.asciidocFormattedSigClose(), preservedLinebreak)), commentString, s.methods().String(), simpleLinebreak, simpleLinebreak, "'''", simpleLinebreak)
}

func (afc AFunctionComment) String() string {
	if afc != "" {
		return formattedComment(string(afc))
	}
	return ""
}

func (generics AGenerics) String() string {
	var s string
	var sep = comma + ws
	typeMap := make(map[string][]string, 0)
	for _, g := range generics {
		if typeMap[g.TypeDesc.SrcTypeDefinition] == nil {
			typeMap[g.TypeDesc.SrcTypeDefinition] = make([]string, 0)
		}
		typeMap[g.TypeDesc.SrcTypeDefinition] = append(typeMap[g.TypeDesc.SrcTypeDefinition], g.Name)
	}
	for ts, nameList := range typeMap {
		for _, name := range nameList {
			s += nameFormat(name) + sep
		}
		s = strings.TrimSuffix(s, sep)
		s += ws + typeFormat(ts) + sep
	}
	s = strings.TrimSuffix(s, sep)
	if s != "" {
		s = enclosingBrackets(square, s)
	}
	return s
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
		if isIndentedBlock(current) {
			if inList = isList(current); inList {
				// is list -> change godoc list marker (-) for asciidoc list marker (*)
				originalList = append(originalList, current)
				current = strings.Replace(current, hyphen, asterisk, 1)
				tmpList = append(tmpList, current)
				continue
			} else {
				if !inIndentedBlock {
					originalList = []string{}
					tmpList = []string{}
				}
				originalList = append(originalList, original)
				tmpList, inIndentedBlock = handleIndentedBlock(current, tmpList)
			}
		} else if inList || inIndentedBlock {
			var tmp string
			original, tmp = handleBlocksAndLists(inList, inIndentedBlock, originalList, tmpList, original, len(ac.Lines) == i)
			result = strings.Replace(result, original, tmp, 1)
			inList = false
			inIndentedBlock = false
			i = i - 1
		} else if isCaption(current) {
			// Caption
			current = formatCaption(current)
			result = strings.Replace(result, original, current, 1)
		}
	}

	return result
}
