package api

import (
	"fmt"
	"regexp"
	"strings"
)

// syntax highlighting for asciidoc
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

func (m Module) AnchorID() string {
	return fmt.Sprintf("[[%s]]", m.Name)
}

func (p Package) AnchorID() string {
	return fmt.Sprintf("[[%s]]", p.Name)
}

func (s Struct) FormattedComment() string {
	return FormattedComment(s.Comment)
}

func (s Struct) FormattedSigOpen() string {
	return fmt.Sprintf("[%s]#%s# [[%s]][%s]#%s# [%s]#struct# %s",
		keyword, typ3, s.TypeDefinition.ID(), str1ng, s.Name, keyword, operatorFormat("{"))
}

func (s Struct) FormattedSigClose() string {
	return fmt.Sprintf("%s", operatorFormat("}"))
}

func (f Field) FormattedComment() string {

	return FormattedComment(f.Comment)
}

func (f Field) FormattedName() string {
	return fmt.Sprintf("[%s]#%s#", t3xt, f.Name)
}

// WhiteSpaceBetween adds non-breaking spaces between a fields name and type, to correctly format code blocks in monospace font
func (f Field) WhiteSpaceBetween() string {
	var s string
	if f.ParentStruct != nil {
		n := f.ParentStruct.WhiteSpaceInFields - len([]rune(f.Name))
		for i := 0; i < n; i++ {
			s += "{nbsp}" // nbsp element is necessary, " " would not be preserved
		}
	}
	return s + " " // at the end a white space has to be added like this or any following formatting code like [code]#*#, will lose its effect
}

// FormattedType formats a fields' SrcTypeDefinition.
// It adds links:
// 1. to the package, if the origin package is not the current one (optional). Packages use their names as ID.
// 2. to the fields' type. Fields have a prefixed Hex encoded id embedded in the asciidoc.
func (f Field) FormattedType() string {
	return formatType(f.SrcTypeDefinition, f.TypeDefinition.ID(), f.Link)
}

func (f Field) FormattedMapType() string {
	srcTypeDef := f.SrcTypeDefinition
	keyType := f.MapType.KeyType
	valueType := f.MapType.ValueType
	formattedKeySrcTypeDef := formatType(keyType.SrcTypeDefinition, keyType.TypeDefinition.ID(), keyType.Link)
	formattedValueSrcTypeDef := formatType(valueType.SrcTypeDefinition, valueType.TypeDefinition.ID(), valueType.Link)

	srcTypeDef = strings.Replace(srcTypeDef, keyType.SrcTypeDefinition, formattedKeySrcTypeDef, 1)
	srcTypeDef = strings.Replace(srcTypeDef, valueType.SrcTypeDefinition, formattedValueSrcTypeDef, 1)
	srcTypeDef = strings.Replace(srcTypeDef, "map", fmt.Sprintf("[%s]"+"#map#", keyword), 1)

	return srcTypeDef
}

func formatType(srcTypeDef, refId string, link bool) string {
	var replacement string
	originalString := srcTypeDef
	srcTypeDef = withoutBrackets(withoutAsterisks(srcTypeDef))

	if strings.Contains(srcTypeDef, ".") {
		parts := strings.Split(srcTypeDef, ".")
		if link {
			replacement = fmt.Sprintf("<<%s, [%s]#%s#>>.<<%s, [%s]#%s#>>",
				// remove the asterisk to find the linked id, it's still displayed in the doc
				withoutAsterisks(parts[0]), typ3, parts[0], refId, typ3, parts[1])
		} else {
			replacement = fmt.Sprintf("[%s]#%s#.[%s]#%s#", typ3, parts[0], typ3, parts[1])
		}
	} else if link {
		replacement = fmt.Sprintf("<<%s, [%s]#%s#>>", refId, typ3, srcTypeDef)
	} else {
		replacement = fmt.Sprintf("[%s]#%s#", builtin, srcTypeDef)
	}
	return strings.Replace(originalString, srcTypeDef, replacement, 1)
}

func (fn Function) FormattedComment() string {
	return fmt.Sprintf(FormattedComment(fn.Comment))
}

func (fn Function) FormattedSignature() string {

	return fmt.Sprintf("[%s]#func# [%s]#%s#(%s) %s",
		keyword, nam3, fn.Name, fn.FormattedParameters(), fn.FormattedResults())
}

func (fn Function) FormattedParameters() string {
	var s string
	var c int
	for _, p := range fn.Parameters {
		s += fmt.Sprintf("[%s]#%s# [%s]#%s#", variable, p.Name, typ3, p.SrcTypeDefinition)
		if c < len(fn.Parameters)-1 {
			s = addComma(s)
		}
		c++
	}
	return s
}

func (fn Function) FormattedResults() string {
	var results string
	var s string
	var c int
	for _, r := range fn.Results {
		results += fmt.Sprintf("[%s]#%s#", variable, r.SrcTypeDefinition)
		if c < len(fn.Results)-1 {
			results = addComma(results)
		}
		c++
	}
	if len(fn.Results) > 1 {
		s = fmt.Sprintf("(%s)", results)
	}

	return s
}

func addComma(s string) string {
	return fmt.Sprintf("%s, ", s)
}

func withoutAsterisks(s string) string {
	return strings.Replace(s, "*", "", -1)
}

func FormattedComment(s string) string {
	s = strings.Trim(s, "\n")
	return fmt.Sprintf("[%s]#// %s#", comment, s)
}

func operatorFormat(s string) string {
	return fmt.Sprintf("[%s]#%s# ", operator, s)
}

func withoutBrackets(s string) string {
	re := regexp.MustCompile(`\[\d*\]`)
	return re.ReplaceAllString(s, "")
}
