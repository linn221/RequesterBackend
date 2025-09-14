# RequesterBackend

A Go-based REST API for recording HTTP requests, managing notes, attachments, programs, and endpoints for bug hunting and security analysis.

## Overview

This backend service provides a comprehensive API for managing security testing workflows, allowing users to:
- Record and analyze HTTP requests
- Manage bug bounty programs and their endpoints
- Attach files and notes to various resources
- Import HAR files for bulk request analysis
- Track job progress for long-running operations

## Architecture

The application follows a clean architecture pattern with clear separation of concerns:

```
├── api/           # HTTP layer (routes, middleware)
├── handlers/      # HTTP handlers (controllers)
├── services/      # Business logic layer
├── models/        # Data models and database entities
├── utils/         # Utility functions and helpers
└── openapi.yaml   # API specification
```

## Key Features

### 🎯 **Program Management**
- Create, read, update, delete bug bounty programs
- Track program scope, domains, and notes
- Associate endpoints and requests with programs

### 🔗 **Endpoint Management**
- Manage API endpoints with method, URI, and type information
- Link endpoints to programs
- Track endpoint-specific notes and attachments

### 📝 **Request Recording**
- Record HTTP requests and responses
- Store headers, body, status codes, and timing information
- Generate hashes for request/response deduplication
- Support for bulk import via HAR files

### 📎 **Notes & Attachments**
- Add notes to programs, endpoints, and requests
- Upload and manage file attachments
- Polymorphic relationships for flexible note/attachment linking

### 🔄 **Job Management**
- Track long-running import jobs
- Monitor progress of HAR file imports
- Support for different job types (import_har, import_xml)

## Data Models

### Program
Represents a bug bounty program or security testing project.

```go
type Program struct {
    Id        int       `gorm:"primaryKey"`
    Name      string    `gorm:"size:255;not null"`
    URL       string    `gorm:"size:500"`
    Notes     string    `gorm:"type:text"`
    Scope     string    `gorm:"type:text"`
    Domains   string    `gorm:"type:text"`
    CreatedAt time.Time `gorm:"autoCreateTime"`
    UpdatedAt time.Time `gorm:"autoUpdateTime"`
    
    // Associations
    Endpoints  []Endpoint  `gorm:"foreignKey:ProgramId"`
    Requests   []MyRequest `gorm:"foreignKey:ProgramId"`
    Attachments []Attachment `gorm:"polymorphic:Reference;polymorphicValue:programs"`
}
```

### Endpoint
Represents an API endpoint within a program.

```go
type Endpoint struct {
    Id           int          `gorm:"primaryKey"`
    ProgramId    int          `gorm:"index;not null"`
    Method       string       `gorm:"size:10;not null"`
    Domain       string       `gorm:"size:255;not null"`
    URI          string       `gorm:"type:text;not null"`
    EndpointType EndpointType `gorm:"size:20;not null;default:'API'"`
    Notes        string       `gorm:"type:text"`
    CreatedAt    time.Time    `gorm:"autoCreateTime"`
    UpdatedAt    time.Time    `gorm:"autoUpdateTime"`
    
    // Associations
    Program     *Program     `gorm:"foreignKey:ProgramId"`
    Requests    []MyRequest  `gorm:"foreignKey:EndpointId"`
    Attachments []Attachment `gorm:"polymorphic:Reference;polymorphicValue:endpoints"`
}
```

### MyRequest
Represents an HTTP request/response pair.

```go
type MyRequest struct {
    Id          int       `gorm:"primaryKey"`
    ProgramId   *int      `gorm:"index"`
    ImportJobId int       `gorm:"not null;index"`
    EndpointId  int       `gorm:"not null;index"`
    Sequence    int       `gorm:"not null"`
    URL         string    `gorm:"type:text;not null"`
    Method      string    `gorm:"size:10;not null"`
    Domain      string    `gorm:"size:255;not null"`
    
    // Request data
    ReqHeaders string `gorm:"type:text"`
    ReqBody    string `gorm:"type:longtext"`
    
    // Response data
    ResStatus  int    `gorm:"not null"`
    ResHeaders string `gorm:"type:text"`
    ResBody    string `gorm:"type:longtext"`
    RespSize   int    `gorm:"not null"`
    LatencyMs  int64  `gorm:"not null"`
    
    // Metadata
    RequestTime string `gorm:"size:50"`
    ReqHash1    string `gorm:"size:64;index"`
    ReqHash     string `gorm:"size:64;index"`
    ResHash     string `gorm:"size:64;index"`
    ResBodyHash string `gorm:"size:64;index"`
    
    CreatedAt time.Time `gorm:"autoCreateTime"`
    UpdatedAt time.Time `gorm:"autoUpdateTime"`
    
    // Associations
    Attachments []Attachment `gorm:"polymorphic:Reference;polymorphicValue:requests"`
}
```

## API Endpoints

### Authentication
- `GET /start_session?secret=...` - Start a new session

### Programs
- `POST /programs` - Create a program
- `GET /programs` - List all programs
- `GET /programs/{id}` - Get program details
- `PUT /programs/{id}` - Update a program
- `DELETE /programs/{id}` - Delete a program

