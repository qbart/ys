package ys

import (
	"gopkg.in/yaml.v3"
)

// yamlNode wraps yaml.Node for internal use.
type yamlNode = yaml.Node

// Validate parses YAML data and validates it against the given schema.
func Validate(data []byte, schema Schema) (Result, error) {
	var doc yaml.Node
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return Result{}, err
	}

	// yaml.Unmarshal wraps content in a document node
	if doc.Kind != yaml.DocumentNode || len(doc.Content) == 0 {
		return Result{
			OK:     false,
			Errors: []SchemaError{{Path: "", Line: 0, Message: "empty document"}},
		}, nil
	}

	root := doc.Content[0]
	errors := schema.validate(root, "")

	return Result{
		OK:     len(errors) == 0,
		Errors: errors,
	}, nil
}
