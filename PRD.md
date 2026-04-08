# PRD: `ys` — YAML Schema Validator for Go

## Overview

`ys` is a Go library that validates YAML documents against schemas defined in Go code. Schemas use a developer-friendly DSL that reads naturally and supports nested objects, arrays, required/optional fields, and type checking.

## Goals

- **Developer-friendly DSL**: Schema definitions should be easy to read and write in Go
- **Rich error reporting**: Return structured errors with line numbers from the YAML source
- **Type safety**: Validate strings, integers, floats, booleans, arrays, objects
- **Nested structures**: Support deeply nested objects and arrays with their own schemas
- **Required/Optional**: Fields can be marked as required or optional at any nesting level

## API Design

### Schema Definition

```go
schema := ys.Object(
    ys.Required("name", ys.String()),
    ys.Required("age", ys.Int()),
    ys.Optional("email", ys.String()),
    ys.Required("address", ys.Object(
        ys.Required("street", ys.String()),
        ys.Required("city", ys.String()),
        ys.Optional("zip", ys.String()),
    )),
    ys.Optional("tags", ys.Array(ys.String())),
)
```

### Validation

```go
result, err := ys.Validate(yamlBytes, schema)
if err != nil {
    // parse error or internal error
}
if !result.OK {
    for _, e := range result.Errors {
        fmt.Printf("line %d: %s\n", e.Line, e.Message)
    }
}
```

### Core Types

```go
type Result struct {
    OK     bool
    Errors []SchemaError
}

type SchemaError struct {
    Path    string // e.g. "address.street"
    Line    int    // line number in YAML source
    Message string // human-readable error message
}
```

### Schema Types

| Function | Description |
|----------|-------------|
| `ys.String()` | Validates string values |
| `ys.Int()` | Validates integer values |
| `ys.Float()` | Validates float values |
| `ys.Bool()` | Validates boolean values |
| `ys.Object(fields...)` | Validates an object with given fields |
| `ys.Array(itemSchema)` | Validates an array where each item matches schema |
| `ys.Any()` | Accepts any value |

### Field Definitions

| Function | Description |
|----------|-------------|
| `ys.Required(name, schema)` | Field must be present and match schema |
| `ys.Optional(name, schema)` | Field may be absent; if present must match schema |

## Non-Goals (v1)

- Custom validation functions / constraints (min, max, regex, etc.)
- Schema composition / references
- YAML generation from schema
- Loading schemas from YAML/JSON files

## Implementation Plan

1. Project scaffolding (go module, basic types)
2. String type validation
3. Int type validation
4. Float type validation
5. Bool type validation
6. Any type validation
7. Required/Optional field definitions
8. Object validation (flat)
9. Nested object validation
10. Array validation
11. Line number tracking in errors
12. Edge cases (empty docs, null values, type mismatches)
13. Complex nested schemas (arrays of objects, etc.)
14. Final integration tests with realistic YAML samples
