# Test Documentation

## Running Tests

### Run all tests
```bash
go test ./...
```

### Run tests with coverage
```bash
go test -cover ./...
```

### Run tests with verbose output
```bash
go test -v ./...
```

### Run tests for a specific package
```bash
go test ./internal/services
go test ./internal/api/handlers
go test ./internal/middleware
```

### Generate coverage report
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Test Structure

### Unit Tests
- `auth_service_test.go` - Authentication service tests
- `auth_handler_test.go` - Authentication HTTP handler tests
- `auth_test.go` - JWT middleware tests

### Mock Objects
Tests use mock implementations to avoid database dependencies:
- `MockUserRepository` - In-memory user repository
- `MockAuthService` - Mock authentication service
- `MockOrgRepository` - Mock organization repository

## Test Coverage

Current test coverage focuses on critical paths:
- ✅ User registration
- ✅ User login
- ✅ JWT token generation
- ✅ JWT token validation
- ✅ Authentication middleware
- ✅ Health check endpoint

## Adding New Tests

When adding new features, create corresponding test files:

```go
package yourpackage

import "testing"

func TestYourFunction(t *testing.T) {
    // Arrange
    input := "test"
    expected := "expected"
    
    // Act
    result := YourFunction(input)
    
    // Assert
    if result != expected {
        t.Errorf("Expected %s but got %s", expected, result)
    }
}
```

## Test Database (Future)

For integration tests that require a database:

1. Use `testcontainers-go` for PostgreSQL
2. Or use an in-memory SQLite database
3. Run migrations before tests
4. Clean up after each test

Example:
```go
func TestWithDB(t *testing.T) {
    // Setup test database
    db := setupTestDB(t)
    defer db.Close()
    
    // Run tests
    // ...
}
```

## Continuous Integration

Tests should run on:
- Every commit
- Every pull request
- Before deployment

GitHub Actions example:
```yaml
name: Tests
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
        with:
          go-version: 1.22
      - run: go test -v ./...
```

## Test Best Practices

1. **Use table-driven tests** for multiple scenarios
2. **Mock external dependencies** (database, APIs)
3. **Test error cases** as well as success cases
4. **Keep tests fast** - use mocks instead of real services
5. **Test public APIs** not internal implementation
6. **One assertion per test** when possible
7. **Clear test names** that describe what they test

## Future Test Additions

Priority tests to add:
- [ ] Twilio webhook handler tests
- [ ] OpenAI service tests (with mocked API)
- [ ] Conversation service tests
- [ ] Repository tests (with test database)
- [ ] Stripe service tests (with mocked Stripe)
- [ ] End-to-end API tests
