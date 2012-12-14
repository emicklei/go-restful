package restful

import ()

type ParameterKind int

const (
	PATH ParameterKind = iota
	QUERY
	BODY
)

type Parameter struct {
	Name, Description string
	Kind              ParameterKind
}

type ParameterBuilder struct {
	parameter *Parameter
}

func (self *ParameterBuilder) Name(name string) *ParameterBuilder {
	self.parameter.Name = name
	return self
}
func (self *ParameterBuilder) Description(desc string) *ParameterBuilder {
	self.parameter.Description = desc
	return self
}
func (self *ParameterBuilder) Kind(kind ParameterKind) *ParameterBuilder {
	self.parameter.Kind = kind
	return self
}
func (self *ParameterBuilder) Build() Parameter {
	return *self.parameter
}
