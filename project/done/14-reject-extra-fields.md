# Task 14: Reject Extra Fields Not in Schema

## Description
When validating a YAML document against an object schema, any keys present in the YAML
that are not defined in the schema should produce a validation error. Only fields explicitly
defined in the schema (via Required() or Optional()) should be accepted.

## Acceptance Criteria
- Extra keys in YAML objects produce SchemaError with path, line, and message
- Works at all nesting levels (root, nested objects, objects in arrays)
- Error message clearly states the field is unknown
- All existing tests continue to pass
- New unit tests and integration tests cover the feature

## Approach (TDD)
1. Write failing tests for extra field detection
2. Implement the check in objectSchema.validate()
3. Update any existing tests that rely on extra fields being silently ignored
4. Add integration test with testdata YAML files
5. Verify all tests pass
