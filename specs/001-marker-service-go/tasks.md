# Tasks: MarkerService Go Port

**Input**: Design documents from `/specs/001-marker-service-go/`  
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/, quickstart.md

**Tests**: Test-Driven Development per Constitution Principle II. Tests written FIRST, approved, then code implemented.

**Organization**: Tasks organized by user story (US1-US5) to enable independent implementation, testing, and MVP-first delivery.

## Format: `[ID] [P?] [Story?] Description`

- **[P]**: Task can run in parallel (different files, no dependencies on incomplete tasks)
- **[Story]**: User story label (US1, US2, US3, etc.) - REQUIRED for story phases only
- **Paths**: Absolute from project root (marker-service-go/)

---

## Phase 1: Setup (Project Initialization)

**Purpose**: Project structure, dependencies, build system

- [ ] T001 Create Go module and basic project structure with go.mod, go.sum, main.go, models/, handlers/, repository/, service/, tests/ directories per plan.md
- [ ] T002 [P] Initialize dependencies: add github.com/gin-gonic/gin, github.com/swaggo/swag, github.com/swaggo/gin-swagger, github.com/stretchr/testify to go.mod
- [ ] T003 [P] Create Makefile with targets: build, test, test-race, run, clean for project development
- [ ] T004 Create README.md with project overview, setup instructions, running the service, and API usage examples

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story implementation

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [ ] T005 Create models/errors.go defining Error and ErrorResponse types for API error responses in handler.go
- [ ] T006 Create models/marker.go defining MarkerCategory enum (General, GroupCategory1_Marker1, GroupCategory1_Marker2, GroupCategory2_Marker1) and string conversion helpers
- [ ] T007 Create models/marker.go with Marker struct (id, userId, name, category, latitude, longitude, description, createdAt, updatedAt) matching data-model.md specification
- [ ] T008 Create models/marker.go with CreateMarkerRequest, UpdateMarkerRequest, MarkerResponse structs with JSON tags and validation constraints from spec.md
- [ ] T009 Create repository/marker_repository.go interface IMarkerRepository with methods: CreateAsync, GetByIdAsync, GetByCategoryAsync, UpdateAsync, DeleteAsync per data-model.md
- [ ] T010 Create repository/marker_repository.go InMemoryMarkerRepository implementation using map[int64]map[string]Marker with sync.RWMutex for thread-safe access per data-model.md
- [ ] T011 Create service/marker_service.go interface IMarkerService with methods: CreateMarkerAsync, GetMarkerAsync, GetAllMarkersAsync, UpdateMarkerAsync, DeleteMarkerAsync per spec.md endpoints
- [ ] T012 Implement marker ID generation in service/marker_service.go: format "{Category}_v1.0_{UnixTimestamp}" using time.Now().Unix() per research.md decision #3
- [ ] T013 Create handlers/middleware.go with ExtractUserID middleware to parse X-User-Id header, validate positive integer, and inject into request context
- [ ] T014 Create handlers/marker_handler.go with HTTP handler functions: CreateMarker, GetMarker, GetAllMarkers, UpdateMarker, DeleteMarker matching API contract in contracts/markers.openapi.yaml
- [ ] T015 Create handlers/swagger_handler.go to serve Swagger UI at /swagger using swaggo/gin-swagger
- [ ] T016 Create main.go with Gin router setup: register all handlers, middleware, add X-User-Id header validation, configure Swagger, listen on :8080

**Checkpoint**: Foundation ready - user story implementation can now begin independently in parallel

---

## Phase 3: User Story 1 - Create and Manage Markers (Priority: P1) 🎯 MVP

**Goal**: Users can create geotagged markers with unique ID generation, returning complete marker data with timestamps.

**Independent Test**: POST /api/v1/markers endpoint creates marker with correct ID format, stores all fields (name, category, coordinates, description), returns 201 Created with Location header, and subsequent GET retrieves identical marker.

### Tests (TDD: Written First)

