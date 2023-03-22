package api

import "fmt"

func (m Module) AnchorID() string {
	return fmt.Sprintf("[[%s]]", m.Name)
}

func (p Package) AnchorID() string {
	return fmt.Sprintf("[[%s]]", p.Name)
}

func (v Variable) AnchorID() string {
	return fmt.Sprintf("[[%s]]", v.Name)
}

func (f Field) StringWithLinks() string {
	return fmt.Sprintf("<<%s,%s>>", f.TypeDefinition.ID(), "")
}
