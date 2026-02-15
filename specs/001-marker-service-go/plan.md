# Implementation Plan: MarkerService Go Port

**Branch**: `001-marker-service-go` | **Date**: 2026-02-15 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/001-marker-service-go/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

Port the MarkerService microservice from C# .NET to Go, maintaining identical functionality for marker CRUD operations with geolocation storage, category management, and in-memory data persistence. The service provides a RESTful API with user isolation via X-User-Id headers and interactive Swagger documentation. This port uses Go's native concurrency primitives (goroutines, sync.RWMutex) for thread-safe operations instead of .NET's ConcurrentDictionary.

## Technical Context

**Language/Version**: Go 1.21+ (modern concurrency, error handling, native HTTP server)  
**Primary Dependencies**: 
  - `github.com/gin-gonic/gin` (HTTP router/middleware framework - lightweight, comparable to .NET's routing)
  - `github.com/swaggo/swag` + `github.com/swaggo/gin-swagger` (Swagger/OpenAPI documentation)
  - Standard library: `encoding/json`, `sync`, `time`, `net/http`
  
**Storage**: In-memory map structure with RWMutex for thread safety (no external database per Constitution Principle III)  
**Testing**: Go's standard `testing` package + `github.com/stretchr/testify` for assertions  
**Target Platform**: Linux/Mac/Windows standalone server (matches .NET standalone version)  
**Project Type**: Single web service (monolithic API server)  
**Performance Goals**: <100ms average latency for all operations (on local machine with in-memory storage)  
**Constraints**: <500MB memory footprint typical; concurrent-safe marker operations required  
**Scale/Scope**: Single aggregate (User → Markers); 5 REST endpoints; ~1000 LOC expected

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Notes |
|-----------|--------|-------|
| **I. API-First Design** | ✅ PASS | Spec includes complete API contract with 5 endpoints (POST, GET, PATCH, DELETE) and request/response schemas. JSON format confirmed. |
| **II. Test-Driven Development** | ⚠️ CONDITIONAL | Constitution requires tests written first. This will be enforced during task generation (speckit.tasks). No TDD requirement in spec, but Constitution mandates it. |
| **III. In-Memory Simplicity** | ✅ PASS | Spec explicitly requires in-memory storage. No Repository pattern, no database. Nested map structure (UserID → MarkerID → Marker) confirmed. |
| **IV. Contract Stability & Semantic Versioning** | ✅ PASS | MarkerID format includes version: `{category}_v1.0_{timestamp}`. Category is immutable (cannot be updated without ID regeneration). /api/v1 versioning in use. |
| **V. Clear Data Models** | ✅ PASS | Three entities explicitly defined: User (ID from header), Marker, MarkerCategory enum. All constraints mapped. |

**Gate Result**: ✅ **PASS** - All Constitution principles satisfied. No violations requiring justification.

## Project Structure

### Documentation (this feature)

```text
specs/001-marker-service-go/
├── plan.md                  # This file (/speckit.plan command output)
├── spec.md                  # Feature specification
├── research.md              # Phase 0 output (/speckit.plan command)
├── data-model.md            # Phase 1 output (/speckit.plan command)
├── quickstart.md            # Phase 1 output (/speckit.plan command)
├── contracts/               # Phase 1 output (/speckit.plan command)
│   ├── markers.openapi.yaml # OpenAPI schema for marker endpoints
│   └── markers.json         # JSON schema for request/response
└── tasks.md                 # Phase 2 output (/speckit.tasks command)
```

### Source Code (repository root)

```text
# Single Go project structure
marker-service-go/
├── main.go                  # Application entry point
├── go.mod                   # Module definition
├── go.sum                   # Dependency lock file
│
├── models/
│   ├── marker.go            # Marker, MarkerCategory, Request/Response types
│   └── errors.go            # Error types
│
├── handlers/
│   ├── marker_handler.go    # HTTP endpoint handlers for markers
│   ├── swagger_handler.go   # Swagger/OpenAPI documentation handler
│   └── middleware.go        # X-User-Id header validation middleware
│
├── repository/
│   └── marker_repository.go # In-memory marker storage (UserID → MarkerID → Marker)
│
├── service/
│   └── marker_service.go    # Business logic (ID generation, validation, filtering)
│
├── tests/
│   ├── unit/
│   │   ├── service_test.go
│   │   └── repository_test.go
│   ├── contract/
│   │   └── markers_api_test.go
│   └── integration/
│       └── markers_integration_test.go
│
├── docs/
│   ├── swagger.yaml         # Generated Swagger documentation
│   └── README.md            # Setup and usage instructions
│
└── Makefile                 # Build, test, run targets
```

**Structure Decision**: Single Go project structure selected. No separation into backend/frontend as per spec (API-only service). Follows Clean Architecture principles: models (entities) → repository (data access) → service (business logic) → handlers (HTTP layer).

## Complexity Tracking

> **No violations** - Constitution Check passed all gates. No justification needed.
