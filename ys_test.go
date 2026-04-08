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
