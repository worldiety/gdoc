package golang

import (
	"fmt"
	"github.com/worldiety/gdoc/internal/api"
	"log"
	"strings"
)

// syntax highlighting for Asciidoc
const (
	keyword      = "keyword"
	t3xt         = "text"
	background   = "background"
	builtin      = "builtin"
	str1ng       = "string"
	numb3r       = "number"
	comment      = "comment"
	functionDecl = "functionDecl"
	functionCall = "functionCall"
	variable     = "variable"
	Const        = "constant"
	operator     = "operator"
	control      = "control"
	preprocessor = "preprocessor"
	other        = "other"
	typ3         = "type"
	nam3         = "name"
	code         = "code"
)

func (m AModule) AnchorID() string {
	return fmt.Sprintf("[[%s]]", m.Name)
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
	return fmt.Sprintf("[[%s]]", p.Name)
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
	return formattedComment(s.Comment)
}

func (s AStruct) asciidocFormattedSigOpen() string {
	return fmt.Sprintf("[%s]#%s# [[%s]][%s]#%s# [%s]#struct# %s +\n",
		keyword, typ3, s.TypeDefinition.ID(), str1ng, s.Name, keyword, operatorFormat("{"))
}

func (s AStruct) asciidocFormattedSigClose() string {
	return fmt.Sprintf("%s", operatorFormat("}"))
}

func (f AField) asciidocFormattedComment() string {
	if f.Comment != "" {
		return formattedComment(f.Comment)
	}
	return ""
}

func (f AField) asciidocFormattedName() string {
	return fmt.Sprintf("[%s]#%s#", t3xt, f.Name)
}

// AsciidocWhiteSpaceBetween adds non-breaking spaces between a fields name and type, to correctly format code blocks in monospace font
func (f AField) asciidocWhiteSpaceBetween() string {
	var s string
	if f.ParentStruct != nil {
		n := f.ParentStruct.WhiteSpaceInFields - len([]rune(f.Name))
		for i := 0; i < n; i++ {
			s += "{nbsp}" // nbsp element is necessary, " " would not be preserved
		}
	}
	return s + " " // at the end a white space has to be added like this or any following formatting code like [code]#*#, will lose its effect
}

// AsciidocFormattedType formats a fields' SrcTypeDefinition.
// It adds links:
// 1. to the package, if the origin package is not the current one (optional). Packages use their names as ID.
// 2. to the fields' type. Fields have a prefixed Hex encoded id embedded in the Asciidoc.
func (f AField) asciidocFormattedType() string {
	var s string
	if f.TypeDesc.Map() {
		s = f.asciidocFormattedMapType() + linebreak
	} else {
		s = f.TypeDescription().typeString() + linebreak
	}
	return s
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
func (fn AFunction) asciidocFormattedComment() string {
	return fmt.Sprintf(formattedComment(fn.Comment))
}

func (fn AFunction) asciidocFormattedSignature() string {

	return fmt.Sprintf("[%s]#func# [%s]#%s#(%s) %s",
		keyword, nam3, fn.Name, fn.asciidocFormattedParameters(), fn.asciidocFormattedResults())
}

func (fn AFunction) asciidocFormattedParameters() string {
	//var s string
	//var c int
	//for _, p := range fn.Parameters {
	//	s += p.TypeDesc
	//	if c < len(fn.Parameters)-1 {
	//		s = addComma(s)
	//	}
	//	c++
	//}
	//return s
	return ""
}

func (fn AFunction) asciidocFormattedResults() string {
	//var results string
	//var s string
	//var c int
	//for _, r := range fn.Results {
	//	results += fmt.Sprintf("[%s]#%s#", variable, r.TypeDesc.SrcTypeDefinition)
	//	if c < len(fn.Results)-1 {
	//		results = addComma(results)
	//	}
	//	c++
	//}
	//if len(fn.Results) > 1 {
	//	s = fmt.Sprintf("(%s)", results)
	//}

	return ""
}

// AsciidocFormattedMapType formats the key and value types of a map to asciidoc format
func (f AField) asciidocFormattedMapType() string {
	srcTypeDef := f.TypeDesc.SrcTypeDefinition
	keyType := f.TypeDesc.MapType.KeyType
	valueType := f.TypeDesc.MapType.ValueType
	formattedKeySrcTypeDef := NewATypeDesc(*keyType).typeString()
	formattedValueSrcTypeDef := NewATypeDesc(*valueType).typeString()

	srcTypeDef = strings.Replace(srcTypeDef, keyType.SrcTypeDefinition, formattedKeySrcTypeDef, 1)
	srcTypeDef = strings.Replace(srcTypeDef, valueType.SrcTypeDefinition, formattedValueSrcTypeDef, 1)
	srcTypeDef = strings.Replace(srcTypeDef, "map", fmt.Sprintf("[%s]"+"#map#", keyword), 1)

	return srcTypeDef
}
