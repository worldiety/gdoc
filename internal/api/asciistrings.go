package api

import (
	"fmt"
	"strings"
)

func (m Module) AnchorID() string {
	return fmt.Sprintf("[[%s]]", m.Name)
}

func (p Package) AnchorID() string {
	return fmt.Sprintf("[[%s]]", p.Name)
}

func (s Struct) FormattedComment() string {
	return fmt.Sprintf("[comment]#//%s#", string(s.Comment))
}

func (s Struct) FormattedSigOpen() string {
	return fmt.Sprintf("[keyword]#type# [[%s]][string]#%s# [keyword]#struct# [operator]#{#", s.TypeDefinition.ID(), s.Name)
}

func (s Struct) FormattedSigClose() string {
	return fmt.Sprintf("[operator]#}#")
}

func (f Field) FormattedComment() string {

	return fmt.Sprintf("[comment]#// %s#", f.Comment)
}

func (f Field) FormattedName() string {
	return fmt.Sprintf("[text]#%s#", f.Name)
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
	return s + " "
}

// FormattedType formats a fields' SrcTypeDefinition.
// It adds links:
// 1. to the package, if the origin package is not the current one (optional). Packages use their names as ID.
// 2. to the fields' type. Fields have a prefixed Hex encoded id embedded in the asciidoc.
func (f Field) FormattedType() string {
	if strings.Contains(f.SrcTypeDefinition, ".") {
		parts := strings.Split(f.SrcTypeDefinition, ".")
		if f.Link {
			return fmt.Sprintf("<<%s, [type]#%s#>>.<<%s, [type]#%s#>>",
				// remove the asterisk to find the linked id, it's still displayed in the doc
				removeAsterisks(parts[0]), parts[0], f.TypeDefinition.ID(), parts[1])
		}
		return fmt.Sprintf("[type]#%s#.[type]#%s#", parts[0], parts[1])
	}

	if f.Link {
		return fmt.Sprintf("<<%s, [type]#%s#>>", f.TypeDefinition.ID(), f.SrcTypeDefinition)
	}
	return fmt.Sprintf("[type]#%s#", f.SrcTypeDefinition)
}

func removeAsterisks(s string) string {
	return strings.Replace(s, "*", "", -1)
}
