// +build "go1.1","!go1.10"

package restful

import "encoding/json"

// Read unmarshalls the value from JSON
func (e entityJSONAccess) Read(req *Request, v interface{}) error {
	decoder := json.NewDecoder(req.Request.Body)
	decoder.UseNumber()
	return decoder.Decode(v)
}
