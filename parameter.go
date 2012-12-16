package restful

const (
	PATH_PARAMETER = iota
	QUERY_PARAMETER
	BODY_PARAMETER
)

// Parameter is for documententing the parameter used in a Http Request
// Parameter kinds are Path,Query and Body
type Parameter struct {
	name, description, dataType string
	kind                        int
	required                    bool
	allowableValues             map[string]string
	allowMultiple               bool
}

func (self *Parameter) bePath() *Parameter {
	self.kind = PATH_PARAMETER
	return self
}
func (self *Parameter) beQuery() *Parameter {
	self.kind = QUERY_PARAMETER
	return self
}
func (self *Parameter) beBody() *Parameter {
	self.kind = BODY_PARAMETER
	return self
}

// Required sets the required field and return the receiver
func (self *Parameter) Required(required bool) *Parameter {
	self.required = required
	return self
}

// AllowMultiple sets the allowMultiple field and return the receiver
func (self *Parameter) AllowMultiple(multiple bool) *Parameter {
	self.allowMultiple = multiple
	return self
}

// AllowableValues sets the allowableValues field and return the receiver
func (self *Parameter) AllowableValues(values map[string]string) *Parameter {
	self.allowableValues = values
	return self
}

// DataType sets the dataType field and return the receiver
func (self *Parameter) DataType(typeName string) *Parameter {
	self.dataType = typeName
	return self
}
