package swagger

import (
	"bytes"
	"encoding/json"
)

// NamedModelProperty associates a name to a ModelProperty
type NamedModelProperty struct {
	Name     string
	Property ModelProperty
}

// ModelPropertyList encapsulates a list of NamedModelProperty (association)
type ModelPropertyList struct {
	List []NamedModelProperty
}

// NewModelPropertyList returns a new empty ModelPropertyList
func NewModelPropertyList() *ModelPropertyList {
	return &ModelPropertyList{[]NamedModelProperty{}}
}

// MarshalJSON writes the ModelPropertyList as if it was a map[string]ModelProperty
func (l ModelPropertyList) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString("{\n")
	for i, each := range l.List {
		buf.WriteString("\"")
		buf.WriteString(each.Name)
		buf.WriteString("\": ")
		json.NewEncoder(&buf).Encode(each.Property)
		if i < len(l.List)-1 {
			buf.WriteString(",\n")
		}
	}
	buf.WriteString("}")
	return buf.Bytes(), nil
}

// At returns the ModelPropety by its name unless absent, then ok is false
func (l *ModelPropertyList) At(name string) (p ModelProperty, ok bool) {
	for _, each := range l.List {
		if each.Name == name {
			return each.Property, true
		}
	}
	return p, false
}

// Put add or replaces a ModelProperty with this name
func (l *ModelPropertyList) Put(name string, prop ModelProperty) {
	// maybe replace existing
	for i, each := range l.List {
		if each.Name == name {
			// replace
			l.List[i] = NamedModelProperty{Name: name, Property: prop}
			return
		}
	}
	// add
	l.List = append(l.List, NamedModelProperty{Name: name, Property: prop})
}

// Do enumerates all the properties, each with its assigned name
func (l *ModelPropertyList) Do(block func(name string, value ModelProperty)) {
	for _, each := range l.List {
		block(each.Name, each.Property)
	}
}
