package api

import (
	"regexp"
	"strings"
)

func (td TypeDesc) Map() bool {
	return strings.Contains(td.SrcTypeDefinition, "map[")
}

func (td TypeDesc) Array() bool {
	re := regexp.MustCompile(`\[\d*]`)
	return re.MatchString(td.SrcTypeDefinition)
}

func (td TypeDesc) Identifier() string {
	if include, _ := td.IncludesPkg(); include {
		return withoutAsterisks(removeBrackets(strings.Split(td.SrcTypeDefinition, ".")[1]))
	}
	return removeBrackets(withoutAsterisks(td.SrcTypeDefinition))
}

func (td TypeDesc) IncludesPkg() (bool, string) {
	includes := strings.Contains(td.SrcTypeDefinition, ".")
	var pkgName string
	if includes {
		pkgName = td.PkgName()
	}
	return includes, pkgName
}

func (td TypeDesc) PkgName() string {
	return removeBrackets(withoutAsterisks(strings.Split(td.SrcTypeDefinition, ".")[0]))
}

func (td TypeDesc) Prefix() string {
	var s string
	if td.Array() {
		var prefixEnd int
		if td.lastAsteriskIndex() > td.lastClosedArrayBracketIndex() {
			prefixEnd = td.lastAsteriskIndex()
		} else {
			prefixEnd = td.lastClosedArrayBracketIndex()
		}
		arr := td.SrcTypeDefinition[:prefixEnd+1]
		for i, r := range arr {
			s += string(r)
			if r == '*' && (i+1 < len(arr) && (arr[i+1] == '*' || arr[i+1] == '[')) {
				s += " "
			}
		}

		return s
	} else if td.Pointer {
		return td.SrcTypeDefinition[:td.lastAsteriskIndex()+1]
	}

	return s
}

func (td TypeDesc) lastAsteriskIndex() int {
	return strings.LastIndex(td.SrcTypeDefinition, "*")
}

func (td TypeDesc) lastClosedArrayBracketIndex() int {
	return strings.LastIndex(td.SrcTypeDefinition, "]")
}

func (td TypeDesc) mapParts() []string {
	if td.Map() {
		tmp := strings.Replace(td.SrcTypeDefinition, "map[", "", 1)
		tmpArr := strings.Split(tmp, "]")
		return tmpArr
	}

	return nil
}

func (td TypeDesc) MapSrcDefs() (string, string) {
	if td.Map() {
		return td.mapParts()[0], td.mapParts()[1]
	}
	return "", ""
}
