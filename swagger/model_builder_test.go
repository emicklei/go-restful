package swagger

import (
	"encoding/json"
	"testing"
)

func TestJsonTags(t *testing.T) {
	type X struct {
		A string
		B string `json:"-"`
		C int    `json:",string"`
		D int    `json:","`
	}

	expected := `{
  "swagger.X": {
   "id": "swagger.X",
   "required": [
    "A",
    "C",
    "D"
   ],
   "properties": {
    "A": {
     "type": "string",
     "description": ""
    },
    "C": {
     "type": "string",
     "description": "(int as string)"
    },
    "D": {
     "type": "int",
     "description": ""
    }
   }
  }
 }`

	testJsonFromStruct(t, X{}, expected)
}

func TestJsonTagOmitempty(t *testing.T) {
	type X struct {
		A int `json:",omitempty"`
		B int `json:"C,omitempty"`
	}

	expected := `{
  "swagger.X": {
   "id": "swagger.X",
   "properties": {
    "A": {
     "type": "int",
     "description": ""
    },
    "C": {
     "type": "int",
     "description": ""
    }
   }
  }
 }`

	testJsonFromStruct(t, X{}, expected)
}

func TestJsonTagName(t *testing.T) {
	type X struct {
		A string `json:"B"`
	}

	expected := `{
  "swagger.X": {
   "id": "swagger.X",
   "required": [
    "B"
   ],
   "properties": {
    "B": {
     "type": "string",
     "description": ""
    }
   }
  }
 }`

	testJsonFromStruct(t, X{}, expected)
}

func TestAnonymousStruct(t *testing.T) {
	type X struct {
		A struct {
			B int
		}
	}

	expected := `{
  "swagger.X": {
   "id": "swagger.X",
   "required": [
    "A"
   ],
   "properties": {
    "A": {
     "type": "swagger.X.A",
     "description": ""
    }
   }
  },
  "swagger.X.A": {
   "id": "swagger.X.A",
   "required": [
    "B"
   ],
   "properties": {
    "B": {
     "type": "int",
     "description": ""
    }
   }
  }
 }`

	testJsonFromStruct(t, X{}, expected)
}

func TestAnonymousPtrStruct(t *testing.T) {
	type X struct {
		A *struct {
			B int
		}
	}

	expected := `{
  "swagger.X": {
   "id": "swagger.X",
   "required": [
    "A"
   ],
   "properties": {
    "A": {
     "type": "swagger.X.A",
     "description": ""
    }
   }
  },
  "swagger.X.A": {
   "id": "swagger.X.A",
   "required": [
    "B"
   ],
   "properties": {
    "B": {
     "type": "int",
     "description": ""
    }
   }
  }
 }`

	testJsonFromStruct(t, X{}, expected)
}

func TestAnonymousArrayStruct(t *testing.T) {
	type X struct {
		A []struct {
			B int
		}
	}

	expected := `{
  "swagger.X": {
   "id": "swagger.X",
   "required": [
    "A"
   ],
   "properties": {
    "A": {
     "type": "array",
     "description": "",
     "items": {
      "$ref": "swagger.X.A"
     }
    }
   }
  },
  "swagger.X.A": {
   "id": "swagger.X.A",
   "required": [
    "B"
   ],
   "properties": {
    "B": {
     "type": "int",
     "description": ""
    }
   }
  }
 }`

	testJsonFromStruct(t, X{}, expected)
}

func TestAnonymousPtrArrayStruct(t *testing.T) {
	type X struct {
		A *[]struct {
			B int
		}
	}

	expected := `{
  "swagger.X": {
   "id": "swagger.X",
   "required": [
    "A"
   ],
   "properties": {
    "A": {
     "type": "array",
     "description": "",
     "items": {
      "$ref": "swagger.X.A"
     }
    }
   }
  },
  "swagger.X.A": {
   "id": "swagger.X.A",
   "required": [
    "B"
   ],
   "properties": {
    "B": {
     "type": "int",
     "description": ""
    }
   }
  }
 }`

	testJsonFromStruct(t, X{}, expected)
}

func jsonFromSwaggerService(sample interface{}) string {
	sws := newSwaggerService(Config{})
	decl := ApiDeclaration{Models: map[string]Model{}}
	sws.addModelFromSampleTo(&Operation{}, true, sample, &decl)

	output, _ := json.MarshalIndent(decl.Models, " ", " ")
	return string(output)
}

func testJsonFromStruct(t *testing.T, sample interface{}, expectedJson string) {
	output := jsonFromSwaggerService(sample)
	if output != expectedJson {
		t.Error("output != expected\nexpected:", expectedJson)
	}
	t.Log("output:\n", output)
}
