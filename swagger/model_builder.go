package swagger

import (
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
			sm.Properties[jsonName] = prop
		}
	}

	// update model builder with completed model
	b.Models[modelName] = sm
}

func (b modelBuilder) buildProperty(field reflect.StructField, sm *Model, modelName string) (string, ModelProperty) {
	jsonName := field.Name
	fieldType := field.Type
	fieldKind := fieldType.Kind()
	prop := ModelProperty{}

	if fieldKind == reflect.Struct {
		// embedded struct
		sub := modelBuilder{map[string]Model{}}
		sub.addModel(fieldType, "")
		subKey := sub.keyFrom(fieldType)
		// merge properties from sub
		subModel := sub.Models[subKey]
		for k, v := range subModel.Properties {
			sm.Properties[k] = v
			sm.Required = append(sm.Required, k)
		}
		// empty name signals skip property
		return "", prop
	}

	required := true
	// see if a tag overrides this
	if jsonTag := field.Tag.Get("json"); jsonTag != "" {
		s := strings.Split(jsonTag, ",")
		if s[0] == "-" {
			return "", prop
		} else if s[0] != "" {
			jsonName = s[0]
		}
		if len(s) > 1 {
			switch s[1] {
			case "string":
				prop.Description = "(" + fieldType.String() + " as string)"
				fieldType = reflect.TypeOf("")
			case "omitempty":
				required = false
			}
		}
	}
	prop.Type = b.jsonSchemaType(fieldType.String()) // may include pkg path
	//if format := b.jsonSchemaFormat(fieldType.String()); len(format) > 0 {
	//	prop.Format = format
	//}

	if required {
		sm.Required = append(sm.Required, jsonName)
	}
	if fieldKind == reflect.Slice || fieldKind == reflect.Array {
		// list like
		prop.Type = "array"
		elemName := b.getElementTypeName(modelName, jsonName, fieldType.Elem())
		prop.Items = map[string]string{"$ref": elemName}
		// add|overwrite model for element type
		b.addModel(fieldType.Elem(), elemName)
	} else if fieldKind == reflect.Ptr {
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
	} else if fieldType.Name() == "" { // override type of anonymous structs
		prop.Type = modelName + "." + jsonName
		b.addModel(fieldType, prop.Type)
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
	return strings.Contains("int int32 int64 float32 float64 bool string byte", modelName)
}

func (b modelBuilder) jsonSchemaType(modelName string) string {
	schemaMap := map[string]string{
		"int":       "integer",
		"float64":   "number",
		"bool":      "boolean",
		"time.Time": "date",
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
		"int32":   "int32",
		"int64":   "int64",
		"float64": "double",
	}
	mapped, ok := schemaMap[modelName]
	if ok {
		return mapped
	} else {
		return "" // no format
	}
}
