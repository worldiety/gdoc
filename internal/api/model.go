package api

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
)

type TypeOrigin int

const (
	BuiltIn TypeOrigin = iota
	LocalCustom
	ExternalCustom
	ExternalNonCustom
)

type RefId struct {
	ImportPath ImportPath
	Identifier string
}

func NewRefID(importPath, identifier string) RefId {
	return RefId{
		ImportPath: importPath,
		Identifier: identifier,
	}
}

func (id RefId) ID() string {
	tmp := sha256.Sum224([]byte(id.ImportPath + id.Identifier))
	return fmt.Sprintf("gd%s", hex.EncodeToString(tmp[:]))
}

func (id RefId) Named() bool {
	return id.Identifier != ""
}

func (id RefId) PackageName() string {
	lastSlashIdx := strings.LastIndex(id.ImportPath, "/")
	return id.ImportPath[lastSlashIdx+1:]
}

// Stereotype as usually interpreted in found context but not expressed in language explicitly.
type Stereotype string

const (
	StereotypeConstructor     = "constructor"
	StereotypeMethod          = "method"
	StereotypeSingleton       = "singleton"
	StereotypeEnum            = "enum"
	StereotypeEnumElement     = "enumElement"
	StereotypeDestructor      = "destructor"
	StereotypeExecutable      = "executable"
	StereotypeStruct          = "struct"
	StereotypeClass           = "class"
	StereotypeProperty        = "property"
	StereotypeParameter       = "parameter"
	StereotypeParameterIn     = "in"
	StereotypeParameterOut    = "out"
	StereotypeParameterResult = "result"
)

type ImportPath = string
type Module struct {
	Readme   string
	Name     string
	Packages map[ImportPath]*Package
}

type Package struct {
	PackageDefinition RefId
	Readme            string
	Doc               string
	Name              string
	Imports           Imports
	Stereotypes       []Stereotype
	Types             map[string]RefId
	Consts            map[string]*Constant
	Vars              map[string]*Variable
	Functions         map[string]*Function
	Structs           map[string]*Struct
}
type Struct struct {
	TypeDefinition     RefId
	Comment            string
	Name               string
	Fields             []*Field
	WhiteSpaceInFields int
}

type Function struct {
	TypeDefinition RefId
	Name           string
	Comment        string
	Signature      string
	Parameters     map[string]*Field
	Results        map[string]*Field
}

type Field struct {
	TypeDesc     *TypeDesc
	Name         string
	Comment      string
	ParentStruct *Struct // the struct, this field is a property of
	MapType      *MapType
	Stereotypes  []Stereotype
}

type MapType struct {
	KeyType, ValueType *TypeDesc
}

type TypeDesc struct {
	TypeDefinition    RefId
	SrcTypeDefinition string
	Pointer           bool
	Link              bool
	TypeOrigin        TypeOrigin
}

func removeBrackets(str string) string {
	re := regexp.MustCompile(`\[\d*]`)
	return re.ReplaceAllString(str, "")
}

func withoutAsterisks(s string) string {
	return strings.Replace(s, "*", "", -1)
}

type Constant Field
type Variable Field
type Import string
type Imports []Import