### Endpoints
- `POST /endpoints` - Create an endpoint
- `GET /endpoints` - List all endpoints
- `GET /endpoints/{id}` - Get endpoint details
- `PUT /endpoints/{id}` - Update an endpoint
- `DELETE /endpoints/{id}` - Delete an endpoint

### Requests
- `GET /requests` - List requests with filtering
- `GET /requests/{id}` - Get request details

### Notes
- `POST /notes` - Create a note
- `GET /notes` - List notes with filtering
- `GET /notes/{id}` - Get note details
- `PATCH /notes/{id}` - Update note value
- `DELETE /notes/{id}` - Delete a note

### Attachments
- `POST /attachments` - Upload an attachment
- `GET /attachments/{id}` - Get attachment details
- `DELETE /attachments` - Delete an attachment

### Import & Jobs
- `POST /import_har` - Import HAR file
- `GET /jobs` - List all jobs

## Getting Started

### Prerequisites
- Go 1.19 or later
- MySQL/PostgreSQL database
- Git

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd RequesterBackend
```

2. Install dependencies:
```bash
go mod tidy
```

3. Set up environment variables:
```bash
export DB_HOST=localhost
export DB_PORT=3306
export DB_USER=your_username
export DB_PASSWORD=your_password
export DB_NAME=requester_backend
export UPLOAD_DIR=./uploads
export MAX_FILE_SIZE=10485760
```

4. Run the application:
```bash
go run main.go
```

The API will be available at `http://localhost:8081`

### Database Setup

The application uses GORM for database operations. Database tables will be automatically created when the application starts.

## Development

### Project Structure

```
RequesterBackend/
├── api/
│   ├── app.go          # Application initialization
│   └── routes.go       # Route registration
├── handlers/
│   ├── handlers.go     # Generic handlers
│   ├── programHandler.go    # Program CRUD handlers
│   ├── endpointHandler.go   # Endpoint CRUD handlers
│   ├── attachmentHandler.go # Attachment handlers
│   ├── types.go        # DTOs and mapping functions
│   └── helpers.go      # Handler utilities
├── services/
│   ├── programService.go    # Program business logic
│   ├── endpointService.go   # Endpoint business logic
│   ├── attachmentService.go # Attachment business logic
│   └── helpers.go      # Service utilities
├── models/
│   ├── Program.go      # Program model
│   ├── Endpoint.go     # Endpoint model
│   ├── MyRequest.go    # Request model
│   └── Attachment.go   # Attachment model
├── utils/
│   └── httpHelper.go   # HTTP response utilities
├── main.go
├── go.mod
├── go.sum
└── openapi.yaml
```

### Code Patterns

#### Handler Pattern
Handlers follow a consistent pattern:
1. Parse and validate input
2. Call service layer
3. Handle errors
4. Return appropriate HTTP response

```go
func (h *ProgramHandler) Create(w http.ResponseWriter, r *http.Request) {
    input, err := parseJson[ProgramInput](r)
    if err != nil {
        utils.RespondError(w, err)
        return
    }

    id, err := h.Service.Create(r.Context(), input.ToModel())
    if err != nil {
        utils.RespondError(w, err)
        return
    }

    utils.OkCreated(w, id)
}
```

#### Service Pattern
Services contain business logic and database operations:
1. Validate input
2. Perform database operations
3. Return results or errors

```go
func (s *ProgramService) Create(ctx context.Context, program *models.Program) (int, error) {
    if err := s.validate(s.DB.WithContext(ctx), 0, program); err != nil {
        return 0, err
    }
    if err := s.DB.WithContext(ctx).Create(program).Error; err != nil {
        return 0, err
    }
    return program.Id, nil
}
```

#### DTO Mapping
DTOs are used for API input/output with mapping functions:

```go
func (input *ProgramInput) ToModel() *models.Program {
    return &models.Program{
        Name:    input.Name,
        URL:     input.URL,
        Scope:   input.Scope,
        Domains: input.Domains,
        Notes:   input.Note,
    }
}

func ToProgramList(program *models.Program) *ProgramList {
    return &ProgramList{
        Id:   program.Id,
        Name: program.Name,
        URL:  program.URL,
    }
}
```

### Database Relationships

The application uses GORM for ORM functionality with the following relationships:

- **Program** has many **Endpoints** and **Requests**
- **Endpoint** belongs to **Program** and has many **Requests**
- **MyRequest** belongs to **Program** and **Endpoint**
- **Attachment** has polymorphic relationships with **Program**, **Endpoint**, and **MyRequest**

### Error Handling

The application uses a consistent error handling pattern:
- Business logic errors are returned from services
- HTTP handlers convert errors to appropriate HTTP status codes
- Validation errors are returned as 400 Bad Request
- Not found errors are returned as 404 Not Found
- Server errors are returned as 500 Internal Server Error

### Validation

Input validation is handled using the `go-playground/validator` package:
- Struct tags define validation rules
- `parseJson` function automatically validates input
- Validation errors are returned as user-friendly messages

## API Documentation

Complete API documentation is available in `openapi.yaml` and can be viewed using tools like Swagger UI or Postman.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

[Add your license information here]

## Support

For support and questions, please [create an issue](link-to-issues) or contact the development team.
