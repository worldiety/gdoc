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

// Doc for a const block
const (
	StereotypeConstructor     = "constructor" // the constructor of a struct
	StereotypeMethod          = "method"
	StereotypeSingleton       = "singleton"
	StereotypeEnum            = "enum"
	StereotypeEnumElement     = "enumElement"
	StereotypeDestructor      = "destructor" // the destructor of a struct
	StereotypeExecutable      = "executable"
	StereotypeStruct          = "struct"
	StereotypeClass           = "class"
	StereotypeProperty        = "property"
	StereotypeParameter       = "parameter"
	StereotypeParameterIn     = "in"
	StereotypeParameterOut    = "out"
	StereotypeParameterResult = "result"
	StereotypeGeneric         = "generic"
)

type ImportPath = string

// A Module contains various [Package] Pointers.
// There is also
//   - a readme and
//   - a name
//
// # A caption
//
// Is defined as in markdown.
// And we have also indented stuff like so:
//
//	   Module.Packages can be accessed directly.
//		        and is formatted like pre.
type Module struct {
	Readme   string
	Name     string
	Packages map[ImportPath]*Package
}

type List[T, X any, V Constant] struct {
}

func (l *List[T, X, V]) Add(t T) {

}

type Person struct {
	Firstname, Last string
}

type AddressBookPerson Person

type AddressBookPerson2 struct {
	Firstname, Last string
}

type BetterPerson = Person

var x = struct{ Firstname, Last string }{}

func bla() {
	x = AddressBookPerson{Firstname: StereotypeProperty}
	x = Person{}

	var p Person = Person{}
	p = Person(AddressBookPerson{})
	p = BetterPerson{}
	_ = p
}

type Vector = List[any, any, Constant]
type Test2 func() string

type Package struct {
	PackageDefinition RefId
	Readme            string
	Doc               string
	Name              string
	Imports           Imports
	Stereotypes       []Stereotype
	Types             map[string]RefId
	Consts            []ConstantBlock
	Vars              map[string]*Variable
	Functions         map[string]*Function
	Structs           map[string]*Struct
}
type Struct struct {
	TypeDefinition     RefId
	Comment            string
	Name               string
	Fields             []*Field
	Methods            []*Method
	Generics           Generics
	Constructors       []*Function
	WhiteSpaceInFields int
}

type Generics []*Field
type Method struct {
	*Function
	Recv *Recv
}

type Recv struct {
	*Field
	Name       string
	TypeString string
}

func NewRecv(f *Field, name, ts string) *Recv {
	return &Recv{
		Field:      f,
		Name:       name,
		TypeString: ts,
	}
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
	TypeDesc *TypeDesc
	Name     string
	Comment  string
	Doc      string
	// ParentStruct test
	ParentStruct *Struct // the struct, this field is a property of
	Stereotypes  []Stereotype
}

func NewField(name, comment string, doc string, t *TypeDesc, parent *Struct) *Field {
	return &Field{
		TypeDesc:     t,
		Name:         name,
		Comment:      comment,
		Doc:          doc,
		ParentStruct: parent,
	}
}

type MapType struct {
	KeyType, ValueType *TypeDesc
}

type TypeDesc struct {
	TypeDefinition    RefId
	SrcTypeDefinition string
	Pointer           bool
	Linebreak         bool
	MapType           *MapType
	TypeOrigin        TypeOrigin
}

func NewTypeDesc(ref RefId, srcTypeDef string, pointer bool, mapType *MapType) *TypeDesc {
	return &TypeDesc{
		TypeDefinition:    ref,
		SrcTypeDefinition: srcTypeDef,
		Pointer:           pointer,
		MapType:           mapType,
	}
}

func removeBrackets(str string) string {
	re := regexp.MustCompile(`\[\d*]`)
	return re.ReplaceAllString(str, "")
}

func withoutAsterisks(s string) string {
	return strings.Replace(s, "*", "", -1)
}

func NewVariable(name, comment, doc string, t *TypeDesc) *Variable {
	return &Variable{
		TypeDesc: t,
		Name:     name,
		Doc:      doc,
		Comment:  comment,
	}
}

type ConstantBlock struct {
	Doc     string
	Content []Constant
}

func NewConstantBlock(consts []Constant, doc string) ConstantBlock {
	return ConstantBlock{
		Doc:     doc,
		Content: consts,
	}
}

type Constant struct {
	RefId   RefId
	Value   any
	Comment string
}

func NewConstant(refId RefId, comment string, value any) Constant {
	return Constant{
		RefId:   refId,
		Value:   value,
		Comment: comment,
	}
}

type Variable Field
type Import string
type Imports []Import

func TestFunc(s string) string {
	return fmt.Sprintf("Hello %s", s)
}
