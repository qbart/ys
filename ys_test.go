package ys

import (
	"os"
	"testing"
)

func readTestData(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile("testdata/" + name)
	if err != nil {
		t.Fatalf("failed to read testdata/%s: %v", name, err)
	}
	return data
}

func TestResult_EmptyDocument(t *testing.T) {
	result, err := Validate([]byte(""), Object())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.OK {
		t.Error("expected not OK for empty document")
	}
	if len(result.Errors) != 1 {
		t.Fatalf("expected 1 error, got %d", len(result.Errors))
	}
	if result.Errors[0].Message != "empty document" {
		t.Errorf("unexpected message: %s", result.Errors[0].Message)
	}
}

func TestSchema_Types(t *testing.T) {
	_ = String()
	_ = Int()
	_ = Float()
	_ = Bool()
	_ = Any()
	_ = Object()
	_ = Array(String())
}

// --- String validation ---

func TestString_Valid(t *testing.T) {
	schema := Object(Required("name", String()))
	result, err := Validate([]byte("name: hello"), schema)
	if err != nil {
		t.Fatal(err)
	}
	if !result.OK {
		t.Errorf("expected OK, got errors: %v", result.Errors)
	}
}

func TestString_Invalid(t *testing.T) {
	schema := Object(Required("name", String()))
	result, err := Validate([]byte("name: 42"), schema)
	if err != nil {
		t.Fatal(err)
	}
	if result.OK {
		t.Error("expected validation to fail")
	}
	if len(result.Errors) != 1 {
		t.Fatalf("expected 1 error, got %d", len(result.Errors))
	}
	if result.Errors[0].Path != "name" {
		t.Errorf("expected path 'name', got %q", result.Errors[0].Path)
	}
	if result.Errors[0].Message != "expected string, got int" {
		t.Errorf("unexpected message: %s", result.Errors[0].Message)
	}
}

func TestString_QuotedNumber(t *testing.T) {
	schema := Object(Required("val", String()))
	result, err := Validate([]byte(`val: "42"`), schema)
	if err != nil {
		t.Fatal(err)
	}
	if !result.OK {
		t.Errorf("quoted number should be string, got errors: %v", result.Errors)
	}
}

// --- Int validation ---

func TestInt_Valid(t *testing.T) {
	schema := Object(Required("age", Int()))
	result, err := Validate([]byte("age: 30"), schema)
	if err != nil {
		t.Fatal(err)
	}
	if !result.OK {
		t.Errorf("expected OK, got errors: %v", result.Errors)
	}
}

func TestInt_Invalid_String(t *testing.T) {
	schema := Object(Required("age", Int()))
	result, err := Validate([]byte("age: thirty"), schema)
	if err != nil {
		t.Fatal(err)
	}
	if result.OK {
		t.Error("expected validation to fail")
	}
	if result.Errors[0].Message != "expected int, got string" {
		t.Errorf("unexpected message: %s", result.Errors[0].Message)
	}
}

func TestInt_Invalid_Float(t *testing.T) {
	schema := Object(Required("age", Int()))
	result, err := Validate([]byte("age: 3.14"), schema)
	if err != nil {
		t.Fatal(err)
	}
	if result.OK {
		t.Error("expected validation to fail for float value")
	}
}

func TestInt_Negative(t *testing.T) {
	schema := Object(Required("temp", Int()))
	result, err := Validate([]byte("temp: -10"), schema)
	if err != nil {
		t.Fatal(err)
	}
	if !result.OK {
		t.Errorf("negative int should be valid, got errors: %v", result.Errors)
	}
}

// --- Float validation ---

func TestFloat_Valid(t *testing.T) {
	schema := Object(Required("rating", Float()))
	result, err := Validate([]byte("rating: 4.5"), schema)
	if err != nil {
		t.Fatal(err)
	}
	if !result.OK {
		t.Errorf("expected OK, got errors: %v", result.Errors)
	}
}