- [ ] T017 [P] [US1] Create tests/unit/service_test.go TestGenerateMarkerID: verify ID format "{category}_v1.0_{timestamp}" with different categories and timestamps
- [ ] T018 [P] [US1] Create tests/unit/service_test.go TestValidateCreateRequest: verify validation of required fields (name, category, latitude, longitude), coordinate ranges (±90, ±180)
- [ ] T019 [P] [US1] Create tests/unit/repository_test.go TestCreateMarker: verify in-memory storage creates new user entry if needed, stores marker with correct ID, concurrent creates are thread-safe
- [ ] T020 [US1] Create tests/contract/markers_api_test.go TestCreateMarkerSuccess: POST /api/v1/markers returns 201, Location header set, response has all fields including generated ID
- [ ] T021 [US1] Create tests/contract/markers_api_test.go TestCreateMarkerValidation: missing required fields return 400, invalid category returns 400, coordinates out of range return 400, missing X-User-Id header returns 400
- [ ] T022 [US1] Create tests/contract/markers_api_test.go TestCreateMarkerOptionalDescription: POST without description field succeeds, response description is null (not empty string)
- [ ] T023 [US1] Create tests/integration/markers_integration_test.go TestCreateAndRetrieveMarker: full flow POST /api/v1/markers then GET /api/v1/markers/{id} returns identical marker

### Implementation

- [ ] T024 [US1] Implement MarkerService.CreateMarkerAsync in service/marker_service.go: validate input, generate ID via GenerateMarkerID, set userId from context, set createdAt/updatedAt to time.Now().UTC(), call repository.CreateAsync
- [ ] T025 [US1] Implement InMemoryMarkerRepository.CreateAsync in repository/marker_repository.go: acquire write lock, create user entry if needed (map[userId]), store marker with ID as key, return marker, ensure concurrent-safe with RWMutex
- [ ] T026 [US1] Implement CreateMarker handler in handlers/marker_handler.go: extract X-User-Id, parse JSON request body, validate, call service.CreateMarkerAsync, return 201 with Location header and marshaled MarkerResponse
- [ ] T027 [US1] Add request validation to handlers/marker_handler.go CreateMarker: verify latitude ∈ [-90, 90], longitude ∈ [-180, 180], name not empty, category is valid enum value per FR-005, FR-004
- [ ] T028 [US1] Implement marker timestamp serialization in models/marker.go: use json.RawMessage or custom MarshalJSON for ISO 8601 UTC format (e.g., "2026-02-15T13:44:53Z") per data-model.md serialization

---

## Phase 4: User Story 2 - Retrieve Markers by Category (Priority: P1) 🎯 MVP

**Goal**: Users can retrieve all their markers or filter by category, with enforced user isolation so users only see their own data.

**Independent Test**: GET /api/v1/markers endpoint returns all markers for authenticated user, GET /api/v1/markers?category=N returns only markers matching category N, GET as different user returns empty set (user isolation verified).

### Tests (TDD: Written First)

