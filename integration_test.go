package ys

import (
	"os"
	"testing"
)

func loadTestData(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile("testdata/" + name)
	if err != nil {
		t.Fatalf("failed to read testdata/%s: %v", name, err)
	}
	return data
}

// personSchema is a reusable schema for simple person documents
func personSchema() Schema {
	return Object(
		Required("name", String()),
		Required("age", Int()),
		Optional("email", String()),
	)
}

// personWithAddressSchema adds nested address
func personWithAddressSchema() Schema {
	return Object(
		Required("name", String()),
		Required("age", Int()),
		Optional("email", String()),
		Required("address", Object(
			Required("street", String()),
			Required("city", String()),
			Optional("zip", String()),
		)),
	)
}

// companySchema is a complex schema for the company test data
func companySchema() Schema {
	return Object(
		Required("name", String()),
		Required("active", Bool()),
		Required("employees", Array(Object(
			Required("name", String()),
			Required("age", Int()),
			Required("department", String()),
		))),
		Required("address", Object(
			Required("street", String()),
			Required("city", String()),
			Optional("zip", String()),
		)),
		Required("rating", Float()),
	)
}

func TestIntegration_SimpleValid(t *testing.T) {
	data := loadTestData(t, "simple_valid.yaml")
	result, err := Validate(data, personSchema())
	if err != nil {
		t.Fatal(err)
	}
	if !result.OK {
		t.Errorf("simple_valid.yaml should pass, got errors: %v", result.Errors)
	}
}

func TestIntegration_SimpleMissingRequired(t *testing.T) {
	data := loadTestData(t, "simple_missing_required.yaml")
	result, err := Validate(data, personSchema())
	if err != nil {
		t.Fatal(err)
	}
	if result.OK {
		t.Error("simple_missing_required.yaml should fail")
	}
	if len(result.Errors) != 1 {
		t.Fatalf("expected 1 error, got %d: %v", len(result.Errors), result.Errors)
	}
	e := result.Errors[0]
	if e.Path != "age" {
		t.Errorf("expected path 'age', got %q", e.Path)
	}
	if e.Message != `missing required field "age"` {
		t.Errorf("unexpected message: %s", e.Message)
	}
}

func TestIntegration_SimpleWrongType(t *testing.T) {
	data := loadTestData(t, "simple_wrong_type.yaml")
	result, err := Validate(data, personSchema())
	if err != nil {
		t.Fatal(err)
	}
	if result.OK {
		t.Error("simple_wrong_type.yaml should fail")
	}
	if len(result.Errors) != 1 {
		t.Fatalf("expected 1 error, got %d: %v", len(result.Errors), result.Errors)
	}
	if result.Errors[0].Path != "age" {
		t.Errorf("expected path 'age', got %q", result.Errors[0].Path)
	}
}

func TestIntegration_NestedValid(t *testing.T) {
	data := loadTestData(t, "nested_valid.yaml")
	result, err := Validate(data, personWithAddressSchema())
	if err != nil {
		t.Fatal(err)
	}
	if !result.OK {
		t.Errorf("nested_valid.yaml should pass, got errors: %v", result.Errors)
	}
}

func TestIntegration_NestedMissingRequired(t *testing.T) {
	data := loadTestData(t, "nested_missing_required.yaml")
	result, err := Validate(data, personWithAddressSchema())
	if err != nil {
		t.Fatal(err)
	}
	if result.OK {
		t.Error("nested_missing_required.yaml should fail")
	}
	// Missing address.city
	found := false
	for _, e := range result.Errors {
		if e.Path == "address.city" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected error at 'address.city', got: %v", result.Errors)
	}
}

func TestIntegration_ArrayValid(t *testing.T) {
	schema := Object(
		Required("name", String()),
		Optional("tags", Array(String())),
	)
	data := loadTestData(t, "array_valid.yaml")
	result, err := Validate(data, schema)
	if err != nil {
		t.Fatal(err)
	}
	if !result.OK {
		t.Errorf("array_valid.yaml should pass, got errors: %v", result.Errors)
	}
}

func TestIntegration_ArrayWrongItemType(t *testing.T) {
	schema := Object(
		Required("name", String()),
		Optional("tags", Array(String())),
	)
	data := loadTestData(t, "array_wrong_item_type.yaml")
	result, err := Validate(data, schema)
	if err != nil {
		t.Fatal(err)
	}
	if result.OK {
		t.Error("array_wrong_item_type.yaml should fail")
	}
	if len(result.Errors) != 2 {
		t.Fatalf("expected 2 errors, got %d: %v", len(result.Errors), result.Errors)
	}
}

func TestIntegration_ComplexValid(t *testing.T) {
	data := loadTestData(t, "complex_valid.yaml")
	result, err := Validate(data, companySchema())
	if err != nil {
		t.Fatal(err)
	}
	if !result.OK {
		t.Errorf("complex_valid.yaml should pass, got errors: %v", result.Errors)
	}
}

