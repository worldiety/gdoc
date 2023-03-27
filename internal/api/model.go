package api

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

type Link int

const (
	None Link = iota
	MapKey
	MapValue
	MapAll
	ArrayType
	FieldType
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
	Parameters     map[string]*Parameter
	Results        map[string]*Parameter
}

type Field struct {
	Name              string
	Comment           string
	ParentStruct      *Struct // the struct, this field is a property of
	MapType           *MapType
	Link              Link
	TypeDefinition    RefId
	SrcTypeDefinition string
	Stereotypes       []Stereotype
}

type MapType struct {
	KeyType, ValueType TypeDesc
}

type TypeDesc struct {
	TypeDefinition    RefId
	SrcTypeDefinition string
	Link              Link
}

type Constant Field
type Variable Field
type Parameter Field
type Import string
type Imports []Import
