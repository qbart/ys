package ys

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// stringSchema validates that a value is a string.
type stringSchema struct{}

func String() Schema { return stringSchema{} }

func (s stringSchema) validate(node *yamlNode, path string) []SchemaError {
	if node.Kind != yaml.ScalarNode || node.Tag != "!!str" {
		return []SchemaError{{
			Path:    path,
			Line:    node.Line,
			Message: fmt.Sprintf("expected string, got %s", describeNode(node)),
		}}
	}
	return nil
}

// intSchema validates that a value is an integer.
type intSchema struct{}

func Int() Schema { return intSchema{} }

func (s intSchema) validate(node *yamlNode, path string) []SchemaError {
	if node.Kind != yaml.ScalarNode || node.Tag != "!!int" {
		return []SchemaError{{
			Path:    path,
			Line:    node.Line,
			Message: fmt.Sprintf("expected int, got %s", describeNode(node)),
		}}
	}
	return nil
}

// floatSchema validates that a value is a float.
type floatSchema struct{}

func Float() Schema { return floatSchema{} }

func (s floatSchema) validate(node *yamlNode, path string) []SchemaError {
	if node.Kind != yaml.ScalarNode || node.Tag != "!!float" {
		return []SchemaError{{
			Path:    path,
			Line:    node.Line,
			Message: fmt.Sprintf("expected float, got %s", describeNode(node)),
		}}
	}
	return nil
}

// boolSchema validates that a value is a boolean.
type boolSchema struct{}

func Bool() Schema { return boolSchema{} }

func (s boolSchema) validate(node *yamlNode, path string) []SchemaError {
	if node.Kind != yaml.ScalarNode || node.Tag != "!!bool" {
		return []SchemaError{{
			Path:    path,
			Line:    node.Line,
			Message: fmt.Sprintf("expected bool, got %s", describeNode(node)),
		}}
	}
	return nil
}

// anySchema accepts any value.
type anySchema struct{}

func Any() Schema { return anySchema{} }

func (s anySchema) validate(node *yamlNode, path string) []SchemaError {
	return nil
}

// objectSchema validates a YAML mapping with defined fields.
type objectSchema struct {
	fields []Field
}

func Object(fields ...Field) Schema {
	return objectSchema{fields: fields}
}

func (s objectSchema) validate(node *yamlNode, path string) []SchemaError {
	if node.Kind != yaml.MappingNode {
		return []SchemaError{{
			Path:    path,
			Line:    node.Line,
			Message: fmt.Sprintf("expected object, got %s", describeNode(node)),
		}}
	}

	// Build a map of key -> value node, and key -> key node (for line numbers)
	keys := make(map[string]*yamlNode)
	vals := make(map[string]*yamlNode)
	for i := 0; i < len(node.Content)-1; i += 2 {
		keyNode := node.Content[i]
		valNode := node.Content[i+1]
		keys[keyNode.Value] = keyNode
		vals[keyNode.Value] = valNode
	}

	var errs []SchemaError

	for _, field := range s.fields {
		fieldPath := joinPath(path, field.Name)
		valNode, exists := vals[field.Name]

		if !exists {
			if field.Required {
				errs = append(errs, SchemaError{
					Path:    fieldPath,
					Line:    node.Line,
					Message: fmt.Sprintf("missing required field %q", field.Name),
				})
			}
			continue
		}

		// Check for null values
		if valNode.Tag == "!!null" {
			if field.Required {
				errs = append(errs, SchemaError{
					Path:    fieldPath,
					Line:    valNode.Line,
					Message: fmt.Sprintf("field %q is null but is required", field.Name),
				})
			}
			continue
		}

		errs = append(errs, field.Schema.validate(valNode, fieldPath)...)
	}

	// Check for extra fields not defined in the schema
	defined := make(map[string]bool, len(s.fields))
	for _, field := range s.fields {
		defined[field.Name] = true
	}
	for i := 0; i < len(node.Content)-1; i += 2 {
		keyNode := node.Content[i]
		if !defined[keyNode.Value] {
			fieldPath := joinPath(path, keyNode.Value)
			errs = append(errs, SchemaError{
				Path:    fieldPath,
				Line:    keyNode.Line,
				Message: fmt.Sprintf("unknown field %q", keyNode.Value),
			})
		}
	}

	return errs
}

// arraySchema validates a YAML sequence where each item matches a schema.
type arraySchema struct {
	itemSchema Schema
}

func Array(items Schema) Schema {
	return arraySchema{itemSchema: items}
}

func (s arraySchema) validate(node *yamlNode, path string) []SchemaError {
	if node.Kind != yaml.SequenceNode {
		return []SchemaError{{
			Path:    path,
			Line:    node.Line,
			Message: fmt.Sprintf("expected array, got %s", describeNode(node)),
		}}
	}

	var errs []SchemaError
	for i, item := range node.Content {
		itemPath := fmt.Sprintf("%s[%d]", path, i)
		errs = append(errs, s.itemSchema.validate(item, itemPath)...)
	}

	return errs
}

// helpers

func joinPath(base, field string) string {
	if base == "" {
		return field
	}
	return base + "." + field
}

func describeNode(node *yamlNode) string {
	switch node.Kind {
	case yaml.ScalarNode:
		switch node.Tag {
		case "!!str":
			return "string"
		case "!!int":
			return "int"
		case "!!float":
			return "float"
		case "!!bool":
			return "bool"
		case "!!null":
			return "null"
		default:
			return "scalar(" + node.Tag + ")"
		}
	case yaml.MappingNode:
		return "object"
	case yaml.SequenceNode:
		return "array"
	default:
		return "unknown"
	}
}
