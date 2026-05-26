package prmaven

import (
	"encoding/json"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestReportSchemaTracksJSONContractFields(t *testing.T) {
	data, err := os.ReadFile("../../schema/prmaven-report.schema.json")
	if err != nil {
		t.Fatal(err)
	}

	var schema map[string]any
	if err := json.Unmarshal(data, &schema); err != nil {
		t.Fatal(err)
	}

	assertSchemaProperties(t, schema, jsonFields(reflect.TypeOf(Report{})), "properties")
	assertSchemaRequired(t, schema, requiredJSONFields(reflect.TypeOf(Report{})), "required")
	defs := schemaObject(t, schema, "$defs")
	assertSchemaProperties(t, defs, jsonFields(reflect.TypeOf(Summary{})), "summary.properties")
	assertSchemaRequired(t, defs, requiredJSONFields(reflect.TypeOf(Summary{})), "summary.required")
	assertSchemaProperties(t, defs, jsonFields(reflect.TypeOf(Module{})), "module.properties")
	assertSchemaRequired(t, defs, requiredJSONFields(reflect.TypeOf(Module{})), "module.required")
	assertSchemaProperties(t, defs, jsonFields(reflect.TypeOf(Finding{})), "finding.properties")
	assertSchemaRequired(t, defs, requiredJSONFields(reflect.TypeOf(Finding{})), "finding.required")
}

func assertSchemaProperties(t *testing.T, schema map[string]any, expected []string, path string) {
	t.Helper()

	properties := schemaMapAtPath(t, schema, path)
	for _, field := range expected {
		if _, ok := properties[field]; !ok {
			t.Fatalf("%s missing JSON field %q", path, field)
		}
	}
}

func assertSchemaRequired(t *testing.T, schema map[string]any, expected []string, path string) {
	t.Helper()

	requiredValue := schemaObjectAtPath(t, schema, path)
	required, ok := requiredValue.([]any)
	if !ok {
		t.Fatalf("%s is %T, want []any", path, requiredValue)
	}

	got := map[string]bool{}
	for _, value := range required {
		name, ok := value.(string)
		if !ok {
			t.Fatalf("%s contains %T, want string", path, value)
		}
		got[name] = true
	}
	for _, field := range expected {
		if !got[field] {
			t.Fatalf("%s missing required JSON field %q", path, field)
		}
	}
}

func schemaObjectAtPath(t *testing.T, schema map[string]any, path string) any {
	t.Helper()

	parts := strings.Split(path, ".")
	var current any = schema
	for _, part := range parts {
		object, ok := current.(map[string]any)
		if !ok {
			t.Fatalf("%s parent is %T, want object", path, current)
		}
		current = object[part]
	}
	return current
}

func schemaMapAtPath(t *testing.T, schema map[string]any, path string) map[string]any {
	t.Helper()

	value := schemaObjectAtPath(t, schema, path)
	object, ok := value.(map[string]any)
	if !ok {
		t.Fatalf("%s is %T, want object", path, value)
	}
	return object
}

func schemaObject(t *testing.T, schema map[string]any, key string) map[string]any {
	t.Helper()

	value, ok := schema[key]
	if !ok {
		t.Fatalf("schema missing key %q", key)
	}
	object, ok := value.(map[string]any)
	if !ok {
		t.Fatalf("schema key %q is %T, want object", key, value)
	}
	return object
}

func jsonFields(typ reflect.Type) []string {
	fields := make([]string, 0, typ.NumField())
	for i := range typ.NumField() {
		if name := jsonFieldName(typ.Field(i)); name != "" {
			fields = append(fields, name)
		}
	}
	return fields
}

func requiredJSONFields(typ reflect.Type) []string {
	fields := make([]string, 0, typ.NumField())
	for i := range typ.NumField() {
		field := typ.Field(i)
		tag := field.Tag.Get("json")
		if strings.Contains(tag, "omitempty") {
			continue
		}
		if name := jsonFieldName(field); name != "" {
			fields = append(fields, name)
		}
	}
	return fields
}

func jsonFieldName(field reflect.StructField) string {
	tag := field.Tag.Get("json")
	if tag == "" || tag == "-" {
		return ""
	}
	return strings.Split(tag, ",")[0]
}
