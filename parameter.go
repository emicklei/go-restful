package restful

// Copyright 2013 Ernest Micklei. All rights reserved.
// Use of this source code is governed by a license
// that can be found in the LICENSE file.

const (
	PATH_PARAMETER = iota // indicator of Request parameter type
	QUERY_PARAMETER
	BODY_PARAMETER
	HEADER_PARAMETER
)

// Parameter is for documententing the parameter used in a Http Request
// ParameterData kinds are Path,Query and Body
type Parameter struct {
	data *ParameterData
}

// ParameterData represents the state of a Parameter.
// It is made public to make it accessible to e.g. the Swagger package.
type ParameterData struct {
	Name, Description, DataType string
	Kind                        int
	Required                    bool
	AllowableValues             map[string]string
	AllowMultiple               bool
}

// Data returns the state of the Parameter
func (p *Parameter) Data() ParameterData {
	return *p.data
}

// Kind returns the parameter type indicator (see const for valid values)
func (p *Parameter) Kind() int {
	return p.data.Kind
}

func (p *Parameter) bePath() *Parameter {
	p.data.Kind = PATH_PARAMETER
	return p
}
func (p *Parameter) beQuery() *Parameter {
	p.data.Kind = QUERY_PARAMETER
	return p
}
func (p *Parameter) beBody() *Parameter {
	p.data.Kind = BODY_PARAMETER
	return p
}

func (p *Parameter) beHeader() *Parameter {
	p.data.Kind = HEADER_PARAMETER
	return p
}

// Required sets the required field and return the receiver
func (p *Parameter) Required(required bool) *Parameter {
	p.data.Required = required
	return p
}

// AllowMultiple sets the allowMultiple field and return the receiver
func (p *Parameter) AllowMultiple(multiple bool) *Parameter {
	p.data.AllowMultiple = multiple
	return p
}

// AllowableValues sets the allowableValues field and return the receiver
func (p *Parameter) AllowableValues(values map[string]string) *Parameter {
	p.data.AllowableValues = values
	return p
}

// DataType sets the dataType field and return the receiver
func (p *Parameter) DataType(typeName string) *Parameter {
	p.data.DataType = typeName
	return p
}
