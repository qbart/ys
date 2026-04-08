package ys

// Schema defines the interface for all schema types.
type Schema interface {
	validate(node *yamlNode, path string) []SchemaError
}

// Result is returned from Validate and contains the validation outcome.
type Result struct {
	OK     bool
	Errors []SchemaError
}

// SchemaError represents a single validation error with location info.
type SchemaError struct {
	Path    string
	Line    int
	Message string
}

// Field represents a named field in an object schema.
type Field struct {
	Name     string
	Schema   Schema
	Required bool
}

// Required creates a required field definition.
func Required(name string, schema Schema) Field {
	return Field{Name: name, Schema: schema, Required: true}
}

// Optional creates an optional field definition.
func Optional(name string, schema Schema) Field {
	return Field{Name: name, Schema: schema, Required: false}
}
