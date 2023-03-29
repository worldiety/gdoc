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
