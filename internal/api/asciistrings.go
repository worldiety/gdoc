package api

import (
	"fmt"
	"regexp"
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

func (m Module) AnchorID() string {
	return fmt.Sprintf("[[%s]]", m.Name)
}

func (p Package) AnchorID() string {
	return fmt.Sprintf("[[%s]]", p.Name)
}

func (s Struct) AsciidocFormattedComment() string {
	return formattedComment(s.Comment)
}

func (s Struct) AsciidocFormattedSigOpen() string {
	return fmt.Sprintf("[%s]#%s# [[%s]][%s]#%s# [%s]#struct# %s",
		keyword, typ3, s.TypeDefinition.ID(), str1ng, s.Name, keyword, operatorFormat("{"))
}

func (s Struct) AsciidocFormattedSigClose() string {
	return fmt.Sprintf("%s", operatorFormat("}"))
}

func (f Field) AsciidocFormattedComment() string {

	return formattedComment(f.Comment)
}

func (f Field) AsciidocFormattedName() string {
	return fmt.Sprintf("[%s]#%s#", t3xt, f.Name)
}

// AsciidocWhiteSpaceBetween adds non-breaking spaces between a fields name and type, to correctly format code blocks in monospace font
func (f Field) AsciidocWhiteSpaceBetween() string {
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
func (f Field) AsciidocFormattedType() string {
	return formatType(f.SrcTypeDefinition, f.TypeDefinition.ID(), f.Link)
}

// AsciidocFormattedMapType formats the key and value types of a map to asciidoc format
func (f Field) AsciidocFormattedMapType() string {
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

// formatType formats struct and array types for asciidoc
func formatType(srcTypeDef, refId string, link bool) string {
	var replacement string
	originalString := srcTypeDef
	srcTypeDef = withoutBrackets(withoutAsterisks(srcTypeDef))

	if strings.Contains(srcTypeDef, ".") {
		parts := strings.Split(srcTypeDef, ".")
		if link {
			// custom type from external package from this project
			replacement = fmt.Sprintf("<<%s, [%s]#%s#>>.<<%s, [%s]#%s#>>",
				// remove the asterisk to find the linked id, it's still displayed in the doc
				withoutAsterisks(parts[0]), typ3, parts[0], refId, typ3, parts[1])
		} else {
			// from external package, but not from this project
			replacement = fmt.Sprintf("[%s]#%s#.[%s]#%s#", typ3, parts[0], typ3, parts[1])
		}
	} else if link {
		// custom type from current package
		replacement = fmt.Sprintf("<<%s, [%s]#%s#>>", refId, typ3, srcTypeDef)
	} else {
		// if not from external package and not in this package, it's a built-in type
		replacement = fmt.Sprintf("[%s]#%s#", builtin, srcTypeDef)
	}
	return strings.Replace(originalString, srcTypeDef, replacement, 1)
}

func (fn Function) AsciidocFormattedComment() string {
	return fmt.Sprintf(formattedComment(fn.Comment))
}

func (fn Function) AsciidocFormattedSignature() string {

	return fmt.Sprintf("[%s]#func# [%s]#%s#(%s) %s",
		keyword, nam3, fn.Name, fn.AsciidocFormattedParameters(), fn.AsciidocFormattedResults())
}

func (fn Function) AsciidocFormattedParameters() string {
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

func (fn Function) AsciidocFormattedResults() string {
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

func formattedComment(s string) string {
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