func TestFloat_Invalid(t *testing.T) {
	schema := Object(Required("rating", Float()))
	result, err := Validate([]byte("rating: not_a_number"), schema)
	if err != nil {
		t.Fatal(err)
	}
	if result.OK {
		t.Error("expected validation to fail")
	}
	if result.Errors[0].Message != "expected float, got string" {
		t.Errorf("unexpected message: %s", result.Errors[0].Message)
	}
}

func TestFloat_RejectsInt(t *testing.T) {
	schema := Object(Required("rating", Float()))
	result, err := Validate([]byte("rating: 5"), schema)
	if err != nil {
		t.Fatal(err)
	}
	// YAML tags: 5 is !!int, 5.0 is !!float — float schema should reject int
	if result.OK {
		t.Error("expected float schema to reject bare integer")
	}
}

// --- Bool validation ---

func TestBool_Valid(t *testing.T) {
	schema := Object(Required("active", Bool()))
	result, err := Validate([]byte("active: true"), schema)
	if err != nil {
		t.Fatal(err)
	}
	if !result.OK {
		t.Errorf("expected OK, got errors: %v", result.Errors)
	}
}

func TestBool_ValidFalse(t *testing.T) {
	schema := Object(Required("active", Bool()))
	result, err := Validate([]byte("active: false"), schema)
	if err != nil {
		t.Fatal(err)
	}
	if !result.OK {
		t.Errorf("expected OK, got errors: %v", result.Errors)
	}
}

func TestBool_Invalid_String(t *testing.T) {
	schema := Object(Required("active", Bool()))
	result, err := Validate([]byte(`active: "yes"`), schema)
	if err != nil {
		t.Fatal(err)
	}
	if result.OK {
		t.Error("expected validation to fail for string 'yes'")
	}
}

// --- Any validation ---

func TestAny_AcceptsString(t *testing.T) {
	schema := Object(Required("val", Any()))
	result, err := Validate([]byte("val: hello"), schema)
	if err != nil {
		t.Fatal(err)
	}
	if !result.OK {
		t.Errorf("Any() should accept string, got errors: %v", result.Errors)
	}
}

func TestAny_AcceptsInt(t *testing.T) {
	schema := Object(Required("val", Any()))
	result, err := Validate([]byte("val: 42"), schema)
	if err != nil {
		t.Fatal(err)
	}
	if !result.OK {
		t.Errorf("Any() should accept int, got errors: %v", result.Errors)
	}
}

func TestAny_AcceptsObject(t *testing.T) {
	schema := Object(Required("val", Any()))
	yaml := "val:\n  nested: true"
	result, err := Validate([]byte(yaml), schema)
	if err != nil {
		t.Fatal(err)
	}
	if !result.OK {
		t.Errorf("Any() should accept object, got errors: %v", result.Errors)
	}
}

// --- Object validation with Required/Optional ---

func TestObject_SimpleValid(t *testing.T) {
	schema := Object(
		Required("name", String()),
		Required("age", Int()),
		Optional("email", String()),
	)
	data := readTestData(t, "simple_valid.yaml")
	result, err := Validate(data, schema)
	if err != nil {
		t.Fatal(err)
	}
	if !result.OK {
		t.Errorf("expected OK, got errors: %v", result.Errors)
	}
}

func TestObject_MissingRequired(t *testing.T) {
	schema := Object(
		Required("name", String()),
		Required("age", Int()),
		Optional("email", String()),
	)
	data := readTestData(t, "simple_missing_required.yaml")
	result, err := Validate(data, schema)
	if err != nil {
		t.Fatal(err)
	}
	if result.OK {
		t.Error("expected validation to fail for missing required field")
	}
	if len(result.Errors) != 1 {
		t.Fatalf("expected 1 error, got %d: %v", len(result.Errors), result.Errors)
	}
	if result.Errors[0].Path != "age" {
		t.Errorf("expected path 'age', got %q", result.Errors[0].Path)
	}
}

