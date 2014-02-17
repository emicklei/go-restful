package swagger

import (
	"encoding/json"
	"reflect"
	"strings"
)

type modelBuilder struct {
	Models map[string]Model
}

func (b modelBuilder) addModel(st reflect.Type, nameOverride string) {
	modelName := b.keyFrom(st)
	if nameOverride != "" {
		modelName = nameOverride
	}
	// no models needed for primitive types
	if b.isPrimitiveType(modelName) {
		return
	}
	// see if we already have visited this model
	if _, ok := b.Models[modelName]; ok {
		return
	}
	sm := Model{modelName, []string{}, map[string]ModelProperty{}}

	// reference the model before further initializing (enables recursive structs)
	b.Models[modelName] = sm

	// check for structure or primitive type
	if st.Kind() != reflect.Struct {
		return
	}
	for i := 0; i < st.NumField(); i++ {
		field := st.Field(i)
		jsonName, prop := b.buildProperty(field, &sm, modelName)
		// add if not ommitted
		if len(jsonName) != 0 {
			// update Required
			if b.isPropertyRequired(field) {
				sm.Required = append(sm.Required, jsonName)
			}
			sm.Properties[jsonName] = prop
		}
	}

	// update model builder with completed model
	b.Models[modelName] = sm
}

func (b modelBuilder) isPropertyRequired(field reflect.StructField) bool {
	required := true
	if jsonTag := field.Tag.Get("json"); jsonTag != "" {
		s := strings.Split(jsonTag, ",")
		if len(s) > 1 && s[1] == "omitempty" {
			return false
		}
	}
	return required
}

func (b modelBuilder) buildProperty(field reflect.StructField, model *Model, modelName string) (jsonName string, prop ModelProperty) {
	jsonName = b.jsonNameOfField(field)
	if len(jsonName) == 0 {
		// empty name signals skip property
		return "", prop
	}
	fieldType := field.Type
	fieldKind := fieldType.Kind()

	if jsonTag := field.Tag.Get("json"); jsonTag != "" {
		s := strings.Split(jsonTag, ",")
		if len(s) > 1 && s[1] == "string" {
			prop.Description = "(" + fieldType.String() + " as string)"
			fieldType = reflect.TypeOf("")
		}
	}

	prop.Type = b.jsonSchemaType(fieldType.String()) // may include pkg path
	if b.isPrimitiveType(fieldType.String()) {
		prop.Format = b.jsonSchemaFormat(fieldType.String())
		return jsonName, prop
	}

	marshalerType := reflect.TypeOf((*json.Marshaler)(nil)).Elem()
	if fieldType.Implements(marshalerType) {
		prop.Type = "string"
		return jsonName, prop
	}

	if fieldKind == reflect.Struct {
		return b.buildStructTypeProperty(field, jsonName, model)
	}

	if fieldKind == reflect.Slice || fieldKind == reflect.Array {
		return b.buildArrayTypeProperty(field, jsonName, modelName)
	}

	if fieldKind == reflect.Ptr {
		return b.buildPointerTypeProperty(field, jsonName, modelName)
	}

	if fieldType.Name() == "" { // override type of anonymous structs
		nestedTypeName := modelName + "." + jsonName
		prop.Type = nestedTypeName
		b.addModel(fieldType, nestedTypeName)
	}
	return jsonName, prop
}

func (b modelBuilder) buildStructTypeProperty(field reflect.StructField, jsonName string, model *Model) (nameJson string, prop ModelProperty) {
	fieldType := field.Type
	// check for anonymous
	if len(fieldType.Name()) == 0 {
		// anonymous
		anonType := model.Id + "." + jsonName
		b.addModel(fieldType, anonType)
		prop.Type = anonType
		return jsonName, prop
	}
	if field.Name == fieldType.Name() {
		// embedded struct
		sub := modelBuilder{map[string]Model{}}
		sub.addModel(fieldType, "")
		subKey := sub.keyFrom(fieldType)
		// merge properties from sub
		subModel := sub.Models[subKey]
		for k, v := range subModel.Properties {
			model.Properties[k] = v
			model.Required = append(model.Required, k)
		}
		// empty name signals skip property
		return "", prop
	}
	// simple struct
	b.addModel(fieldType, "")
	prop.Type = fieldType.String()
	return jsonName, prop
}

