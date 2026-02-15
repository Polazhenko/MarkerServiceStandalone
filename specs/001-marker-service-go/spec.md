# Feature Specification: MarkerService Go Port

**Feature Branch**: `001-marker-service-go`  
**Created**: 2026-02-15  
**Status**: Draft  
**Input**: User description: "Port MarkerService microservice from C# .NET to Go, replicating all CRUD operations, category management, geolocation storage, and RESTful API with Swagger documentation using in-memory storage"

## User Scenarios & Testing

### User Story 1 - Create and Manage Markers (Priority: P1)

A user with a numeric user ID needs to create markers at specific geographic locations (latitude/longitude) with a category and description. The system assigns each marker a unique ID based on its category and creation timestamp, returning it to the user for reference.

**Why this priority**: Creating markers is the core functionality of the service. All other features depend on having markers in the system. This is the critical MVP feature.

**Independent Test**: Can be fully tested by calling the Create Marker endpoint, verifying the marker is stored with correct geolocation data, category, and a uniquely generated ID.

**Acceptance Scenarios**:

1. **Given** a user with ID 1, **When** posting a CreateMarkerRequest with name "Deer Sighting", category "GroupCategory1_Marker1", latitude 45.5, longitude -93.2, **Then** the marker is created with ID format "{category}_v1.0_{unix_timestamp}", status 201, and returned to user.
2. **Given** a marker has been created, **When** creating another marker with the same category, **Then** the new marker gets a different ID due to different timestamp.
3. **Given** a user creates a marker, **When** no description provided, **Then** the marker is created successfully with null description field.

---

### User Story 2 - Retrieve Markers by Category (Priority: P1)

A user needs to retrieve their markers with optional filtering by category. The system returns all markers for a user, optionally filtered to a specific category.

**Why this priority**: Retrieval is equally critical as creation. Users need to query their data. Without this, markers created in US1 cannot be retrieved, making the service non-functional.

**Independent Test**: Can be fully tested by creating markers with different categories and filtering retrieval by category, verifying correct subset is returned.

**Acceptance Scenarios**:

1. **Given** user 1 has markers in multiple categories, **When** retrieving with `GET /api/v1/markers?category=2`, **Then** only markers with category 2 are returned.
2. **Given** user 1 has markers, **When** retrieving without category parameter, **Then** all user 1's markers are returned.
3. **Given** user 2 and user 1 both have markers, **When** user 1 retrieves their markers, **Then** only user 1's markers are returned, not user 2's.
4. **Given** a user with no markers, **When** retrieving, **Then** empty list is returned.

---

### User Story 3 - Update Marker Details (Priority: P2)

A user needs to update a marker's location (latitude/longitude), name, and description after creation. The system allows partial updates without requiring all fields.

**Why this priority**: Updates allow users to correct mistakes or provide additional information. Less critical than CRUD basics but important for user experience.

**Independent Test**: Can be fully tested by creating a marker, updating specific fields, and verifying only changed fields are updated while others remain unchanged.

**Acceptance Scenarios**:

1. **Given** marker exists, **When** updating only the name, **Then** name changes, latitude/longitude/description remain unchanged.
2. **Given** marker exists, **When** updating latitude and longitude, **Then** geolocation updates, other fields unchanged.
3. **Given** marker exists, **When** clearing description (null value), **Then** description becomes null.
4. **Given** marker exists, **When** attempting to update category, **Then** category is NOT changed (category is immutable per .NET implementation).

---

### User Story 4 - Delete Markers (Priority: P2)

A user needs to remove markers they no longer need. The system deletes the marker and returns success/not found status appropriately.

**Why this priority**: Deletion is important for data management but less critical than creation and retrieval for MVP.

**Independent Test**: Can be fully tested by creating a marker, deleting it, and verifying it no longer exists in retrieval.

**Acceptance Scenarios**:

1. **Given** marker exists, **When** deleting it, **Then** HTTP 204 No Content is returned.
2. **Given** marker has been deleted, **When** attempting to retrieve it, **Then** HTTP 404 Not Found is returned.
3. **Given** marker for user 1 exists, **When** user 2 attempts to delete it, **Then** marker is not found (user isolation enforced).
4. **Given** non-existent marker, **When** attempting to delete, **Then** HTTP 404 Not Found is returned.

---

### User Story 5 - API Documentation and Discoverability (Priority: P3)

Developers need to discover and understand all available API endpoints without manual documentation hunting. The service provides interactive API documentation.

**Why this priority**: Documentation is important for usability but not critical for MVP functionality. Users can still use the API with manual testing or separate documentation.

**Independent Test**: Can be fully tested by accessing the Swagger UI endpoint and verifying all CRUD operations are documented with proper request/response schemas.

**Acceptance Scenarios**:

1. **Given** API is running, **When** accessing `/swagger`, **Then** interactive Swagger UI is displayed.
2. **Given** Swagger UI is open, **When** viewing endpoint definitions, **Then** all CRUD endpoints (Create, Get, List, Update, Delete) are documented.
3. **Given** Swagger documentation, **When** examining request schemas, **Then** correct required and optional fields are shown for each endpoint.

---

### Edge Cases

- What happens when a marker is created for a user who has never been seen before? (New user entry should be created)
- How does the system handle concurrent marker creation from the same user? (Thread-safe concurrent dictionary should handle this)
- What happens when retrieving markers with an invalid category filter? (Should return only exact matches or empty set)
- What happens when location coordinates are at boundary values (±180°, ±90°)? (Should accept valid coordinate ranges)
- What happens if a marker update request has no fields to update? (Should return unchanged marker or accept empty update)
- How does the system respond when user ID is missing or invalid in the header? (Should reject with 400 Bad Request)