func TestObject_WrongType(t *testing.T) {
	schema := Object(
		Required("name", String()),
		Required("age", Int()),
	)
	data := readTestData(t, "simple_wrong_type.yaml")
	result, err := Validate(data, schema)
	if err != nil {
		t.Fatal(err)
	}
	if result.OK {
		t.Error("expected validation to fail for wrong type")
	}
	if len(result.Errors) != 1 {
		t.Fatalf("expected 1 error, got %d: %v", len(result.Errors), result.Errors)
	}
	if result.Errors[0].Path != "age" {
		t.Errorf("expected path 'age', got %q", result.Errors[0].Path)
	}
}

func TestObject_OptionalMissing_OK(t *testing.T) {
	schema := Object(
		Required("name", String()),
		Optional("email", String()),
	)
	result, err := Validate([]byte("name: hello"), schema)
	if err != nil {
		t.Fatal(err)
	}
	if !result.OK {
		t.Errorf("missing optional field should be OK, got errors: %v", result.Errors)
	}
}

func TestObject_OptionalPresent_WrongType(t *testing.T) {
	schema := Object(
		Required("name", String()),
		Optional("count", Int()),
	)
	result, err := Validate([]byte("name: hello\ncount: nope"), schema)
	if err != nil {
		t.Fatal(err)
	}
	if result.OK {
		t.Error("optional field with wrong type should fail")
	}
}

func TestObject_ExpectedObjectGotScalar(t *testing.T) {
	schema := Object(Required("name", String()))
	result, err := Validate([]byte("just a string"), schema)
	if err != nil {
		t.Fatal(err)
	}
	if result.OK {
		t.Error("expected validation to fail when YAML is scalar not object")
	}
	if result.Errors[0].Message != "expected object, got string" {
		t.Errorf("unexpected message: %s", result.Errors[0].Message)
	}
}

// --- Nested object validation ---

func TestNested_Valid(t *testing.T) {
	schema := Object(
		Required("name", String()),
		Required("age", Int()),
		Required("address", Object(
			Required("street", String()),
			Required("city", String()),
			Optional("zip", String()),
		)),
	)
	data := readTestData(t, "nested_valid.yaml")
	result, err := Validate(data, schema)
	if err != nil {
		t.Fatal(err)
	}
	if !result.OK {
		t.Errorf("expected OK, got errors: %v", result.Errors)
	}
}

func TestNested_MissingRequiredNested(t *testing.T) {
	schema := Object(
		Required("name", String()),
		Required("age", Int()),
		Required("address", Object(
			Required("street", String()),
			Required("city", String()),
			Optional("zip", String()),
		)),
	)
	data := readTestData(t, "nested_missing_required.yaml")
	result, err := Validate(data, schema)
	if err != nil {
		t.Fatal(err)
	}
	if result.OK {
		t.Error("expected validation to fail for missing nested required field")
	}
	if len(result.Errors) != 1 {
		t.Fatalf("expected 1 error, got %d: %v", len(result.Errors), result.Errors)
	}
	if result.Errors[0].Path != "address.city" {
		t.Errorf("expected path 'address.city', got %q", result.Errors[0].Path)
	}
}

func TestNested_DeeplyNested(t *testing.T) {
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
	data := readTestData(t, "deeply_nested.yaml")
	result, err := Validate(data, schema)
	if err != nil {
		t.Fatal(err)
	}
	if !result.OK {
		t.Errorf("expected OK, got errors: %v", result.Errors)
	}
}

func TestNested_DeeplyNested_MissingField(t *testing.T) {
	schema := Object(
		Required("config", Object(
			Required("database", Object(
				Required("host", String()),
				Required("port", Int()),
				Required("credentials", Object(
					Required("username", String()),
					Required("password", String()),
					Required("token", String()),
				)),
			)),
		)),
	)
	data := readTestData(t, "deeply_nested.yaml")
	result, err := Validate(data, schema)
	if err != nil {
		t.Fatal(err)
	}
	if result.OK {
		t.Error("expected validation to fail")
	}
	found := false
	for _, e := range result.Errors {
		if e.Path == "config.database.credentials.token" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected error at path 'config.database.credentials.token', got: %v", result.Errors)
	}
}

