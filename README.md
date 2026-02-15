# MarkerService

A RESTful API service for managing geotagged markers with user isolation, category management, and in-memory data storage.
Original implementation on .Net core, ported to GoLang, using github spec.kit + Copilot agent

## Overview

MarkerService is a port of the C# .NET MarkerService microservice to Go 1.21+. It provides:

- **CRUD Operations**: Create, Read, Update, Delete markers
- **Category Management**: Organize markers by categories
- **Geolocation**: Store latitude/longitude coordinates
- **User Isolation**: Each user can only see their own markers
- **RESTful API**: Clean REST design with proper HTTP status codes
- **Swagger Documentation**: Interactive API documentation via Swagger UI
- **Thread-Safe Storage**: In-memory storage with concurrent-safe operations

## Quick Start

### Prerequisites

- Go 1.21 or later
- Make (optional, for running Makefile targets)

### Setup

```bash
# Navigate to project root
cd marker-service-go

# Download dependencies
go mod download

# Build the project
go build -o ../bin/marker-service

# Or use Make
make build
```

### Running the Service

```bash
# Run directly
go run main.go

# Or use Make
make run

# Service will start on http://localhost:8080
```

## API Endpoints

### Create Marker
```http
POST /api/v1/markers
Headers: X-User-Id: 1
Body: {
  "name": "Deer Sighting",
  "category": 2,
  "latitude": 45.5,
  "longitude": -93.2,
  "description": "Large buck spotted"
}
```

### Get All Markers
```http
GET /api/v1/markers
Headers: X-User-Id: 1
Query: ?category=2 (optional)
```

### Get Specific Marker
```http
GET /api/v1/markers/{markerId}
Headers: X-User-Id: 1
```

### Update Marker
```http
PATCH /api/v1/markers/{markerId}
Headers: X-User-Id: 1
Body: {
  "name": "Updated Name",
  "latitude": 45.6,
  "longitude": -93.1,
  "description": "Updated description"
}
```

### Delete Marker
```http
DELETE /api/v1/markers/{markerId}
Headers: X-User-Id: 1
```

## Marker Categories

The system supports four marker categories:

- **1**: General (generic/uncategorized)
- **2**: GroupCategory1_Marker1
- **3**: GroupCategory1_Marker2
- **4**: GroupCategory2_Marker1

## Marker ID Format

Marker IDs are generated automatically in the format:

```
{MarkerCategory}_v1.0_{UnixTimestamp}
```

Example: `GroupCategory1_Marker1_v1.0_1707999893`

## Testing

### Run All Tests
```bash
go test ./...
# or
make test
```

### Run Tests with Race Detector
```bash
go test -race ./...
# or
make test-race
```

### Integration Testing
The `quickstart.md` file includes 10 test scenarios that can be run using cURL to verify the API functionality.

## Project Structure

```
marker-service-go/
├── main.go                      # Application entry point
├── go.mod                       # Module definition
├── go.sum                       # Dependency lock file
├── models/
│   ├── marker.go               # Data structures and DTOs
│   └── errors.go               # Error types
├── handlers/
│   ├── marker_handler.go       # HTTP endpoints
│   ├── middleware.go           # Request middleware
│   └── swagger_handler.go      # Swagger UI handler
├── repository/
│   └── marker_repository.go    # In-memory storage
├── service/
│   └── marker_service.go       # Business logic
├── tests/
│   ├── unit/                   # Unit tests
│   ├── contract/               # API contract tests
│   └── integration/            # Integration tests
└── docs/
    └── swagger.yaml            # Generated Swagger spec
```

## API Documentation

The API is documented with Swagger/OpenAPI. Access the interactive documentation at:

```
http://localhost:8080/swagger/index.html
```

## Technical Decisions

Key technical decisions for this implementation:

1. **HTTP Framework**: Gin (lightweight, fast, good middleware support)
2. **Storage**: In-memory map with sync.RWMutex (no external database)
3. **Testing**: Three-layer approach (unit/contract/integration) with TDD
4. **Documentation**: Swagger/OpenAPI with code-generated specs
5. **Timestamps**: ISO 8601 UTC format
6. **Error Handling**: Standard REST status codes with JSON error bodies

See `specs/001-marker-service-go/research.md` for detailed rationale on each decision.

## Data Model

### User
- **ID**: Integer from X-User-Id header
- **Markers**: One-to-many relationship (each user owns multiple markers)

### Marker
- **ID**: String (auto-generated: `{category}_v1.0_{timestamp}`)
- **UserId**: User ownership reference
- **Name**: String (required)
- **Category**: Enum (1-4)
- **Latitude**: Float (-90 to +90)
- **Longitude**: Float (-180 to +180)
- **Description**: String (optional)
- **CreatedAt**: ISO 8601 UTC timestamp
- **UpdatedAt**: ISO 8601 UTC timestamp

## User Isolation

User isolation is enforced at the repository layer. Users can only:
- Create markers for themselves
- Read their own markers
- Update their own markers
- Delete their own markers

Attempting to access another user's markers returns 404 Not Found.

## Performance

- **Target Latency**: <100ms average for all operations
- **Storage**: In-memory (no database latency)
- **Concurrency**: RWMutex for thread-safe read/write operations

## Development Workflow

1. **Tests First**: Write tests before implementation (TDD)
2. **Phase-by-Phase**: Complete each phase before moving to the next
3. **Parallel Tasks**: Tasks marked with [P] can run simultaneously
4. **Quality Gates**: All tests must pass before commit

## Makefile Targets

```bash
make build          # Build the binary
make test           # Run tests
make test-race      # Run tests with race detector
make run            # Build and run
make clean          # Remove artifacts
make deps           # Download dependencies
make swagger        # Generate Swagger docs
make help           # Display help
```

## Contributing

1. Ensure all tests pass: `go test ./...`
2. Run race detector: `go test -race ./...`
3. Check code formatting: `go fmt ./...`
4. Keep Swagger docs in sync with code

## Future Enhancements

- Persistent storage (PostgreSQL)
- Authentication/authorization
- Rate limiting
- Request logging and tracing
- Caching layer
- Multi-region support

## License

MIT

## Contact

For questions or issues, refer to the specification documents in `specs/001-marker-service-go/`.
