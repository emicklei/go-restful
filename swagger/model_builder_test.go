package swagger

import (
	"encoding/json"
	"os"
	"testing"
)

func TestJsonTags(t *testing.T) {
	type X struct {
		A string
		B string `json:"C"`
		D string `json:"-"`
		E int    `json:",string"`
		F int    `json:",omitempty"`
		G int    `json:"H,omitempty"`
		I int    `json:","`
	}

	expected := `{
  "id": "swagger.X",
  "required": [
   "A",
   "C",
   "E",
   "I"
  ],
  "properties": {
   "A": {
    "type": "string",
    "description": ""
   },
   "C": {
    "type": "string",
    "description": ""
   },
   "E": {
    "type": "string",
    "description": "(int as string)"
   },
   "F": {
    "type": "int",
    "description": ""
   },
   "H": {
    "type": "int",
    "description": ""
   },
   "I": {
    "type": "int",
    "description": ""
   }
  }
 }`

	sws := newSwaggerService(Config{})
	decl := ApiDeclaration{Models: map[string]Model{}}
	sws.addModelFromSampleTo(&Operation{}, true, X{}, &decl)

	output, _ := json.MarshalIndent(decl.Models["swagger.X"], " ", " ")
	if string(output) != expected {
		t.Error("output != expected")
		os.Stdout.WriteString(expected)
	}
	os.Stdout.Write(output)
}