// --- Array validation ---

func TestArray_ValidStrings(t *testing.T) {
	schema := Object(
		Required("name", String()),
		Optional("tags", Array(String())),
	)
	data := readTestData(t, "array_valid.yaml")
	result, err := Validate(data, schema)
	if err != nil {
		t.Fatal(err)
	}
	if !result.OK {
		t.Errorf("expected OK, got errors: %v", result.Errors)
	}
}

func TestArray_WrongItemType(t *testing.T) {
	schema := Object(
		Required("name", String()),
		Optional("tags", Array(String())),
	)
	data := readTestData(t, "array_wrong_item_type.yaml")
	result, err := Validate(data, schema)
	if err != nil {
		t.Fatal(err)
	}
	if result.OK {
		t.Error("expected validation to fail for wrong item types in array")
	}
	// items at index 1 (42) and 2 (true) should fail
	if len(result.Errors) != 2 {
		t.Fatalf("expected 2 errors, got %d: %v", len(result.Errors), result.Errors)
	}
	if result.Errors[0].Path != "tags[1]" {
		t.Errorf("expected path 'tags[1]', got %q", result.Errors[0].Path)
	}
	if result.Errors[1].Path != "tags[2]" {
		t.Errorf("expected path 'tags[2]', got %q", result.Errors[1].Path)
	}
}

func TestArray_Empty(t *testing.T) {
	schema := Object(
		Required("items", Array(String())),
	)
	result, err := Validate([]byte("items: []"), schema)
	if err != nil {
		t.Fatal(err)
	}
	if !result.OK {
		t.Errorf("empty array should be valid, got errors: %v", result.Errors)
	}
}

func TestArray_ExpectedArrayGotScalar(t *testing.T) {
	schema := Object(
		Required("items", Array(String())),
	)
	result, err := Validate([]byte("items: hello"), schema)
	if err != nil {
		t.Fatal(err)
	}
	if result.OK {
		t.Error("expected validation to fail when scalar given for array")
	}
	if result.Errors[0].Message != "expected array, got string" {
		t.Errorf("unexpected message: %s", result.Errors[0].Message)
	}
}

func TestArray_OfInts(t *testing.T) {
	schema := Object(
		Required("nums", Array(Int())),
	)
	result, err := Validate([]byte("nums:\n  - 1\n  - 2\n  - 3"), schema)
	if err != nil {
		t.Fatal(err)
	}
	if !result.OK {
		t.Errorf("expected OK, got errors: %v", result.Errors)
	}
}

// --- Line number tracking ---

func TestLineNumbers_WrongType(t *testing.T) {
	yaml := "name: John\nage: \"not a number\"\nemail: john@test.com"
	schema := Object(
		Required("name", String()),
		Required("age", Int()),
		Optional("email", String()),
	)
	result, err := Validate([]byte(yaml), schema)
	if err != nil {
		t.Fatal(err)
	}
	if result.OK {
		t.Fatal("expected validation to fail")
	}
	if result.Errors[0].Line != 2 {
		t.Errorf("expected error on line 2, got line %d", result.Errors[0].Line)
	}
}

func TestLineNumbers_NestedError(t *testing.T) {
	yaml := "name: Jane\naddress:\n  street: 123 Main St\n  city: 42"
	schema := Object(
		Required("name", String()),
		Required("address", Object(
			Required("street", String()),
			Required("city", String()),
		)),
	)
	result, err := Validate([]byte(yaml), schema)
	if err != nil {
		t.Fatal(err)
	}
	if result.OK {
		t.Fatal("expected validation to fail")
	}
	if result.Errors[0].Line != 4 {
		t.Errorf("expected error on line 4, got line %d", result.Errors[0].Line)
	}
}

func TestLineNumbers_ArrayItem(t *testing.T) {
	yaml := "tags:\n  - hello\n  - 42\n  - world"
	schema := Object(
		Required("tags", Array(String())),
	)
	result, err := Validate([]byte(yaml), schema)
	if err != nil {
		t.Fatal(err)
	}
	if result.OK {
		t.Fatal("expected validation to fail")
	}
	if result.Errors[0].Line != 3 {
		t.Errorf("expected error on line 3, got line %d", result.Errors[0].Line)
	}
}

