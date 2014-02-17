package swagger

import (
	"testing"
	"time"
)

type YesNo bool

func (y YesNo) MarshalJSON() ([]byte, error) {
	if y {
		return []byte("yes"), nil
	}
	return []byte("no"), nil
}

// clear && go test -v -test.run TestCustomMarshaller_Issue96 ...swagger
func TestCustomMarshaller_Issue96(t *testing.T) {
	type Vote struct {
		What YesNo
	}
	testJsonFromStruct(t, Vote{}, `{
  "swagger.Vote": {
   "id": "swagger.Vote",
   "required": [
    "What"
   ],
   "properties": {
    "What": {
     "type": "string",
     "description": "",
     "format": ""
    }
   }
  }
 }`)
}

// clear && go test -v -test.run TestPrimitiveTypes ...swagger
func TestPrimitiveTypes(t *testing.T) {
	type Prims struct {
		f float64
		t time.Time
	}
	testJsonFromStruct(t, Prims{}, `{
  "swagger.Prims": {
   "id": "swagger.Prims",
   "required": [
    "f",
    "t"
   ],
   "properties": {
    "f": {
     "type": "number",
     "description": "",
     "format": "double"
    },
    "t": {
     "type": "string",
     "description": "",
     "format": "date-time"
    }
   }
  }
 }`)
}

// clear && go test -v -test.run TestS1 ...swagger
func TestS1(t *testing.T) {
	type S1 struct {
		Id string
	}
	testJsonFromStruct(t, S1{}, `{
  "swagger.S1": {
   "id": "swagger.S1",
   "required": [
    "Id"
   ],
   "properties": {
    "Id": {
     "type": "string",
     "description": "",
     "format": ""
    }
   }
  }
 }`)
}

// clear && go test -v -test.run TestS2 ...swagger
func TestS2(t *testing.T) {
	type S2 struct {
		Ids []string
	}
	testJsonFromStruct(t, S2{}, `{
  "swagger.S2": {
   "id": "swagger.S2",
   "required": [
    "Ids"
   ],
   "properties": {
    "Ids": {
     "type": "array",
     "description": "",
     "items": {
      "$ref": "string"
     },
     "format": ""
    }
   }
  }
 }`)
}

// clear && go test -v -test.run TestS3 ...swagger
func TestS3(t *testing.T) {
	type NestedS3 struct {
		Id string
	}
	type S3 struct {
		Nested NestedS3
	}
	testJsonFromStruct(t, S3{}, `{
  "swagger.NestedS3": {
   "id": "swagger.NestedS3",
   "required": [
    "Id"
   ],
   "properties": {
    "Id": {
     "type": "string",
     "description": "",
     "format": ""
    }
   }
  },
  "swagger.S3": {
   "id": "swagger.S3",
   "required": [
    "Nested"
   ],
   "properties": {
    "Nested": {
     "type": "swagger.NestedS3",
     "description": "",
     "format": ""
    }
   }
  }
 }`)
}

type sample struct {
	id       string `swagger:"required"` // TODO
	items    []item
	rootItem item `json:"root"`
}

type item struct {
	itemName string `json:"name"`
}

