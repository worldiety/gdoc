package api

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"
)

type RefId struct {
	ImportPath string
	Identifier string
}

func (id RefId) ID() string {
	tmp := sha256.Sum224([]byte(id.ImportPath + id.Identifier))
	return hex.EncodeToString(tmp[:])
}

func (id RefId) Named() bool {
	return id.Identifier == ""
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

type Import string

type Imports []Import

type BaseType string

type ImportPath = string
type Module struct {
	Readme   string
	Name     string
	Packages map[ImportPath]*Package
}

type Package struct {
	Readme      string
	Doc         string
	Name        string
	Imports     Imports
	Stereotypes []Stereotype
	Types       map[int]BaseType
	Consts      map[string]*Constant
	Vars        map[string]*Variable
	Functions   map[string]*Function
}
type Comment string
type Struct struct {
	ID         RefId
	Comment    Comment
	Name       string
	Fields     []*Field
	Methods    []*Method
	Interfaces []*Interface
}

type Function struct {
	ID         RefId
	Name       string
	Comment    string
	Signature  string
	Parameters map[string]*Parameter
	Results    map[string]*Parameter
}

type Method struct {
	Function
	Receiver *Struct
}

type Constructor struct {
	ID         RefId
	Comment    string
	parameters []*Parameter
}

type Bla struct {
	Yolo func(x int)
	*Bla
}
type Field struct {
	Name              string
	Comment           string
	TypeDefinition    RefId
	SrcTypeDefinition string
	Stereotypes       []Stereotype
}

type Interface struct {
	ID        RefId
	Name      string
	Comment   string
	Methods   []*Method
	Functions []*Function
}
type Constant Field
type Variable Field
type Enum Field
type Example struct {
	ID    RefId
	Value string
}
type Parameter Field

func (p Imports) Len() int           { return len(p) }
func (p Imports) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p Imports) Less(i, j int) bool { return p[i] < p[j] }

func (p Imports) Sort() {
	sort.Sort(p)
}