func (b modelBuilder) buildArrayTypeProperty(field reflect.StructField, jsonName, modelName string) (nameJson string, prop ModelProperty) {
	fieldType := field.Type
	prop.Type = "array"
	elemName := b.getElementTypeName(modelName, jsonName, fieldType.Elem())
	prop.Items = map[string]string{"$ref": elemName}
	// add|overwrite model for element type
	b.addModel(fieldType.Elem(), elemName)
	return jsonName, prop
}

func (b modelBuilder) buildPointerTypeProperty(field reflect.StructField, jsonName, modelName string) (nameJson string, prop ModelProperty) {
	fieldType := field.Type

	// override type of pointer to list-likes
	if fieldType.Elem().Kind() == reflect.Slice || fieldType.Elem().Kind() == reflect.Array {
		prop.Type = "array"
		elemName := b.getElementTypeName(modelName, jsonName, fieldType.Elem().Elem())
		prop.Items = map[string]string{"$ref": elemName}
		// add|overwrite model for element type
		b.addModel(fieldType.Elem().Elem(), elemName)
	} else {
		// non-array, pointer type
		prop.Type = fieldType.String()[1:] // no star, include pkg path
		elemName := ""
		if fieldType.Elem().Name() == "" {
			elemName = modelName + "." + jsonName
			prop.Type = elemName
		}
		b.addModel(fieldType.Elem(), elemName)
	}
	return jsonName, prop
}

func (b modelBuilder) getElementTypeName(modelName, jsonName string, t reflect.Type) string {
	if t.Name() == "" {
		return modelName + "." + jsonName
	}
	return b.keyFrom(t)
}

func (b modelBuilder) keyFrom(st reflect.Type) string {
	key := st.String()
	if len(st.Name()) == 0 { // unnamed type
		// Swagger UI has special meaning for [
		key = strings.Replace(key, "[]", "||", -1)
	}
	return key
}

func (b modelBuilder) isPrimitiveType(modelName string) bool {
	return strings.Contains("int int32 int64 float32 float64 bool string byte time.Time", modelName)
}

// jsonNameOfField returns the name of the field as it should appear in JSON format
// An empty string indicates that this field is not part of the JSON representation
func (b modelBuilder) jsonNameOfField(field reflect.StructField) string {
	if jsonTag := field.Tag.Get("json"); jsonTag != "" {
		s := strings.Split(jsonTag, ",")
		if s[0] == "-" {
			// empty name signals skip property
			return ""
		} else if s[0] != "" {
			return s[0]
		}
	}
	return field.Name
}

func (b modelBuilder) jsonSchemaType(modelName string) string {
	schemaMap := map[string]string{
		"int":       "integer",
		"int32":     "integer",
		"int64":     "integer",
		"byte":      "string",
		"float64":   "number",
		"float32":   "number",
		"bool":      "boolean",
		"time.Time": "string",
	}
	mapped, ok := schemaMap[modelName]
	if ok {
		return mapped
	} else {
		return modelName // use as is (custom or struct)
	}
}

func (b modelBuilder) jsonSchemaFormat(modelName string) string {
	schemaMap := map[string]string{
		"int":       "int32",
		"int32":     "int32",
		"int64":     "int64",
		"byte":      "byte",
		"float64":   "double",
		"float32":   "float",
		"time.Time": "date-time",
	}
	mapped, ok := schemaMap[modelName]
	if ok {
		return mapped
	} else {
		return "" // no format
	}
}
