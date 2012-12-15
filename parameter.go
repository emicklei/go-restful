package restful

const (
	PATH_PARAMETER = iota
	QUERY_PARAMETER
	BODY_PARAMETER
)

// Parameter is for documententing the parameter used in a Http Request
// Parameter kinds are Path,Query and Body
type Parameter struct {
	Name, Description, DataType string
	kind                        int
	Required                    bool
	AllowableValues             map[string]string
	AllowMultiple               bool
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

// Optional sets the required field
func (self *Parameter) Optional(optional bool) *Parameter {
	self.Required = !optional
	return self
}

// MultipleAllowed sets the AllowMultiple field
func (self *Parameter) MultipleAllowed(multiple bool) *Parameter {
	self.AllowMultiple = multiple
	return self
}

// ValuesAllowed sets the AllowableValues field
func (self *Parameter) ValuesAllowed(values map[string]string) *Parameter {
	self.AllowableValues = values
	return self
}
