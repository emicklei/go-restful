package restful

const (
	PATH_PARAMETER = iota
	QUERY_PARAMETER
	BODY_PARAMETER
)

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
func (self *Parameter) Optional(optional bool) *Parameter {
	self.Required = !optional
	return self
}
func (self *Parameter) MultipleAllowed(multiple bool) *Parameter {
	self.AllowMultiple = multiple
	return self
}
func (self *Parameter) ValuesAllowed(values map[string]string) *Parameter {
	self.AllowableValues = values
	return self
}