// clear && go test -v -test.run TestSampleToModelAsJson ...swagger
func TestSampleToModelAsJson(t *testing.T) {
	testJsonFromStruct(t, sample{items: []item{}}, `{
  "swagger.item": {
   "id": "swagger.item",
   "required": [
    "name"
   ],
   "properties": {
    "name": {
     "type": "string",
     "description": "",
     "format": ""
    }
   }
  },
  "swagger.sample": {
   "id": "swagger.sample",
   "required": [
    "id",
    "items",
    "root"
   ],
   "properties": {
    "id": {
     "type": "string",
     "description": "",
     "format": ""
    },
    "items": {
     "type": "array",
     "description": "",
     "items": {
      "$ref": "swagger.item"
     },
     "format": ""
    },
    "root": {
     "type": "swagger.item",
     "description": "",
     "format": ""
    }
   }
  }
 }`)
}

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
     "description": "",
     "format": ""
    },
    "C": {
     "type": "string",
     "description": "(int as string)",
     "format": ""
    },
    "D": {
     "type": "integer",
     "description": "",
     "format": "int32"
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
     "type": "integer",
     "description": "",
     "format": "int32"
    },
    "C": {
     "type": "integer",
     "description": "",
     "format": "int32"
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
     "description": "",
     "format": ""
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
     "description": "",
     "format": ""
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
     "type": "integer",
     "description": "",
     "format": "int32"
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
     "description": "",
     "format": ""
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
     "type": "integer",
     "description": "",
     "format": "int32"
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
     },
     "format": ""
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
     "type": "integer",
     "description": "",
     "format": "int32"
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
     },
     "format": ""
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
     "type": "integer",
     "description": "",
     "format": "int32"
    }
   }
  }
 }`

	testJsonFromStruct(t, X{}, expected)
}

// go test -v -test.run TestEmbeddedStruct_Issue98 ...swagger
func TestEmbeddedStruct_Issue98(t *testing.T) {
	type Y struct {
		A int
	}
	type X struct {
		Y
	}
	testJsonFromStruct(t, X{}, `{
  "swagger.X": {
   "id": "swagger.X",
   "required": [
    "A"
   ],
   "properties": {
    "A": {
     "type": "integer",
     "description": "",
     "format": "int32"
    }
   }
  }
 }`)
}

type Dataset struct {
	Names []string
}

// clear && go test -v -test.run TestIssue85 ...swagger
func TestIssue85(t *testing.T) {
	anon := struct{ Datasets []Dataset }{}
	testJsonFromStruct(t, anon, `{
  "struct { Datasets ||swagger.Dataset }": {
   "id": "struct { Datasets ||swagger.Dataset }",
   "required": [
    "Datasets"
   ],
   "properties": {
    "Datasets": {
     "type": "array",
     "description": "",
     "items": {
      "$ref": "swagger.Dataset"
     },
     "format": ""
    }
   }
  },
  "swagger.Dataset": {
   "id": "swagger.Dataset",
   "required": [
    "Names"
   ],
   "properties": {
    "Names": {
     "type": "array",
     "description": "",
     "items": {
      "$ref": "string"
     },
     "format": ""
    }
   }
  }
 }`)
}

type File struct {
	History     []File
	HistoryPtrs []*File
}

// go test -v -test.run TestRecursiveStructure ...swagger
func TestRecursiveStructure(t *testing.T) {
	testJsonFromStruct(t, File{}, `{
  "swagger.File": {
   "id": "swagger.File",
   "required": [
    "History",
    "HistoryPtrs"
   ],
   "properties": {
    "History": {
     "type": "array",
     "description": "",
     "items": {
      "$ref": "swagger.File"
     },
     "format": ""
    },
    "HistoryPtrs": {
     "type": "array",
     "description": "",
     "items": {
      "$ref": "swagger.File.HistoryPtrs"
     },
     "format": ""
    }
   }
  },
  "swagger.File.HistoryPtrs": {
   "id": "swagger.File.HistoryPtrs",
   "properties": {}
  }
 }`)
}

type A1 struct {
	B struct {
		Id int
	}
}

// go test -v -test.run TestEmbeddedStructA1 ...swagger
func TestEmbeddedStructA1(t *testing.T) {
	testJsonFromStruct(t, A1{}, `{
  "swagger.A1": {
   "id": "swagger.A1",
   "required": [
    "B"
   ],
   "properties": {
    "B": {
     "type": "swagger.A1.B",
     "description": "",
     "format": ""
    }
   }
  },
  "swagger.A1.B": {
   "id": "swagger.A1.B",
   "required": [
    "Id"
   ],
   "properties": {
    "Id": {
     "type": "integer",
     "description": "",
     "format": "int32"
    }
   }
  }
 }`)
}

type A2 struct {
	C
}
type C struct {
	Id int `json:"B"`
}

// go test -v -test.run TestEmbeddedStructA2 ...swagger
func TestEmbeddedStructA2(t *testing.T) {
	testJsonFromStruct(t, A2{}, `{
  "swagger.A2": {
   "id": "swagger.A2",
   "required": [
    "B"
   ],
   "properties": {
    "B": {
     "type": "integer",
     "description": "",
     "format": "int32"
    }
   }
  }
 }`)
}

type A3 struct {
	B D
}

type D struct {
	Id int
}

// clear && go test -v -test.run TestStructA3 ...swagger
func TestStructA3(t *testing.T) {
	testJsonFromStruct(t, A3{}, `{
  "swagger.A3": {
   "id": "swagger.A3",
   "required": [
    "B"
   ],
   "properties": {
    "B": {
     "type": "swagger.D",
     "description": "",
     "format": ""
    }
   }
  },
  "swagger.D": {
   "id": "swagger.D",
   "required": [
    "Id"
   ],
   "properties": {
    "Id": {
     "type": "integer",
     "description": "",
     "format": "int32"
    }
   }
  }
 }`)
}
