# Research & Technical Decisions: MarkerService Go Port

**Date**: 2026-02-15  
**Feature**: MarkerService Go Port  
**Phase**: 0 (Research & unknowns resolution)

## Decision 1: Go Web Framework Selection

**Decision**: Use `github.com/gin-gonic/gin` as the HTTP router and middleware framework.

**Rationale**:
- Gin is lightweight and fast, comparable to the minimal overhead of ASP.NET Core routing
- Excellent middleware support matching .NET's middleware pipeline pattern
- Built-in JSON marshaling/unmarshaling (like .NET's System.Text.Json)
- Swagger integration via `swaggo/gin-swagger` (equivalent to Swashbuckle in .NET)
- Active community with good documentation
- Minimal dependencies reduce maintenance burden

**Alternatives Considered**:
- **Echo**: Similar features to Gin; slightly more verbose middleware syntax
- **Chi**: Minimalist; less built-in functionality; would require more custom code
- **Gorilla Mux**: More established but heavier; overkill for this use case
- **Standard library http**: Would require manual middleware composition; feasible but verbose

**Result**: Gin selected for balance of simplicity, features, and ecosystem maturity.

---

## Decision 2: In-Memory Storage Implementation

**Decision**: Use nested `map[int]map[string]Marker` with `sync.RWMutex` for thread-safe access.

**Rationale**:
- Go's `sync` package provides efficient concurrency primitives
- RWMutex allows concurrent reads with exclusive write access (matches typical workload)
- Nested map structure mirrors the hierarchical relationship: UserID → MarkerID → Marker
- Zero external dependencies (uses Go standard library only)
- Simpler than Repository pattern (aligns with Constitution Principle III)
- Easy to understand and debug; explicit state management

**Alternatives Considered**:
- **Repository pattern + interface abstraction**: Would add complexity per Constitution; rejected
- **ConcurrentMap wrapper**: Unnecessary abstraction layer; Go's sync primitives sufficient
- **sync.Map**: Optimized for specific patterns; less suitable for nested structure

**Result**: Direct map + RWMutex selected for simplicity and performance.

---

## Decision 3: MarkerID Generation Strategy

**Decision**: Generate MarkerID as `{MarkerCategory}_{v1.0}_{UnixTimestampSeconds}` at creation time.

**Rationale**:
- Matches .NET implementation exactly (ensures parity across language ports)
- Monotonically increasing timestamps ensure uniqueness within same category
- Embeds version info for forward compatibility
- Can be regenerated if category changes (per Constitution Principle IV)
- Human-readable format aids debugging

**Alternatives Considered**:
- **UUID**: Would lose semantic meaning; harder to correlate with creation time
- **Incrementing counter**: Would require distributed consensus for multi-instance; unnecessary for MVP
- **Timestamp only (no category prefix)**: Would lose category information; less useful for debugging

**Result**: Selected format maintains compatibility and semantic clarity.

---

## Decision 4: JSON Marshaling and Timestamp Format

**Decision**: Use ISO 8601 UTC format for all timestamps in API responses (e.g., "2026-02-15T13:44:53.457Z").

**Rationale**:
- ISO 8601 is the standard for JSON APIs and REST services
- Go's `time.Time` marshals to ISO 8601 by default when using standard `encoding/json`
- Matches .NET's default JSON serialization behavior
- Timezone-aware; unambiguous interpretation across regions
- Supported by all major JavaScript/frontend libraries

**Alternatives Considered**:
- **Unix timestamp (seconds)**: Less readable; requires client-side conversion
- **Unix timestamp (milliseconds)**: Higher precision but rarely needed for marker timestamps
- **Custom format**: Would complicate client integration

**Result**: ISO 8601 UTC selected for clarity and standard interoperability.

---

## Decision 5: Testing Strategy

**Decision**: Implement three test layers per Constitution Principle II:
1. **Unit tests**: Test service layer functions (ID generation, filtering, validation)
2. **Contract tests**: Test HTTP endpoints with request/response validation (using testify)
3. **Integration tests**: Test full request flow through handlers → service → repository

**Rationale**:
- Constitution Principle II (TDD) mandates tests written first
- Three-layer approach matches .NET project structure
- Contract tests validate API stability
- Integration tests ensure end-to-end functionality
- Go's standard testing package sufficient; testify adds assertion convenience

**Alternatives Considered**:
- **Only unit tests**: Would miss API contract issues; risky for API-first design
- **Only integration tests**: Would miss granular logic bugs; slower feedback
- **BDD framework (Gherkin)**: Unnecessary complexity for this project; standard tests sufficient

**Result**: Three-layer testing approach selected with TDD mandatory per Constitution.

---

## Decision 6: Swagger/OpenAPI Implementation

**Decision**: Use `swaggo/swag` code-generation tool with `swaggo/gin-swagger` for automatic documentation generation.

**Rationale**:
- `swag` parses Go comments and generates OpenAPI specs (like Swagger annotations in .NET)
- Zero runtime overhead; documentation generated at build time
- Automatically stays synchronized with code via comments
- Endpoint `/swagger/index.html` provides interactive UI (matches .NET Swagger experience)
- Standard Go ecosystem practice

**Alternatives Considered**:
- **Manual OpenAPI YAML**: Would require manual updates; easy to fall out of sync
- **go-restful + swagger**: More heavyweight; overkill for single service
- **No automated documentation**: Violates Constitution Principle I (API-First Design)

**Result**: swaggo selected for automatic synchronization and minimal overhead.

---

## Decision 7: HTTP Status Codes and Error Handling

**Decision**: Return standard REST status codes with JSON error responses:
- 201 Created (POST marker)
- 200 OK (GET, PATCH marker/list)
- 204 No Content (DELETE marker)
- 400 Bad Request (validation errors)
- 404 Not Found (marker/user not found)
- 500 Internal Server Error (unexpected errors)

**Rationale**:
- Matches .NET implementation exactly
- Enables client-side error handling via status code
- JSON error body provides human-readable details
- No custom error codes; standard REST semantics

**Alternatives Considered**:
- **Custom HTTP status codes**: Would complicate clients; no advantage
- **All responses 200 with error in body**: Breaks HTTP semantics; complicates proxies/load balancers
- **No error bodies**: Would force clients to infer errors; poor experience

**Result**: Standard REST status codes + JSON error body selected per specification.

---

## Decision 8: Concurrency Safety for User Isolation

**Decision**: Use nested maps with RWMutex to enforce user isolation at storage layer; no additional access control layer needed.

**Rationale**:
- X-User-Id header extracted in middleware (single point of truth)
- Repository operations are user-scoped (takes UserID parameter)
- Service layer assumes valid UserID from middleware
- No need for complex authorization; simple header-based scoping sufficient
- Matches .NET in-memory repository pattern

**Alternatives Considered**:
- **RBAC/middleware auth layer**: Unnecessary complexity for header-based user ID
- **Async/await at repository level**: Go's concurrency model (goroutines) handles naturally
- **Optimistic locking**: Unnecessary for in-memory single-instance service

**Result**: Simple user-scoped repository operations selected; middleware validates header presence.

---

## Decision 9: Request/Response Serialization and Optional Fields

**Decision**: Use pointers for optional fields in Go structs (e.g., `Description *string`) to properly serialize null vs missing.

**Rationale**:
- Go's JSON marshaler treats zero values (empty string) as meaningful
- Pointers allow distinction between "not provided" and "empty string"
- Matches .NET's nullable reference behavior
- Cleaner than custom JSON marshaling logic

**Example**:
```go
type CreateMarkerRequest struct {
    Name        string
    Category    int
    Latitude    float32
    Longitude   float32
    Description *string // Optional; nil means not provided
}
```

**Alternatives Considered**:
- **Custom UnmarshalJSON**: Would require verbose boilerplate per struct
- **`omitempty` struct tags**: Doesn't distinguish null from missing on unmarshal
- **json.RawMessage**: Unnecessary complexity; pointers simpler

**Result**: Pointer-based optional fields selected for simplicity and clarity.

---

## Decision 10: Module and Package Naming

**Decision**: Go module path: `github.com/marker-service-go` (or same as repository name)
Package structure follows domain organization: `models`, `repository`, `service`, `handlers`

**Rationale**:
- Clear separation of concerns matches layered architecture
- Package names describe their responsibility (not acronyms)
- Easy to onboard new developers; self-documenting structure
- Follows Go conventions (golang.org/doc/effective_go)

**Alternatives Considered**:
- **Flat package structure**: Would mix concerns; harder to navigate as project grows
- **Nested packages (`internal/repository`, etc.)**: Go's `internal/` pattern is useful for larger projects; not needed here
- **Feature-based organization**: Makes less sense for single-service with single domain

**Result**: Flat, responsibility-based package organization selected.

---

## Open Questions (All Resolved)

None remaining. All technical decisions documented with rationale.

---

## Summary Table

| Area | Decision | Tool/Library | Rationale |
|------|----------|--------------|-----------|
| HTTP Framework | Gin | `gin-gonic/gin` | Lightweight, good middleware, Swagger support |
| Storage | Map + RWMutex | Go stdlib `sync` | Simple, thread-safe, no dependencies |
| Testing | Three-layer TDD | stdlib `testing` + `testify` | Constitutional requirement + comprehensive coverage |
| Swagger | Code generation | `swaggo/swag` | Automatic sync, zero runtime overhead |
| Error handling | REST status codes | stdlib `net/http` | Standard semantics, JSON error bodies |
| Serialization | Pointers for optional | Go stdlib `encoding/json` | Clear null vs missing distinction |

