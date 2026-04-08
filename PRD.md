# PRD: yamlschema - YAML Schema Validator for Go

## Overview

`yamlschema` is a Go library (package `ys`) that allows developers to define YAML schemas in Go code and validate YAML documents against them. The library prioritizes ease of use and code readability.

**Module path:** `github.com/qbart/yamlschema`
**Package name:** `ys`

## Core API Design

### Schema Definition

```go
schema := ys.Schema{
    "name":    ys.String().Required(),
    "age":     ys.Int().Optional(),
    "email":   ys.String().Required(),
    "tags":    ys.Array(ys.String()),
    "address": ys.Object(ys.Schema{
        "street": ys.String().Required(),
        "city":   ys.String().Required(),
        "zip":    ys.String().Optional(),
    }),
}
```

### Validation

```go
result, err := schema.Validate(yamlBytes)
// err is for parse/system errors (invalid YAML, etc.)
// result contains validation outcome
```

### Result

```go
type Result struct {
    OK     bool
    Errors []SchemaError
}

type SchemaError struct {
    Line    int
    Path    string   // e.g. "address.street"
    Message string
}
```

## Type System

| Builder     | Description                |
|-------------|----------------------------|
| `ys.String()` | String value             |
| `ys.Int()`    | Integer value            |
| `ys.Float()`  | Float value              |
| `ys.Bool()`   | Boolean value            |
| `ys.Array(T)` | Array of type T          |
| `ys.Object(S)`| Nested object with schema|
| `ys.Any()`    | Any value (skip type check)|

Each type supports:
- `.Required()` — field must be present (default)
- `.Optional()` — field may be absent

## YAML Seed Test Cases

### Seed 1: Simple flat object (string fields only)
```yaml
name: "Alice"
email: "alice@example.com"
```

### Seed 2: Mixed types (string, int, bool)
```yaml
name: "Bob"
age: 30
active: true
```

### Seed 3: Optional fields (present and absent)
```yaml
name: "Charlie"
# bio is optional and absent
```

### Seed 4: Required field missing
```yaml
# name is required but missing
age: 25
```

### Seed 5: Wrong type
```yaml
name: 123        # expected string
age: "not a num" # expected int
```

### Seed 6: Nested object
```yaml
name: "Diana"
address:
  street: "123 Main St"
  city: "Springfield"
  zip: "62704"
```

### Seed 7: Nested object with missing required field
```yaml
name: "Eve"
address:
  street: "456 Elm St"
  # city is required but missing
```

### Seed 8: Array of scalars
```yaml
name: "Frank"
tags:
  - "go"
  - "yaml"
  - "dev"
```

### Seed 9: Array with wrong element type
```yaml
name: "Grace"
tags:
  - 1
  - 2
  - 3
```

### Seed 10: Array of objects
```yaml
name: "Hank"
friends:
  - name: "Ivy"
    age: 28
  - name: "Jack"
    age: 32
```

### Seed 11: Deeply nested objects
```yaml
company:
  name: "Acme"
  hq:
    address:
      street: "789 Oak Ave"
      city: "Metropolis"
```

### Seed 12: Empty document
```yaml
```

### Seed 13: Null values
```yaml
name: null
age: null
```

### Seed 14: Extra fields not in schema
```yaml
name: "Kelly"
unknown_field: "should this error?"
```

### Seed 15: Complex mixed (arrays of objects with nested arrays)
```yaml
name: "Leo"
projects:
  - title: "Alpha"
    tags:
      - "web"
      - "api"
    members:
      - name: "Mia"
      - name: "Noah"
  - title: "Beta"
    tags:
      - "cli"
    members:
      - name: "Olivia"
```

### Seed 16: Float values
```yaml
name: "Pat"
score: 9.5
ratio: 0.75
```

### Seed 17: Boolean validation
```yaml
enabled: true
verbose: "yes"  # not a bool
```

### Seed 18: Any type
```yaml
metadata: "could be anything"
data: 42
flag: true
```

## Edge Cases to Cover

- Invalid YAML syntax (returns error, not result)
- Empty byte slice
- YAML with only comments
- Unicode field names and values
- Very large nested structures
- Array at root level (not an object) — should error
- Numeric strings vs actual numbers
- Trailing whitespace / multiline strings

## Implementation Plan (TDD Increments)

1. **Project scaffolding** — go.mod, empty package, first test file
2. **Schema type + String().Required()** — validate a single required string field
3. **Result/SchemaError structs** — return structured errors with line numbers
4. **Multiple string fields** — validate documents with several fields
5. **Optional fields** — `.Optional()` modifier, absent fields pass
6. **Required field missing** — proper error with path and line
7. **Int type** — `ys.Int()`, type mismatch errors
8. **Bool type** — `ys.Bool()`
9. **Float type** — `ys.Float()`
10. **Nested objects** — `ys.Object(schema)`
11. **Arrays of scalars** — `ys.Array(ys.String())`
12. **Arrays of objects** — `ys.Array(ys.Object(schema))`
13. **Deeply nested structures** — multi-level nesting
14. **Any type** — `ys.Any()`
15. **Edge cases** — invalid YAML, empty input, null values, extra fields, etc.
16. **Line number tracking** — accurate line reporting for all error types
17. **Complex integration tests** — seeds 15, full coverage