func TestLineNumbers_MultipleErrors(t *testing.T) {
	yaml := "name: 123\nage: hello\nactive: 42"
	schema := Object(
		Required("name", String()),
		Required("age", Int()),
		Required("active", Bool()),
	)
	result, err := Validate([]byte(yaml), schema)
	if err != nil {
		t.Fatal(err)
	}
	if result.OK {
		t.Fatal("expected validation to fail")
	}
	if len(result.Errors) != 3 {
		t.Fatalf("expected 3 errors, got %d: %v", len(result.Errors), result.Errors)
	}
	expectedLines := []int{1, 2, 3}
	for i, e := range result.Errors {
		if e.Line != expectedLines[i] {
			t.Errorf("error %d: expected line %d, got %d", i, expectedLines[i], e.Line)
		}
	}
}

// --- Edge cases ---

func TestEdge_NullRequiredField(t *testing.T) {
	schema := Object(
		Required("name", String()),
		Required("age", Int()),
	)
	data := readTestData(t, "null_values.yaml")
	result, err := Validate(data, schema)
	if err != nil {
		t.Fatal(err)
	}
	if result.OK {
		t.Error("expected validation to fail for null required field")
	}
	foundNull := false
	for _, e := range result.Errors {
		if e.Path == "name" {
			foundNull = true
		}
	}
	if !foundNull {
		t.Errorf("expected error for null 'name' field, got: %v", result.Errors)
	}
}

func TestEdge_NullOptionalField(t *testing.T) {
	schema := Object(
		Required("age", Int()),
		Optional("name", String()),
	)
	data := readTestData(t, "null_values.yaml")
	result, err := Validate(data, schema)
	if err != nil {
		t.Fatal(err)
	}
	if !result.OK {
		t.Errorf("null optional field should be OK, got errors: %v", result.Errors)
	}
}

func TestEdge_EmptyDocument(t *testing.T) {
	schema := Object(Required("name", String()))
	data := readTestData(t, "empty.yaml")
	result, err := Validate(data, schema)
	if err != nil {
		t.Fatal(err)
	}
	if result.OK {
		t.Error("expected validation to fail for empty document")
	}
}

func TestEdge_InvalidYAML(t *testing.T) {
	_, err := Validate([]byte(":\n  :\n    - [invalid"), Object())
	if err != nil {
		// parse error is acceptable
		return
	}
	// If it parses, it should at least not panic
}

func TestEdge_ObjectWithNoFields(t *testing.T) {
	schema := Object()
	result, err := Validate([]byte("name: hello\nage: 30"), schema)
	if err != nil {
		t.Fatal(err)
	}
	if !result.OK {
		t.Errorf("object with no field constraints should accept any object, got errors: %v", result.Errors)
	}
}

func TestEdge_BoolAsString(t *testing.T) {
	// In YAML, unquoted true/false are bools. Quoted "true" is a string
	schema := Object(Required("val", String()))
	result, err := Validate([]byte("val: true"), schema)
	if err != nil {
		t.Fatal(err)
	}
	if result.OK {
		t.Error("unquoted true should be bool, not string")
	}
}

func TestEdge_BoolAsStringQuoted(t *testing.T) {
	schema := Object(Required("val", String()))
	result, err := Validate([]byte(`val: "true"`), schema)
	if err != nil {
		t.Fatal(err)
	}
	if !result.OK {
		t.Errorf("quoted 'true' should be string, got errors: %v", result.Errors)
	}
}

func TestEdge_ZeroValues(t *testing.T) {
	schema := Object(
		Required("count", Int()),
		Required("rate", Float()),
		Required("active", Bool()),
		Required("name", String()),
	)
	result, err := Validate([]byte("count: 0\nrate: 0.0\nactive: false\nname: \"\""), schema)
	if err != nil {
		t.Fatal(err)
	}
	if !result.OK {
		t.Errorf("zero values should be valid, got errors: %v", result.Errors)
	}
}

