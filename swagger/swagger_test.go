package swagger

import (
	"encoding/json"
	"os"
	"testing"
)

func TestApi(t *testing.T) {
	value := Api{Path: "/", Description: "Some Path", Operations: []Operation{}, Models: map[string]Model{}}
	output, _ := json.MarshalIndent(value, " ", " ")
	print(string(output))
}

type sample struct {
	id       string
	items    []item
	rootItem item `json:"root"`
}

type item struct {
	itemName string `json:"name"`
}

func TestModelToJsonSchema(t *testing.T) {
	api := new(Api)
	api.Models = map[string]Model{}
	op := new(Operation)
	addModelFromSample(api, op, true, sample{items: []item{}})
	output, _ := json.MarshalIndent(api, " ", " ")
	os.Stdout.Write(output)
}
