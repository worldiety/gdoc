package golang

import (
	"fmt"
	"github.com/worldiety/gdoc/internal/api"
	"log"
	"strings"
)

// syntax highlighting for Asciidoc
const (
	keyword         = "keyword"
	t3xt            = "text"
	background      = "background"
	builtin         = "builtin"
	str1ng          = "string"
	numb3r          = "number"
	comment         = "comment"
	functionDecl    = "functionDecl"
	functionCall    = "functionCall"
	variable        = "variable"
	Const           = "constant"
	operator        = "operator"
	control         = "control"
	preprocessor    = "preprocessor"
	other           = "other"
	typ3            = "type"
	nam3            = "name"
	info            = "information"
	code            = "code"
	functionComment = "functionComment"
	funcTitle       = "func"
	structTitle     = "struct"
	mapPrefix       = "map"
)

func (m AModule) AnchorID() string {
	return enclosingDoubleBrackets(square, m.Name)
}

func (m AModule) title() string {
	return title(moduleTitlePrefix, m.Name, "", 1)
}

func (m AModule) readme() string {
	if m.Readme != "" {
		return readme(m.Readme, 2)
	}
	return ""
}

func (p APackage) AnchorID() string {
	return enclosingDoubleBrackets(square, p.Name)
}

func (p APackage) title() string {
	return title(packageTitlePrefix, p.AnchorID(), p.Name, 2)
}

func (p APackage) readme() string {
	if p.Readme != "" {
		return readme(p.Readme, 3)
	}
	return ""
}

func (s AStruct) asciidocFormattedComment() string {
	return s.comment().String()
}

func (s AStruct) comment() AComment {
	return AComment(s.Comment)
}

func (s AStruct) asciidocFormattedSigOpen() string {
	return fmt.Sprintf("%s%s%s%s%s%s%s%s%s%s%s%s",
		enclosingBrackets(square, keyword), enclose(formatDelimiter, typ3), ws, enclosingDoubleBrackets(square, s.TypeDefinition.ID()),
		enclosingBrackets(square, str1ng), enclose(formatDelimiter, s.Name), ws, enclosingBrackets(square, keyword),
		enclose(formatDelimiter, structTitle), ws, operatorFormat("{"), preservedLinebreak)
}

func (s AStruct) asciidocFormattedSigClose() string {
	return fmt.Sprintf("%s", operatorFormat("}"))
}

// AsciidocWhiteSpaceBetween adds non-breaking spaces between a fields name and type, to correctly format code blocks in monospace font
func (f AField) asciidocWhiteSpaceBetween() string {
	var s string
	if f.ParentStruct != nil {
		n := f.ParentStruct.WhiteSpaceInFields - len([]rune(f.Name))
		for i := 0; i < n; i++ {
			s += nbsp // nbsp element is necessary, " " would not be preserved
		}
	}
	return s + ws // at the end a white space has to be added like this or any following formatting code like [code]#*#, will lose its effect
}

// AsciidocFormattedType formats a fields' SrcTypeDefinition.
// It adds links:
// 1. to the package, if the origin package is not the current one (optional). Packages use their names as ID.
// 2. to the fields' type. Fields have a prefixed Hex encoded id embedded in the Asciidoc.
func (f AField) asciidocFormattedType() string {
	var s string
	if f.TypeDesc.Map() {
		s = f.asciidocFormattedMapType()
	} else {
		s = f.TypeDescription().typeString()
	}
	if f.TypeDesc.Linebreak {
		s = addLinebreak(s)
	}
	return s
}

func addLinebreak(s string) string {
	return s + preservedLinebreak
}
func (td ATypeDesc) typeString() string {
	var result string

	switch td.TypeOrigin {
	case api.LocalCustom:
		result += td.localCustomTypeLink()
	case api.ExternalCustom:
		result += td.externalCustomTypeLink()
	case api.ExternalNonCustom:
		result += td.externalNonCustomTypeLink()
	case api.BuiltIn:
		result += td.builtInTypeLink()
	default:
		log.Fatal("unknown enum while formatting type")
	}

	return result
}
func (fn AFunction) asciidocFormattedSignature() string {

	return fmt.Sprintf("%s%s%s%s%s%s%s%s",
		enclosingBrackets(square, keyword), enclose(formatDelimiter, funcTitle), ws, enclosingBrackets(square, nam3),
		enclose(formatDelimiter, fn.Name), enclosingBrackets(round, fn.asciidocFormattedParameters()), ws, fn.asciidocFormattedResults())
}

func (fn AFunction) asciidocFormattedParameters() string {
	var s string
	var c int
	for _, p := range fn.Parameters {
		s += fmt.Sprintf("%s", NewAField(*p).String())
		if c < len(fn.Parameters)-1 {
			s = addComma(s)
		}
		c++
	}
	return s
}

func (fn AFunction) asciidocFormattedResults() string {
	var results string
	var s string
	var c int
	for _, r := range fn.Results {
		results += fmt.Sprintf("%s", NewAField(*r).String())
		if c < len(fn.Results)-1 {
			results = addComma(results)
		}
		c++
	}
	if len(fn.Results) > 1 {
		s = enclosingBrackets(round, results)
	}

	return s
}

// AsciidocFormattedMapType formats the key and value types of a map to asciidoc format
func (f AField) asciidocFormattedMapType() string {
	srcTypeDef := f.TypeDesc.SrcTypeDefinition
	keyType := *f.TypeDesc.MapType.KeyType
	valueType := *f.TypeDesc.MapType.ValueType
	formattedKeySrcTypeDef := NewATypeDesc(keyType).typeString()
	formattedValueSrcTypeDef := NewATypeDesc(valueType).typeString()

	srcTypeDef = strings.Replace(srcTypeDef, keyType.SrcTypeDefinition, formattedKeySrcTypeDef, 1)
	srcTypeDef = strings.Replace(srcTypeDef, valueType.SrcTypeDefinition, formattedValueSrcTypeDef, 1)
	srcTypeDef = strings.Replace(srcTypeDef, mapPrefix, fmt.Sprintf("%s%s", enclosingBrackets(square, keyword),
		enclose(formatDelimiter, mapPrefix)), 1)

	return srcTypeDef
}
