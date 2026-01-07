# Agent Guidelines for GoKICK

This document provides guidelines for maintaining and contributing to the GoKICK project. Follow these guidelines to ensure code quality, test coverage, and documentation standards are maintained.

## Project Overview

GoKICK is a comprehensive Go client library for the [Kick.com](https://kick.com) API. It provides a type-safe, idiomatic Go interface to interact with Kick's streaming platform.

**Key Characteristics:**
- Go 1.25.5+ required
- Comprehensive test coverage (target: 96%+)
- Strict linting with golangci-lint v2.7.2
- Well-documented with examples in `docs/` directory
- Follows Go best practices and idioms

## Critical Requirements

### 1. Tests Must Always Pass

**Requirement:** All tests must pass before any code is considered complete.

**Commands:**
```bash
make test
# or
go test ./...
```

**Guidelines:**
- Run tests locally before committing
- Ensure all existing tests continue to pass when making changes
- Add tests for all new functionality
- Test files follow the pattern: `<module>_test.go`
- Use `testify` for assertions (already in dependencies)

**Test Coverage:**
- **Target:** Maintain 96%+ statement coverage
- Check coverage with: `go test ./... -coverprofile=coverage.out`
- Review coverage reports to identify untested code paths
- Add tests for edge cases, error conditions, and success paths

### 2. Linting Must Always Pass

**Requirement:** All linting checks must pass with zero issues.

**Commands:**
```bash
make lint
# or
golangci-lint run --timeout 5m --config .golangci.yml
```

**Guidelines:**
- Linting is enforced via `.golangci.yml` configuration
- Zero tolerance for linting errors
- Fix all linting issues before committing
- The project uses golangci-lint v2.7.2
- CI/CD will fail on linting errors

**Common Linters Enabled:**
- `errcheck` - Check for unchecked errors
- `govet` - Go vet checks
- `staticcheck` - Static analysis
- `revive` - Go linter
- `testifylint` - Testify-specific linting
- `funlen`, `gocyclo` - Complexity checks
- And many more (see `.golangci.yml`)

### 3. Documentation Must Always Be Updated

**Requirement:** Documentation must be kept up-to-date with code changes.

#### README.md Updates

When adding new features or changing existing functionality:
- Update the main `README.md` if the change affects:
  - Features list
  - Quick start examples
  - Installation instructions
  - API overview
  - Supported endpoints

#### Documentation Directory (`docs/`)

The `docs/` directory contains detailed API documentation:
- Each module has a corresponding markdown file (e.g., `channels.md`, `chat.md`)
- Update relevant documentation files when:
  - Adding new API methods
  - Changing method signatures
  - Adding new types or enums
  - Changing behavior or error handling

**Documentation Files:**
- `docs/README.md` - Feature checklist and overview
- `docs/authentication.md` - OAuth2 and token management
- `docs/channels.md` - Channel operations
- `docs/chat.md` - Chat message operations
- `docs/moderation.md` - Moderation tools
- `docs/livestreams.md` - Livestream data
- `docs/users.md` - User information
- `docs/categories.md` - Category browsing
- `docs/kicks.md` - Kicks leaderboard
- `docs/events.md` - Webhook event subscriptions
- `docs/webhook_events.md` - Webhook payload structures

**Documentation Guidelines:**
- Include code examples for new features
- Document all parameters and return types
- Explain error conditions
- Keep examples up-to-date with actual API usage
- Update `docs/README.md` checklist when adding endpoints

## Code Structure and Patterns

### File Organization

- **Source files:** `<module>.go` (e.g., `channels.go`, `chat.go`)
- **Test files:** `<module>_test.go` (e.g., `channels_test.go`, `chat_test.go`)
- **Enum files:** `<enum>_enum.go` and `<enum>_enum_test.go`
- **Helper files:** `helper_test.go` for shared test utilities

### Code Patterns

#### Client Methods

All API methods follow this pattern:
```go
func (c *Client) MethodName(ctx context.Context, params ...) (Response[ResultType], error) {
    return makeRequest[ResultType](ctx, c, http.MethodXxx, "/api/path", http.StatusOK, body)
}
```

#### Error Handling

- Use the `NewError` function from `error.go` for API errors
- Always return errors, never panic
- Include context in error messages
- Use `fmt.Errorf` with `%w` for error wrapping

#### Request/Response Types

- Use generic `Response[T]` type for API responses
- Define request structs for POST/PATCH requests
- Use `makeRequest` helper for standard requests
- Use `makeAuthRequest` for authentication endpoints

#### Testing Patterns

- Use table-driven tests where appropriate
- Mock HTTP responses using test helpers
- Test both success and error cases
- Test edge cases and boundary conditions
- Use `testify/assert` and `testify/require` for assertions

### Enum Patterns

Enums follow this structure:
- Type definition with constants
- String method for string representation
- Validation methods where needed
- Comprehensive tests for all enum values

## Workflow Guidelines

### Before Making Changes

1. **Pull latest changes:** `git pull`
2. **Run tests:** `make test` - ensure everything passes
3. **Run lint:** `make lint` - ensure no issues
4. **Check coverage:** Review current coverage levels

### Before Committing

1. **Run full test suite:** `make test`
2. **Run linting:** `make lint`
3. **Check test coverage:** Ensure coverage hasn't decreased
4. **Update documentation:**
   - Update relevant `docs/*.md` files
   - Update `docs/README.md` checklist if adding endpoints
   - Update `README.md` if needed
5. **Review changes:** Ensure code follows project patterns

### Code Review Checklist

When reviewing or submitting code:
- [ ] All tests pass (`make test`)
- [ ] Linting passes with zero issues (`make lint`)
- [ ] Test coverage is maintained or improved
- [ ] Documentation is updated (README, docs/)
- [ ] Code follows project patterns and conventions
- [ ] Error handling is appropriate
- [ ] No hardcoded values or secrets
- [ ] Examples in documentation are correct

## CI/CD Integration

The project uses GitHub Actions for:
- **Go tests** (`.github/workflows/go.yml`): Runs tests and uploads coverage to Coveralls
- **Linting** (`.github/workflows/golangci-lint.yml`): Runs golangci-lint on PRs

Both workflows run on pull requests and must pass before merging.

## Dependencies

- **Go:** 1.25.5 or later
- **Testing:** `github.com/stretchr/testify v1.11.1`
- **Linting:** `golangci/golangci-lint v2.7.2`

Keep dependencies minimal and up-to-date. Review `go.mod` and `go.sum` when adding dependencies.

## Common Tasks

### Adding a New API Endpoint

1. Add method to appropriate module file (e.g., `channels.go`)
2. Define request/response types if needed
3. Write comprehensive tests in `*_test.go`
4. Update relevant documentation in `docs/`
5. Update `docs/README.md` checklist
6. Run tests and lint
5. Update `README.md` if it's a major feature

### Adding a New Enum

1. Create `<enum>_enum.go` with type and constants
2. Implement `String()` method
3. Create `<enum>_enum_test.go` with tests
4. Update documentation if enum is user-facing
5. Run tests and lint

### Fixing a Bug

1. Write a test that reproduces the bug
2. Fix the bug
3. Ensure test passes
4. Run full test suite
5. Update documentation if behavior changed

## Best Practices

1. **Idiomatic Go:** Follow Go conventions and idioms
2. **Error Handling:** Always handle errors explicitly
3. **Context:** Use `context.Context` for cancellation and timeouts
4. **Concurrency:** Use mutexes appropriately (see `client.go` for examples)
5. **Documentation:** Comment exported functions and types
6. **Testing:** Aim for high coverage, test edge cases
7. **Simplicity:** Keep code simple and readable
8. **Consistency:** Follow existing code patterns

## Getting Help

- Review existing code for patterns
- Check `docs/` directory for API examples
- Review test files for usage examples
- Check `.golangci.yml` for linting rules
- Review CI/CD workflows for automated checks

## Summary

**Remember:**
- ✅ Tests must always pass
- ✅ Linting must always pass (zero issues)
- ✅ Coverage must be high (96%+)
- ✅ Documentation must always be updated
- ✅ Follow project patterns and conventions

These requirements are non-negotiable and enforced by CI/CD. Code that doesn't meet these standards will not be merged.

