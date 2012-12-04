package swagger

import (
	"encoding/json"
	"testing"
)

func TestApi(t *testing.T) {
	value := Api{Path: "/", Description: "Some Path", Operations: []Operation{}, Models: []Model{}}
	output, _ := json.MarshalIndent(value, " ", " ")
	print(string(output))
}
