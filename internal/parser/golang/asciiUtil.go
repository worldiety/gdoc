package golang

import (
	"fmt"
)

const (
	linebreak = " +\n"
)

func title(prefix, name, anchor string, n int) string {
	if prefix != "" {
		prefix = fmt.Sprintf(" %s ", prefix)
	}

	return fmt.Sprintf("\n%s%s%s%s", lvl(n), prefix, anchor, name)
}

func addComma(s string) string {
	return fmt.Sprintf("%s, ", s)
}

func formattedComment(s string) string {
	return fmt.Sprintf("[%s]#// %s#%s", comment, s, linebreak)
}

func operatorFormat(s string) string {
	return fmt.Sprintf("[%s]#%s# ", operator, s)
}

func lvl(n int) (lvlStr string) {
	for i := 0; i < n; i++ {
		lvlStr += "="
	}
	return
}

func readme(s string, n int) string {
	return fmt.Sprintf("\n%s %s\n%s", lvl(n), readmeTitle, s)
}

func codeBlock(s string) string {
	return fmt.Sprintf("[.code]\n****\n%s\n****", s)
}

func (td ATypeDesc) localCustomTypeLink() string {
	return fmt.Sprintf("<<%s, [%s]#%s%s#>>", td.TypeDefinition.ID(), typ3, td.Prefix(), td.Identifier())
}

func (td ATypeDesc) externalCustomTypeLink() string {
	// custom type from external package from this project
	return fmt.Sprintf("<<%s, [%s]#%s%s#>>.<<%s, [%s]#%s#>>",
		// remove the asterisk to find the linked id, it's still displayed in the doc
		td.PkgName(), typ3, td.Prefix(), td.PkgName(), td.TypeDefinition.ID(), typ3, td.Identifier())
}

func (td ATypeDesc) externalNonCustomTypeLink() string {
	return fmt.Sprintf("[%s]#%s#.[%s]#%s#", typ3, td.PkgName(), typ3, td.Identifier())
}

func (td ATypeDesc) builtInTypeLink() string {
	return fmt.Sprintf("[%s]#%s#", builtin, td.SrcTypeDefinition)
}