func TestIntegration_ComplexErrors(t *testing.T) {
	data := loadTestData(t, "complex_errors.yaml")
	result, err := Validate(data, companySchema())
	if err != nil {
		t.Fatal(err)
	}
	if result.OK {
		t.Error("complex_errors.yaml should fail")
	}

	paths := make(map[string]bool)
	for _, e := range result.Errors {
		paths[e.Path] = true
		// Verify all errors have line numbers > 0
		if e.Line <= 0 {
			t.Errorf("error at %q has no line number: %+v", e.Path, e)
		}
	}

	expected := []string{"name", "active", "employees[0].age", "employees[1].department", "address.city", "rating"}
	for _, p := range expected {
		if !paths[p] {
			t.Errorf("expected error at path %q, not found", p)
		}
	}
}

func TestIntegration_DeeplyNested(t *testing.T) {
	schema := Object(
		Required("config", Object(
			Required("database", Object(
				Required("host", String()),
				Required("port", Int()),
				Required("credentials", Object(
					Required("username", String()),
					Required("password", String()),
				)),
			)),
			Required("cache", Object(
				Required("enabled", Bool()),
				Required("ttl", Int()),
			)),
		)),
	)
	data := loadTestData(t, "deeply_nested.yaml")
	result, err := Validate(data, schema)
	if err != nil {
		t.Fatal(err)
	}
	if !result.OK {
		t.Errorf("deeply_nested.yaml should pass, got errors: %v", result.Errors)
	}
}

func TestIntegration_ArrayOfObjects(t *testing.T) {
	schema := Object(
		Required("users", Array(Object(
			Required("name", String()),
			Required("email", String()),
			Optional("roles", Array(String())),
		))),
	)
	data := loadTestData(t, "array_of_objects.yaml")
	result, err := Validate(data, schema)
	if err != nil {
		t.Fatal(err)
	}
	if !result.OK {
		t.Errorf("array_of_objects.yaml should pass, got errors: %v", result.Errors)
	}
}

func TestIntegration_Empty(t *testing.T) {
	data := loadTestData(t, "empty.yaml")
	result, err := Validate(data, Object(Required("anything", String())))
	if err != nil {
		t.Fatal(err)
	}
	if result.OK {
		t.Error("empty.yaml should fail validation")
	}
}

func TestIntegration_NullValues(t *testing.T) {
	data := loadTestData(t, "null_values.yaml")
	schema := Object(
		Required("name", String()),
		Required("age", Int()),
	)
	result, err := Validate(data, schema)
	if err != nil {
		t.Fatal(err)
	}
	if result.OK {
		t.Error("null_values.yaml should fail — name is null")
	}
	found := false
	for _, e := range result.Errors {
		if e.Path == "name" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected error for null name, got: %v", result.Errors)
	}
}

func TestIntegration_ExtraFields(t *testing.T) {
	data := loadTestData(t, "extra_fields.yaml")
	schema := Object(
		Required("name", String()),
		Required("age", Int()),
	)
	result, err := Validate(data, schema)
	if err != nil {
		t.Fatal(err)
	}
	if result.OK {
		t.Error("extra_fields.yaml should fail — has 'nickname' and 'hobby' not in schema")
	}
	if len(result.Errors) != 2 {
		t.Fatalf("expected 2 errors, got %d: %v", len(result.Errors), result.Errors)
	}
	paths := make(map[string]bool)
	for _, e := range result.Errors {
		paths[e.Path] = true
		if e.Line <= 0 {
			t.Errorf("error at %q has no line number", e.Path)
		}
	}
	if !paths["nickname"] {
		t.Error("expected error for 'nickname'")
	}
	if !paths["hobby"] {
		t.Error("expected error for 'hobby'")
	}
}

func TestIntegration_ExtraFieldsNested(t *testing.T) {
	data := loadTestData(t, "extra_fields_nested.yaml")
	schema := Object(
		Required("name", String()),
		Required("age", Int()),
		Required("address", Object(
			Required("street", String()),
			Required("city", String()),
		)),
	)
	result, err := Validate(data, schema)
	if err != nil {
		t.Fatal(err)
	}
	if result.OK {
		t.Error("extra_fields_nested.yaml should fail — has extra 'country' and 'state' in address")
	}
	if len(result.Errors) != 2 {
		t.Fatalf("expected 2 errors, got %d: %v", len(result.Errors), result.Errors)
	}
	paths := make(map[string]bool)
	for _, e := range result.Errors {
		paths[e.Path] = true
	}
	if !paths["address.country"] {
		t.Error("expected error for 'address.country'")
	}
	if !paths["address.state"] {
		t.Error("expected error for 'address.state'")
	}
}

// TestIntegration_ErrorMessages verifies error messages are human-readable
func TestIntegration_ErrorMessages(t *testing.T) {
	yaml := "name: 42\nage: hello\nitems: notarray"
	schema := Object(
		Required("name", String()),
		Required("age", Int()),
		Required("items", Array(String())),
	)
	result, err := Validate([]byte(yaml), schema)
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range result.Errors {
		if e.Message == "" {
			t.Errorf("error at %q has empty message", e.Path)
		}
		if e.Path == "" {
			t.Errorf("error has empty path: %+v", e)
		}
		t.Logf("  %s (line %d): %s", e.Path, e.Line, e.Message)
	}
}
