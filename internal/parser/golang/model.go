package golang

import (
	"fmt"
	"github.com/worldiety/gdoc/internal/api"
)

type AModule struct {
	Module api.Module
}

func NewAModule(module api.Module) AModule {
	return AModule{Module: module}
}

func (m AModule) String() string {
	return fmt.Sprintf("AModule{Module: %v}", m.Module)
}

type APackage struct {
	api.Package
}

func NewAPackage(packageVal api.Package) APackage {
	return APackage{Package: packageVal}
}

func (p APackage) String() string {
	return fmt.Sprintf("APackage{Package: %v}", p.Package)
}

type ARefId struct {
	api.RefId
}

func NewARefId(refId api.RefId) ARefId {
	return ARefId{RefId: refId}
}

func (r ARefId) String() string {
	return fmt.Sprintf("ARefId{RefId: %v}", r.RefId)
}

type AStruct struct {
	api.Struct
}

func NewAStruct(structVal api.Struct) AStruct {
	return AStruct{Struct: structVal}
}

func (s AStruct) String() string {
	return fmt.Sprintf("AStruct{Struct: %v}", s.Struct)
}

type AFunction struct {
	api.Function
}

func NewAFunction(functionVal api.Function) AFunction {
	return AFunction{Function: functionVal}
}

func (f AFunction) String() string {
	return fmt.Sprintf("AFunction{Function: %v}", f.Function)
}

type AField struct {
	api.Field
}

func NewAField(fieldVal api.Field) AField {
	return AField{Field: fieldVal}
}

func (f AField) String() string {
	return fmt.Sprintf("AField{Field: %v}", f.Field)
}

type AMapType struct {
	api.MapType
}

func NewAMapType(mapTypeVal api.MapType) AMapType {
	return AMapType{MapType: mapTypeVal}
}

func (m AMapType) String() string {
	return fmt.Sprintf("AMapType{MapType: %v}", m.MapType)
}

type ATypeDesc struct {
	api.TypeDesc
}

func NewATypeDesc(typeDescVal api.TypeDesc) ATypeDesc {
	return ATypeDesc{TypeDesc: typeDescVal}
}

func (td ATypeDesc) RefId() ARefId {
	return NewARefId(td.TypeDefinition)
}

func (td ATypeDesc) String() string {
	var s string
	if td.Pointer {
		s += "*"
	}

	s += td.RefId().String()

	return s
}