- [ ] T029 [P] [US2] Create tests/unit/repository_test.go TestGetByCategoryAsync: verify returns all markers when category=null, returns only matching category when specified, returns empty list for user with no markers, is thread-safe with concurrent reads
- [ ] T030 [P] [US2] Create tests/unit/repository_test.go TestUserIsolation: verify user 1 cannot see user 2's markers, GetByCategoryAsync(userId=2) returns empty even if user 1 has markers
- [ ] T031 [US2] Create tests/contract/markers_api_test.go TestGetAllMarkers: GET /api/v1/markers returns 200 with array of markers, response is empty array for user with no markers
- [ ] T032 [US2] Create tests/contract/markers_api_test.go TestGetMarkersFilterByCategory: GET /api/v1/markers?category=2 returns only category 2 markers, GET ?category=999 returns empty (invalid category)
- [ ] T033 [US2] Create tests/contract/markers_api_test.go TestUserIsolation: create markers for user 1, GET /api/v1/markers with X-User-Id=2 returns empty array (not user 1's markers)
- [ ] T034 [US2] Create tests/contract/markers_api_test.go TestGetMarkerById: GET /api/v1/markers/{id} returns marker when it exists and belongs to user, returns 404 when marker doesn't exist or belongs to different user
- [ ] T035 [US2] Create tests/integration/markers_integration_test.go TestCreateMultipleAndFilter: create 5 markers with mixed categories, filter by each category, verify correct subsets returned, user isolation maintained

### Implementation

- [ ] T036 [US2] Implement MarkerService.GetAllMarkersAsync in service/marker_service.go: call repository.GetByCategoryAsync with userId and optional category filter, map results to MarkerResponse DTOs, return collection
- [ ] T037 [US2] Implement InMemoryMarkerRepository.GetByCategoryAsync in repository/marker_repository.go: acquire read lock, return all markers for userId if category null, filter by category if provided, return empty list if user not found
- [ ] T038 [US2] Implement MarkerService.GetMarkerAsync in service/marker_service.go: call repository.GetByIdAsync(markerId, userId), map to MarkerResponse if found, return null if not found or user doesn't own marker
- [ ] T039 [US2] Implement InMemoryMarkerRepository.GetByIdAsync in repository/marker_repository.go: acquire read lock, lookup markers[userId][markerId], return null if not found, ensure user isolation
- [ ] T040 [US2] Implement GetAllMarkers handler in handlers/marker_handler.go: extract X-User-Id and optional category query param, validate category if provided, call service.GetAllMarkersAsync, return 200 with marshaled array
- [ ] T041 [US2] Implement GetMarker handler in handlers/marker_handler.go: extract X-User-Id and markerId path param, call service.GetMarkerAsync, return 200 if found or 404 with error JSON if not found
- [ ] T042 [US2] Add category validation in handlers/marker_handler.go: parse category query param as integer, validate is in [1,2,3,4], reject invalid categories with 400 or silently ignore per error handling decision

---

## Phase 5: User Story 3 - Update Marker Details (Priority: P2)

**Goal**: Users can partially update marker name, coordinates, and description (category is immutable), with proper change tracking via updatedAt timestamp.

**Independent Test**: PATCH /api/v1/markers/{id} with only name field updates only name while preserving other fields; similar for coordinates and description; category updates are rejected or ignored; updatedAt timestamp changes but createdAt remains fixed.

### Tests (TDD: Written First)

- [ ] T043 [P] [US3] Create tests/unit/service_test.go TestPartialUpdate: update only name field, verify other fields unchanged; update only latitude/longitude, verify others unchanged; update only description, verify others unchanged
- [ ] T044 [P] [US3] Create tests/unit/service_test.go TestUpdateCategoryImmutable: attempt to update category field is ignored/rejected, category remains unchanged, marker ID unchanged
- [ ] T045 [P] [US3] Create tests/unit/service_test.go TestUpdateTimestamps: verify updatedAt changes, createdAt unchanged, both in ISO 8601 UTC format
- [ ] T046 [P] [US3] Create tests/unit/repository_test.go TestUpdateMarker: verify partial field updates in repository, merge request fields into existing marker, update timestamps, thread-safe with RWMutex
- [ ] T047 [US3] Create tests/contract/markers_api_test.go TestPatchMarkerPartial: PATCH with only name succeeds, returns 200, other fields unchanged in response
- [ ] T048 [US3] Create tests/contract/markers_api_test.go TestPatchMarkerLocation: PATCH with latitude/longitude succeeds, coordinates updated, description/name unchanged
- [ ] T049 [US3] Create tests/contract/markers_api_test.go TestPatchMarkerDescription: PATCH with description null clears field, PATCH with description string updates field, other fields unchanged
- [ ] T050 [US3] Create tests/contract/markers_api_test.go TestPatchCategoryImmutable: PATCH with category field in request is rejected (400 or ignored), category unchanged, ID unchanged
- [ ] T051 [US3] Create tests/contract/markers_api_test.go TestPatchNotFound: PATCH non-existent marker returns 404, PATCH marker owned by different user returns 404
- [ ] T052 [US3] Create tests/integration/markers_integration_test.go TestUpdateWorkflow: create marker, update name, update location, update description, verify each change persists and others remain unchanged

### Implementation

- [ ] T053 [US3] Implement MarkerService.UpdateMarkerAsync in service/marker_service.go: validate marker exists for user, apply partial updates (only provided fields), ignore category if present, update updatedAt, call repository.UpdateAsync, return updated MarkerResponse
- [ ] T054 [US3] Implement field-by-field update logic in service/marker_service.go: if request.Name != null update it, if Latitude != null update it, if Longitude != null update it, if Description != null update it, explicitly skip Category field
- [ ] T055 [US3] Implement InMemoryMarkerRepository.UpdateAsync in repository/marker_repository.go: acquire write lock, lookup markers[userId][id], merge fields from input marker, update updatedAt timestamp, return merged marker
- [ ] T056 [US3] Implement UpdateMarker handler in handlers/marker_handler.go: extract X-User-Id and markerId, parse UpdateMarkerRequest JSON, call service.UpdateMarkerAsync, return 200 with updated marker or 404 if not found
- [ ] T057 [US3] Add coordinate validation to handlers/marker_handler.go UpdateMarker: if latitude provided verify ∈ [-90, 90], if longitude provided verify ∈ [-180, 180], reject invalid with 400 per FR-005

---

## Phase 6: User Story 4 - Delete Markers (Priority: P2)

**Goal**: Users can delete their markers, receiving 204 No Content response, with proper 404 handling when marker doesn't exist or belongs to another user.

**Independent Test**: DELETE /api/v1/markers/{id} returns 204 No Content for existing marker, subsequent GET returns 404; DELETE non-existent marker returns 404; DELETE marker owned by different user returns 404; user isolation enforced.

### Tests (TDD: Written First)

- [ ] T058 [P] [US4] Create tests/unit/repository_test.go TestDeleteMarker: verify marker deleted, subsequent get returns null, user isolation enforced, thread-safe with RWMutex
- [ ] T059 [US4] Create tests/contract/markers_api_test.go TestDeleteMarkerSuccess: DELETE existing marker returns 204 No Content with empty body
- [ ] T060 [US4] Create tests/contract/markers_api_test.go TestDeleteMarkerNotFound: DELETE non-existent marker returns 404 with error JSON
- [ ] T061 [US4] Create tests/contract/markers_api_test.go TestDeleteUserIsolation: delete marker for user 1, attempt delete as user 2 returns 404 (not 403, per REST convention)
- [ ] T062 [US4] Create tests/contract/markers_api_test.go TestDeleteThenGetReturns404: POST marker, DELETE it, GET same marker returns 404
- [ ] T063 [US4] Create tests/integration/markers_integration_test.go TestDeleteWorkflow: create multiple markers, delete one, verify deleted marker gone, others remain

### Implementation

- [ ] T064 [US4] Implement MarkerService.DeleteMarkerAsync in service/marker_service.go: verify marker exists for user, call repository.DeleteAsync, return true if success or false if not found
- [ ] T065 [US4] Implement InMemoryMarkerRepository.DeleteAsync in repository/marker_repository.go: acquire write lock, verify markers[userId][id] exists, delete from map, return true if deleted or false if not found
- [ ] T066 [US4] Implement DeleteMarker handler in handlers/marker_handler.go: extract X-User-Id and markerId, call service.DeleteMarkerAsync, return 204 No Content if success or 404 if not found
- [ ] T067 [US4] Add error response for DeleteMarker handler: return 404 with JSON error body (not 403) when marker not found or belongs to different user per REST convention in contracts/markers.openapi.yaml

---

## Phase 7: User Story 5 - API Documentation and Discoverability (Priority: P3)

**Goal**: Developers can discover and test API endpoints via interactive Swagger/OpenAPI UI with accurate request/response schemas and all CRUD operations documented.

**Independent Test**: GET /swagger or /swagger/index.html returns HTML with interactive Swagger UI; all 5 endpoints (POST/GET/PATCH/DELETE) visible with complete request/response schemas; X-User-Id header requirement visible; try-it-out functionality works.

### Tests (TDD: Written First)

- [ ] T068 [US5] Create tests/contract/swagger_test.go TestSwaggerEndpointExists: GET /swagger returns 200 with content-type text/html
- [ ] T069 [US5] Create tests/contract/swagger_test.go TestSwaggerUILoads: GET /swagger/index.html returns valid HTML containing swagger-ui references
- [ ] T070 [US5] Create tests/contract/swagger_test.go TestSwaggerSchemaComplete: GET /swagger/doc.json returns OpenAPI schema with all 5 endpoints, correct request/response models
- [ ] T071 [US5] Create tests/contract/swagger_test.go TestSwaggerDocumentation: Swagger spec includes X-User-Id header as required parameter on all endpoints, response examples match API contract

### Implementation

- [ ] T072 [US5] Generate Swagger documentation: install swaggo tools, add swag init command to Makefile, generate docs/ directory with swagger.json and swagger.yaml from code comments
- [ ] T073 [US5] Add Swagger comment annotations to main.go: document API title, version, contact info, license per swaggo format
- [ ] T074 [US5] Add Swagger endpoint handler in handlers/swagger_handler.go: register /swagger and /swagger/index.html routes, serve swaggo ginSwagger handler with custom swagger.json location
- [ ] T075 [US5] Add Swagger doc comments to marker handlers: document each endpoint (POST/GET/PATCH/DELETE) with description, parameters, request body, response codes, examples per swaggo/gin-swagger format
- [ ] T076 [US5] Add Swagger model comments to models/marker.go: document Marker, MarkerResponse, CreateMarkerRequest, UpdateMarkerRequest structs with field descriptions and constraints per swaggo format
- [ ] T077 [US5] Generate final Swagger documentation: run `swag init`, verify swagger.json created in docs/, verify /swagger endpoint serves UI correctly with all endpoints documented

---

## Phase 8: Polish & Cross-Cutting Concerns

**Purpose**: Production readiness, code quality, performance, documentation

- [ ] T078 [P] Implement comprehensive error handling in handlers: all handlers return appropriate 400/404/500 status codes with consistent JSON error format `{"error": "message"}`
- [ ] T079 [P] Implement request validation middleware: validate X-User-Id header presence and format on all requests, return 400 if missing/invalid per FR-002
- [ ] T080 [P] Add request/response logging: add structured logging to handlers showing request method/path, status code, latency, user ID for debugging
- [ ] T081 [P] Implement concurrent safety testing: run `go test -race` to detect data races in repository with RWMutex, fix any races found
- [ ] T082 Add integration test for full API flow: create marker (US1) → retrieve (US2) → update (US3) → delete (US4) → verify gone, all in single test per quickstart.md Test scenario flow
- [ ] T083 Create performance baseline tests: measure average latency for each endpoint (create, read, update, delete) on local machine, verify <100ms per SC-002 success criterion
- [ ] T084 Add concurrent operation test: spawn 10 goroutines creating markers simultaneously, verify no data corruption, all markers created successfully with unique IDs
- [ ] T085 Update Makefile with test targets: `make test` runs all tests, `make test-race` runs with race detector, `make coverage` generates coverage reports
- [ ] T086 Add code documentation: update README.md with API usage examples from quickstart.md, setup instructions, development workflow
- [ ] T087 Create deployment documentation: document how to build release binary, environment variables, configuration options, running on different platforms

---

## Implementation Strategy

### MVP Scope (Recommended Start)

**Minimum Viable Product**: Complete User Stories 1 & 2 (both P1)
- Phase 1 Setup: Initialize project (T001-T004)
- Phase 2 Foundational: Build all core infrastructure (T005-T016)
- Phase 3 US1: Create markers with unique IDs (T017-T028)
- Phase 4 US2: Retrieve markers with filtering and user isolation (T029-T042)

**Delivery**: Runnable API supporting marker creation and retrieval; testable via quickstart.md tests 1-8; full test coverage for US1 & US2.

### Incremental Delivery Plan

| Phase | User Stories | Status | Deliverable |
|-------|-------------|--------|------------|
| 1-2 | Setup + Foundation | Ready | Infrastructure in place |
| 3-4 | **US1 & US2 (P1)** | **MVP** | **Create + Read markers** |
| 5 | US3 (P2) | Enhanced | + Update markers |
| 6 | US4 (P2) | Enhanced | + Delete markers |
| 7 | US5 (P3) | Polish | + Swagger docs |
| 8 | Cross-cutting | Release | + Logging, performance, tests |

**Recommended Path**: Complete MVP (Phases 1-4) first = ~70% of core functionality in ~60% of effort. Then add update/delete (Phases 5-6) for 100% CRUD coverage.

---

## Task Dependencies

### Dependency Graph

```
T001-T004 (Setup)
    ↓
T005-T016 (Foundation)
    ↓
    ├─→ T017-T028 (US1: Create)
    │   ├─→ T029-T042 (US2: Retrieve)
    │   └─→ T043-T057 (US3: Update)
    │       └─→ T058-T067 (US4: Delete)
    │
    └─→ T068-T077 (US5: Swagger)
    
T078-T087 (Polish & QA)
```

### Critical Path

1. **Blocking**: T001-T016 (Foundation) must complete first
2. **P1 Path**: T001-T016 → T017-T042 (US1+US2 complete MVP)
3. **P2 Path**: After P1, proceed T043-T067 (US3+US4)
4. **P3 Path**: After CRUD, add T068-T077 (Documentation)
5. **QA Path**: Final T078-T087 (Polish)

### Parallel Opportunities

**Within Phase 1 (Setup)**:
- T002, T003 can run parallel (dependencies, build system independent)

**Within Phase 2 (Foundation)**:
- T005, T006, T007, T008 can run parallel (all define models independently)
- T009, T010 parallel (repository interface + implementation separate)
- T011, T012 parallel (service interface + ID generation separate)

**Within Phase 3 (US1)**:
- T017, T018, T019 parallel (different unit tests, no interdependencies)
- T020, T021, T022, T023 parallel (different contract tests)
- T024, T025, T026, T027 parallel (implementation in service → repository → handler)

**Within Phase 4 (US2)**:
- T029, T030 parallel (repository tests for different scenarios)
- T031-T035 parallel (contract tests for different endpoints)
- T036-T042 parallel (implementation across service/repository/handler)

**Note**: Tests and implementation for same feature must follow TDD order (tests first T017-T028, then code T024-T028).

---

## Quality Checklist

Before marking each phase complete:

- [ ] All tasks in phase completed
- [ ] All tests pass: `go test ./...`
- [ ] No race conditions: `go test -race ./...`
- [ ] Code formatted: `go fmt ./...`
- [ ] Dependencies up-to-date: `go mod tidy`
- [ ] Swagger docs generated: `swag init` and /swagger accessible
- [ ] Integration tests pass: full end-to-end workflows work
- [ ] Performance baseline met: <100ms per operation
- [ ] Error handling consistent: all 4xx/5xx responses follow pattern
- [ ] Concurrent operations verified: parallel creates/reads safe

---

## Test Coverage Expected

| Phase | Unit Tests | Contract Tests | Integration Tests |
|-------|-----------|-----------------|-------------------|
| Phase 3 (US1) | T017-T019 (3) | T020-T023 (4) | T023 (1) |
| Phase 4 (US2) | T029-T030 (2) | T031-T035 (5) | T035 (1) |
| Phase 5 (US3) | T043-T046 (4) | T047-T051 (5) | T052 (1) |
| Phase 6 (US4) | T058-T059 (2) | T060-T062 (3) | T063 (1) |
| Phase 7 (US5) | - | T068-T071 (4) | - |
| **Total** | **~11 unit** | **~21 contract** | **~5 integration** |

**TDD Enforcement**: Each user story phase has tests written first (T0xx marked [US#]), then implementation (T0xx marked [US#] without T prefix numbers in 20s+ range).

---

## Format Validation

✅ All 87 tasks follow checklist format:
- ✅ Start with `- [ ]` checkbox
- ✅ Have unique Task ID (T001-T087)
- ✅ Include [P] marker where parallel
- ✅ Include [Story] label for US1-US5 phases
- ✅ Have clear description with file path

**Example formats verified**:
- ✅ `- [ ] T001 Create Go module...` (Setup, no story)
- ✅ `- [ ] T005 Create models/errors.go...` (Foundation, no story)
- ✅ `- [ ] T017 [P] [US1] Create tests/unit/service_test.go...` (US1, parallel)
- ✅ `- [ ] T024 [US1] Implement MarkerService.CreateMarkerAsync...` (US1, sequential)
- ✅ `- [ ] T078 [P] Implement comprehensive error handling...` (Polish, parallel)
