package golang

import (
	"fmt"
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
	return fmt.Sprintf("[%s]#// %s# +", comment, s)
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
