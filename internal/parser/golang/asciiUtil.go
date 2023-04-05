package golang

import (
	"fmt"
	"strings"
)

type asciidocUtils string

const (
	preservedLinebreak = " +\n"
	simpleLinebreak    = "\n"
	codeBlockDelimiter = "****"
	codeBlockName      = "[.code]"
	passPrefix         = "pass:"
)

func title(prefix, name, anchor string, n int) string {
	if prefix != "" {
		prefix = fmt.Sprintf(" %s ", prefix)
	}

	return fmt.Sprintf("%s%s%s%s%s", simpleLinebreak, lvl(n), prefix, anchor, name)
}

func addComma(s string) string {
	return fmt.Sprintf("%s, ", s)
}

func formattedComment(s string, function bool) string {
	if function {
		return fmt.Sprintf("[%s]#%s#%s", functionComment, strings.Trim(s, "\n"), preservedLinebreak)
	}
	return fmt.Sprintf("[%s]#// %s#%s", comment, s, preservedLinebreak)
}

func operatorFormat(s string) string {
	return fmt.Sprintf("[%s]#%s#", operator, s)
}

func lvl(n int) (lvlStr string) {
	for i := 0; i < n; i++ {
		lvlStr += "="
	}
	return
}

func readme(s string, n int) string {
	return fmt.Sprintf("\n\n%s %s\n%s", lvl(n), readmeTitle, s)
}

func codeBlock(s string) string {
	return fmt.Sprintf("%s\n%s\n%s\n%s", codeBlockName, codeBlockDelimiter, s, codeBlockDelimiter)
}

func passThrough(s string) string {
	return fmt.Sprintf("%s[%s]", passPrefix, s)
}
