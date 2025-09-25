# Project Conventions and Rules

This document outlines the coding conventions, patterns, and rules followed in the RequesterBackend project.

## Rules
run bash script start.sh to build and run the program after a prompt has been finished processing to save me more time

## API Response Patterns

### Create Operations (POST)
- **Service Layer**: Returns `(int, error)` - only the created resource ID
- **Handler Layer**: Uses `utils.OkCreated(w, id)` 
- **Response**: HTTP 201 with plain text ID (e.g., "123")
- **Example**: `POST /programs` returns "1" as plain text

### Update Operations (PUT/PATCH)
- **Service Layer**: Returns `(int, error)` - the updated resource ID
- **Handler Layer**: Uses `utils.OkUpdated(w)`
- **Response**: HTTP 200 with no body
- **Example**: `PUT /programs/1` returns empty 200 response

### Delete Operations (DELETE)
- **Service Layer**: Returns `(int, error)` - the deleted resource ID
- **Handler Layer**: Uses `utils.OkDeleted(w)`
- **Response**: HTTP 204 with no body
- **Example**: `DELETE /programs/1` returns empty 204 response

### Read Operations (GET)
- **Service Layer**: Returns full model objects `(*Model, error)` or `([]*Model, error)`
- **Handler Layer**: Uses `utils.OkJson(w, dto)` with appropriate DTOs
- **Response**: HTTP 200 with JSON object/array
- **Example**: `GET /programs/1` returns full program details as JSON

## Service Layer Conventions

### Method Signatures
```go
// Create operations
func (s *Service) Create(ctx context.Context, model *Model) (int, error)

// Read operations  
func (s *Service) Get(ctx context.Context, id int) (*Model, error)
func (s *Service) List(ctx context.Context, ...filters) ([]*Model, error)

// Update operations
func (s *Service) Update(ctx context.Context, id int, model *Model) (int, error)

// Delete operations
func (s *Service) Delete(ctx context.Context, id int) (int, error)
```

### Error Handling
- Always return descriptive error messages
- Use `fmt.Errorf()` for error wrapping
- Return `0` for ID when errors occur
- Check for `gorm.ErrRecordNotFound` and return appropriate errors

### Validation
- Include `validate()` methods for business logic validation
- Validate foreign key references before operations
- Check for business rules (e.g., preventing deletion of resources with children)

## Handler Layer Conventions

### Method Structure
```go
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
    // 1. Parse input
    input, err := parseJson[InputType](r)
    if err != nil {
        utils.RespondError(w, err)
        return
    }
    
    // 2. Call service
    id, err := h.Service.Create(r.Context(), input.ToModel())
    if err != nil {
        utils.RespondError(w, err)
        return
    }
    
    // 3. Return appropriate response
    utils.OkCreated(w, id)
}
```

### Response Utilities
- `utils.OkCreated(w, id)` - For create operations (201)
- `utils.OkUpdated(w)` - For update operations (200)
- `utils.OkDeleted(w)` - For delete operations (204)
- `utils.OkJson(w, data)` - For read operations (200)
- `utils.RespondError(w, err)` - For error responses

## Model Layer Conventions

### GORM Hooks
- Use `BeforeCreate` and `BeforeUpdate` hooks for automatic field generation
- Example: Auto-generating slugs, timestamps, etc.

### Polymorphic Relationships
- Use consistent polymorphic reference types: `programs`, `endpoints`, `requests`, `vulns`
- Include in validation rules: `validate:"required,oneof=programs endpoints requests vulns"`

### Field Conventions
- Use `Id` (not `ID`) for primary keys
- Use `CreatedAt` and `UpdatedAt` for timestamps
- Use nullable pointers for optional foreign keys: `*int`

## DTO Layer Conventions

### Input DTOs
- Exclude auto-generated fields (like `Id`, `Slug`, `CreatedAt`, `UpdatedAt`)
- Include validation tags
- Provide `ToModel()` methods for conversion

### List DTOs
- Include only essential fields for listing
- Use `*Type` for nullable fields

### Detail DTOs
- Include all relevant fields
- Include associated data (notes, attachments, images)
- Use proper date formatting: `"2006-01-02T15:04:05Z07:00"`

## Database Conventions

### Migration Order
- Models with no dependencies first
- Models with dependencies in dependency order
- Polymorphic models last

### Naming
- Use singular table names
- Use snake_case for column names
- Use `_id` suffix for foreign keys

## Error Handling Conventions

### HTTP Status Codes
- `200` - Successful GET/PUT operations
- `201` - Successful POST operations
- `204` - Successful DELETE operations
- `400` - Bad request (validation errors)
- `404` - Resource not found
- `500` - Internal server error

### Error Messages
- Use descriptive, user-friendly messages
- Include context when helpful
- Avoid exposing internal implementation details

## File Organization

### Directory Structure
```
handlers/     - HTTP handlers
services/     - Business logic
models/       - Database models
utils/        - Utility functions
api/          - Route registration
config/       - Configuration and migrations
```

### File Naming
- Use descriptive names: `programHandler.go`, `vulnService.go`
- Use camelCase for Go files
- Group related functionality in same file

## Code Quality Rules

### Linting
- Run `go vet ./...` before commits
- Fix all linter warnings and errors
- Use consistent formatting

### Testing
- Write tests for service layer business logic
- Test error conditions and edge cases
- Use table-driven tests where appropriate

### Documentation
- Document public APIs in OpenAPI specification
- Include examples in API documentation
- Keep README.md updated with new endpoints

## Security Conventions

### Input Validation
- Always validate input data
- Use struct validation tags
- Sanitize user input

### Error Exposure
- Don't expose internal errors to clients
- Log detailed errors server-side
- Return generic error messages to clients

## Performance Conventions

### Database Queries
- Use `Preload()` for associations
- Use batch operations for bulk inserts
- Avoid N+1 query problems

### Response Size
- Keep list responses minimal
- Use detail endpoints for full data
- Consider pagination for large datasets

---

**Note**: These conventions ensure consistency across the codebase and should be followed when adding new features or modifying existing code.