## Requirements

### API Contract

**Base URL**: `http://localhost:8080` (Go service, port TBD during planning)

**Endpoints**:

- **POST** `/api/v1/markers` - Create a new marker
  - Header: `X-User-Id: <int>` (required)
  - Request: `{ "name": "string", "category": <int>, "latitude": <float>, "longitude": <float>, "description": "string (optional)" }`
  - Response (201): `{ "id": "string", "name": "string", "category": <int>, "latitude": <float>, "longitude": <float>, "description": "string or null", "created_at": "ISO8601", "updated_at": "ISO8601" }`
  
- **GET** `/api/v1/markers/{markerId}` - Get a specific marker
  - Header: `X-User-Id: <int>` (required)
  - Response (200): Marker object
  - Response (404): `{ "error": "Marker not found" }`
  
- **GET** `/api/v1/markers` - List markers for authenticated user
  - Header: `X-User-Id: <int>` (required)
  - Query params: `category=<int>` (optional)
  - Response (200): `[ { marker object }, ... ]`
  
- **PATCH** `/api/v1/markers/{markerId}` - Update marker
  - Header: `X-User-Id: <int>` (required)
  - Request: `{ "name": "string (optional)", "latitude": <float (optional)>, "longitude": <float (optional)>, "description": "string (optional)" }`
  - Response (200): Updated marker object
  - Response (404): `{ "error": "Marker not found" }`
  
- **DELETE** `/api/v1/markers/{markerId}` - Delete marker
  - Header: `X-User-Id: <int>` (required)
  - Response (204): No Content
  - Response (404): `{ "error": "Marker not found" }`

### Functional Requirements

- **FR-001**: System MUST accept HTTP requests at `/api/v1/markers` endpoints
- **FR-002**: System MUST require `X-User-Id` header (integer) on all requests and enforce user isolation
- **FR-003**: System MUST generate MarkerID in format `{MarkerCategory}_v1.0_{UnixTimestamp}` where UnixTimestamp is creation time
- **FR-004**: System MUST support four marker categories: General (1), GroupCategory1_Marker1 (2), GroupCategory1_Marker2 (3), GroupCategory2_Marker1 (4)
- **FR-005**: System MUST store latitude/longitude as floating-point numbers with standard geographic coordinate ranges (±180° longitude, ±90° latitude)
- **FR-006**: System MUST return HTTP 201 Created for marker creation with Location header
- **FR-007**: System MUST return HTTP 404 Not Found when requesting non-existent markers
- **FR-008**: System MUST return HTTP 204 No Content for successful deletion (no body required)
- **FR-009**: System MUST support partial updates via PATCH - only provided fields should be updated
- **FR-010**: System MUST NOT allow category updates after creation (category is immutable)
- **FR-011**: System MUST maintain in-memory storage with user isolation (nested map: UserID → MarkerID → Marker)
- **FR-012**: System MUST provide RESTful endpoint design with correct HTTP methods and status codes
- **FR-013**: System MUST handle errors gracefully with JSON error messages `{ "error": "description" }`
- **FR-014**: System MUST support filtering markers by category via optional query parameter
- **FR-015**: System MUST provide Swagger/OpenAPI documentation endpoint at `/swagger` or similar
- **FR-016**: System MUST use ISO 8601 format for all timestamps (created_at, updated_at)
- **FR-017**: System MUST be async/concurrent-safe for in-memory operations

### Key Entities

- **Marker**: Represents a geotagged point of interest with name, category, coordinates, timestamp metadata, and user association. Attributes: id (generated), userId (from header), name, category, latitude, longitude, description (optional), createdAt, updatedAt.
- **MarkerCategory**: Enumeration defining allowed categories (General, GroupCategory1_Marker1, GroupCategory1_Marker2, GroupCategory2_Marker1). Used for marker classification and ID generation.
- **User**: Implicit entity represented only by numeric ID from `X-User-Id` header. No profile data stored, only used for data isolation.

## Success Criteria

### Measurable Outcomes

- **SC-001**: All five CRUD operations (Create, Read/Get, List, Update, Delete) complete successfully and return correct HTTP status codes (201, 200, 204, 404)
- **SC-002**: API responds to all requests in under 100ms average latency (measured on local machine with in-memory storage)
- **SC-003**: System maintains user isolation - user 1 cannot see or modify user 2's markers
- **SC-004**: Marker geolocation data (latitude/longitude) persists accurately with floating-point precision maintained
- **SC-005**: Category filtering correctly returns only markers matching the specified category
- **SC-006**: Swagger/OpenAPI documentation is accessible and documents all endpoints with correct request/response schemas
- **SC-007**: Concurrent requests from multiple users can be processed without data corruption or race conditions
- **SC-008**: Error responses contain clear, actionable error messages in JSON format

## Assumptions

- **Authentication model**: Simple header-based user identification (X-User-Id integer) is sufficient for this standalone service; no complex authentication/authorization required
- **Data persistence**: In-memory storage is acceptable for this version; no database persistence required
- **Concurrency model**: Go's native goroutines and standard library sync utilities are sufficient; no external queueing/message bus required
- **API versioning**: Version 1 is final; no backward compatibility with future versions required
- **Coordinate ranges**: Standard geographic coordinate ranges (±180° longitude, ±90° latitude) are sufficient validation
- **Floating-point precision**: Standard IEEE 754 float precision is acceptable for geolocation coordinates
- **Timestamp format**: ISO 8601 UTC timestamps are acceptable for all date/time fields
- **Rate limiting**: Not required for MVP; can be added in future versions
- **HTTPS/TLS**: Not required for MVP; HTTP is acceptable for this standalone service
