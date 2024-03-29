package golang

import (
	"fmt"
	"strings"
)

const (
	preservedLinebreak = " +\n"
	plusSuffix         = " +"
	simpleLinebreak    = "\n"
	codeBlockDelimiter = "****"
	codeBlockName      = "[.code]"
	mono               = "[mono]"
	passPrefix         = "pass:"
	commentPrefix      = "//"
	hash               = "#"
	boldDelimiter      = "**"
	italicDelimiter    = "__"
	ws                 = " "
	tab                = "\t"
	dot                = "."
	hyphen             = "-"
	asterisk           = "*"
	comma              = ","
	equals             = "="
	nbsp               = "{nbsp}"
)

type bracketType int

const (
	square bracketType = iota
	round
	curly
	angle
)

func title(prefix, anchor, name string, n int) string {
	if prefix != "" {
		prefix = enclose(ws, prefix)
	}

	return fmt.Sprintf("%s%s%s%s%s", simpleLinebreak, lvl(n), prefix, anchor, name)
}

func addComma(s string) string {
	return fmt.Sprintf("%s%s%s", s, comma, ws)
}

func formattedComment(s string) string {

	return fmt.Sprintf("%s%s%s", simpleLinebreak, strings.Trim(s, simpleLinebreak), preservedLinebreak)
}

func operatorFormat(s string) string {
	return fmt.Sprintf("%s%s", enclosingBrackets(square, operator), enclose(hash, s))
}

func readme(s string, n int) string {
	return fmt.Sprintf("%s%s%s%s%s%s", simpleLinebreaks(2), lvl(n), ws, bold(italic(readmeTitle)), simpleLinebreak, s)
}

func codeBlock(s string) string {
	return fmt.Sprintf("%s%s%s%s%s%s%s%s%s",
		simpleLinebreak, codeBlockName, simpleLinebreak, codeBlockDelimiter,
		simpleLinebreak, s, simpleLinebreak, codeBlockDelimiter, simpleLinebreak)
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
	return fmt.Sprintf("%s%s", enclosingBrackets(square, builtin), enclose(hash, s))
}

func typeFormat(s string) string {
	return fmt.Sprintf("%s%s", enclosingBrackets(square, typ3), enclose(hash, s))
}

func nameFormat(s string) string {
	return fmt.Sprintf("%s%s", enclosingBrackets(square, nam3), enclose(hash, s))
}

func keywordFormat(s string) string {
	return fmt.Sprintf("%s%s", enclosingBrackets(square, keyword), enclose(hash, s))
}

func variableFormat(s string) string {
	return fmt.Sprintf("%s%s", enclosingBrackets(square, variable), enclose(hash, s))
}
func formatCaption(s string) string {
	s = strings.Trim(strings.Replace(s, hash, "", -1), ws)
	return lvl(4) + ws + s
}

func formatBlock(tmpList ...string) string {
	var result string
	for _, s := range tmpList {
		result += s + simpleLinebreak
	}
	return result
}

func startsWithEitherPrefix(s string, coll ...string) bool {
	if len(coll) == 0 {
		return true
	}
	for _, pre := range coll {
		if strings.HasPrefix(s, pre) {
			return true
		}
	}
	return false
}

func endsWithEitherSuffix(s string, coll ...string) bool {
	if len(coll) == 0 {
		return true
	}
	for _, suf := range coll {
		if strings.HasSuffix(s, suf) {
			return true
		}
	}
	return false
}

func trimAllPrefixWSAndTabs(s string) string {
	for startsWithEitherPrefix(s, ws, tab) {
		s = strings.TrimPrefix(s, tab)
		s = strings.TrimPrefix(s, ws)
	}
	return s
}

func trimAllSuffixLinebreaks(s string) string {
	for endsWithEitherSuffix(s, simpleLinebreak, plusSuffix, preservedLinebreak) {
		s = strings.TrimSuffix(s, simpleLinebreak)
		s = strings.TrimSuffix(s, plusSuffix)
		s = strings.TrimSuffix(s, preservedLinebreak)
	}
	return s
}

func indent(s string, n int) string {
	if n > 0 {
		s = " " + s
	}
	for i := 0; i < n; i++ {
		s = nbsp + s
	}
	return s
}

func isIndentedBlock(s string) bool {
	return startsWithEitherPrefix(s, ws, tab)
}

func isList(s string) bool {
	return strings.HasPrefix(strings.Trim(s, ws), hyphen)
}

func isCaption(s string) bool {
	return strings.HasPrefix(s, hash)
}

func handleBlocksAndLists(
	inList, inIndentedBlock bool,
	originalList, tmpList []string,
	original string, last bool) (string, string) {
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
	if last {
		original = trimAllSuffixLinebreaks(original)
	}

	return original, tmp
}

func handleIndentedBlock(current string, tmpList []string) ([]string, bool) {
	var count int
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
		break
	}
	return tmpList, count > 0
}
