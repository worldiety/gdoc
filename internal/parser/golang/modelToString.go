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
	return fmt.Sprintf("%s", enclosingDoubleBrackets(angle, fmt.Sprintf("%s,%s%s", id.Identifier, ws, id.Identifier)))
}

func (p APackage) String() string {
	return fmt.Sprintf("%s%s", p.title(), p.readme())
}

func (fn AFunction) String() string {
	return fmt.Sprintf("%s%s%s%s%s", fn.title(), preservedLinebreak,
		codeBlock(fn.asciidocFormattedSignature()), simpleLinebreak, fn.comment().String())
}

func (af AFunctions) String() string {
	return af.title()
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
	return fmt.Sprintf("%s", enclosingDoubleBrackets(angle, fmt.Sprintf("%s,%s%s", r.ID(), ws, r.Identifier)))
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

	return fmt.Sprintf("%s%s", codeBlock(fmt.Sprintf("%s%s%s%s", s.asciidocFormattedSigOpen(),
		fieldsString, s.asciidocFormattedSigClose(), preservedLinebreak)), commentString)
}

func (afc AFunctionComment) String() string {
	if afc != "" {
		return formattedComment(string(afc))
	}
	return ""
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
