package golang

import (
	"fmt"
	"strings"
)

const (
	preservedLinebreak = " +\n"
	simpleLinebreak    = "\n"
	codeBlockDelimiter = "****"
	codeBlockName      = "[.code]"
	passPrefix         = "pass:"
	commentPrefix      = "//"
	formatDelimiter    = "#"
	boldDelimiter      = "**"
	italicDelimiter    = "__"
	ws                 = " "
	dot                = "."
	nbsp               = "{nbsp}"
)

type bracketType int

const (
	square bracketType = iota
	round
	curly
	angle
)

func title(prefix, name, anchor string, n int) string {
	if prefix != "" {
		prefix = fmt.Sprintf("%s", enclose(" ", prefix))
	}

	return fmt.Sprintf("%s%s%s%s%s", simpleLinebreak, lvl(n), prefix, anchor, name)
}

func addComma(s string) string {
	return fmt.Sprintf("%s,%s", s, ws)
}

func formattedComment(s string, function bool) string {
	if function {
		return fmt.Sprintf("%s%s%s%s", simpleLinebreak, enclosingBrackets(square, comment),
			enclose(formatDelimiter, strings.Trim(s, simpleLinebreak)), preservedLinebreak)
	}
	return fmt.Sprintf("%s%s%s", enclosingBrackets(square, codeBlockComment),
		enclose(formatDelimiter, commentPrefix, ws, s), preservedLinebreak)
}

func operatorFormat(s string) string {
	return fmt.Sprintf("%s%s", enclosingBrackets(square, operator), enclose(formatDelimiter, s))
}

func readme(s string, n int) string {
	return fmt.Sprintf("%s%s%s%s%s%s", simpleLinebreaks(2), lvl(n), ws, bold(italic(readmeTitle)), simpleLinebreak, s)
}

func codeBlock(s string) string {
	return fmt.Sprintf("%s%s%s%s%s%s%s%s",
		simpleLinebreak, codeBlockName, simpleLinebreak, codeBlockDelimiter,
		simpleLinebreak, s, simpleLinebreak, codeBlockDelimiter)
}

func passThrough(s string) string {
	return fmt.Sprintf("%s%s", passPrefix, enclosingBrackets(square, s))
}

func simpleLinebreaks(n int) string {
	var s string
	for i := 0; i < n; i++ {
		s += simpleLinebreak
	}
	return s
}

func enclosingBrackets(bt bracketType, s ...string) string {
	var result string
	js := strings.Join(s, "")
	switch bt {
	case square:
		result = fmt.Sprintf("[%s]", js)
	case angle:
		result = fmt.Sprintf("<%s>", js)
	case curly:
		result = fmt.Sprintf("{%s}", js)
	case round:
		result = fmt.Sprintf("(%s)", js)
	}

	return result
}

func enclosingDoubleBrackets(bt bracketType, s ...string) string {
	var result string
	js := strings.Join(s, "")
	switch bt {
	case square:
		result = fmt.Sprintf("[[%s]]", js)
	case angle:
		result = fmt.Sprintf("<<%s>>", js)
	case curly:
		result = fmt.Sprintf("{{%s}}", js)
	case round:
		result = fmt.Sprintf("((%s))", js)
	}

	return result
}

func enclose(outsideString string, s ...string) string {
	js := strings.Join(s, "")
	return fmt.Sprintf("%s%s%s", outsideString, js, outsideString)
}

func lvl(n int) (lvlStr string) {
	for i := 0; i < n; i++ {
		lvlStr += "="
	}
	return
}

func bold(s ...string) string {
	return enclose(boldDelimiter, s...)
}

func italic(s ...string) string {
	return enclose(italicDelimiter, s...)
}

func builtinFormat(s string) string {
	return fmt.Sprintf("%s%s", enclosingBrackets(square, builtin), enclose(formatDelimiter, s))
}

func typeFormat(s string) string {
	return fmt.Sprintf("%s%s", enclosingBrackets(square, typ3), enclose(formatDelimiter, s))
}

func nameFormat(s string) string {
	return fmt.Sprintf("%s%s", enclosingBrackets(square, nam3), enclose(formatDelimiter, s))
}
