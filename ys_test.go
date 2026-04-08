package ys

import (
	"testing"
)

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
	// Just verify the schema constructors don't panic
	_ = String()
	_ = Int()
	_ = Float()
	_ = Bool()
	_ = Any()
	_ = Object()
	_ = Array(String())
}