// --- Complex schemas ---

func TestComplex_ArrayOfObjects(t *testing.T) {
	schema := Object(
		Required("users", Array(Object(
			Required("name", String()),
			Required("email", String()),
			Optional("roles", Array(String())),
		))),
	)
	data := readTestData(t, "array_of_objects.yaml")
	result, err := Validate(data, schema)
	if err != nil {
		t.Fatal(err)
	}
	if !result.OK {
		t.Errorf("expected OK, got errors: %v", result.Errors)
	}
}

func TestComplex_ArrayOfObjects_Errors(t *testing.T) {
	schema := Object(
		Required("users", Array(Object(
			Required("name", String()),
			Required("email", String()),
			Required("age", Int()),
		))),
	)
	data := readTestData(t, "array_of_objects.yaml")
	result, err := Validate(data, schema)
	if err != nil {
		t.Fatal(err)
	}
	if result.OK {
		t.Error("expected validation to fail — missing 'age' in array items")
	}
	// Both users lack 'age'
	if len(result.Errors) != 2 {
		t.Fatalf("expected 2 errors, got %d: %v", len(result.Errors), result.Errors)
	}
	if result.Errors[0].Path != "users[0].age" {
		t.Errorf("expected path 'users[0].age', got %q", result.Errors[0].Path)
	}
	if result.Errors[1].Path != "users[1].age" {
		t.Errorf("expected path 'users[1].age', got %q", result.Errors[1].Path)
	}
}

func TestComplex_FullDocument_Valid(t *testing.T) {
	schema := Object(
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
	data := readTestData(t, "complex_valid.yaml")
	result, err := Validate(data, schema)
	if err != nil {
		t.Fatal(err)
	}
	if !result.OK {
		t.Errorf("expected OK, got errors: %v", result.Errors)
	}
}

func TestComplex_FullDocument_MultipleErrors(t *testing.T) {
	schema := Object(
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
		)),
		Required("rating", Float()),
	)
	data := readTestData(t, "complex_errors.yaml")
	result, err := Validate(data, schema)
	if err != nil {
		t.Fatal(err)
	}
	if result.OK {
		t.Error("expected validation to fail")
	}
	// Expect errors: name(int not string), active(string not bool),
	// employees[0].age(string not int), employees[1].department(missing),
	// address.city(missing), rating(string not float)
	if len(result.Errors) < 5 {
		t.Errorf("expected at least 5 errors, got %d: %v", len(result.Errors), result.Errors)
	}

	// Verify specific paths exist
	paths := make(map[string]bool)
	for _, e := range result.Errors {
		paths[e.Path] = true
	}
	expected := []string{"name", "active", "employees[0].age", "employees[1].department", "address.city", "rating"}
	for _, p := range expected {
		if !paths[p] {
			t.Errorf("expected error at path %q, not found in: %v", p, result.Errors)
		}
	}
}

func TestComplex_NestedArraysOfObjects(t *testing.T) {
	yaml := "teams:\n  - name: Alpha\n    members:\n      - name: Alice\n        role: lead\n      - name: Bob\n        role: dev\n  - name: Beta\n    members:\n      - name: Charlie\n        role: 42"
	schema := Object(
		Required("teams", Array(Object(
			Required("name", String()),
			Required("members", Array(Object(
				Required("name", String()),
				Required("role", String()),
			))),
		))),
	)
	result, err := Validate([]byte(yaml), schema)
	if err != nil {
		t.Fatal(err)
	}
	if result.OK {
		t.Error("expected validation to fail")
	}
	if len(result.Errors) != 1 {
		t.Fatalf("expected 1 error, got %d: %v", len(result.Errors), result.Errors)
	}
	if result.Errors[0].Path != "teams[1].members[0].role" {
		t.Errorf("expected path 'teams[1].members[0].role', got %q", result.Errors[0].Path)
	}
}
