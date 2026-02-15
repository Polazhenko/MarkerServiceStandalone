# Quick Start Guide: MarkerService Go Port

**Date**: 2026-02-15  
**Purpose**: Hands-on walkthrough for testing the MarkerService Go implementation

## Prerequisites

- Go 1.21 or later installed
- cURL or Postman for HTTP requests
- Text editor or IDE
- MarkerService Go repository cloned and built

## Setup

### 1. Build the Project

```bash
cd marker-service-go
go mod download
go build -o bin/marker-service
```

### 2. Run the Server

```bash
./bin/marker-service
```

Expected output:
```
[GIN-debug] Listening and serving HTTP on :8080
```

Server is now running on `http://localhost:8080`

## Testing Workflow

### Test 1: Create a Marker (User Story 1)

**Objective**: Create a marker and verify unique ID generation

```bash
curl -X POST http://localhost:8080/api/v1/markers \
  -H "Content-Type: application/json" \
  -H "X-User-Id: 1" \
  -d '{
    "name": "Deer Sighting",
    "category": 2,
    "latitude": 45.5,
    "longitude": -93.2,
    "description": "Large buck spotted near north fence"
  }'
```

**Expected Response** (HTTP 201):
```json
{
  "id": "GroupCategory1_Marker1_v1.0_1707999893",
  "name": "Deer Sighting",
  "category": 2,
  "latitude": 45.5,
  "longitude": -93.2,
  "description": "Large buck spotted near north fence",
  "created_at": "2026-02-15T13:44:53Z",
  "updated_at": "2026-02-15T13:44:53Z"
}
```

**Verification**:
- [ ] Status code is 201 Created
- [ ] Response includes unique `id` with format `{category}_v1.0_{timestamp}`
- [ ] Marker name, category, coordinates match request
- [ ] Timestamps are ISO 8601 format

---

### Test 2: Retrieve All Markers (User Story 2)

**Objective**: Retrieve markers for a user without filtering

```bash
curl -X GET http://localhost:8080/api/v1/markers \
  -H "X-User-Id: 1"
```

**Expected Response** (HTTP 200):
```json
[
  {
    "id": "GroupCategory1_Marker1_v1.0_1707999893",
    "name": "Deer Sighting",
    "category": 2,
    "latitude": 45.5,
    "longitude": -93.2,
    "description": "Large buck spotted near north fence",
    "created_at": "2026-02-15T13:44:53Z",
    "updated_at": "2026-02-15T13:44:53Z"
  }
]
```

**Verification**:
- [ ] Status code is 200 OK
- [ ] Response is an array of markers
- [ ] All markers belong to User 1
- [ ] Marker count increases with each creation

---

### Test 3: Filter Markers by Category (User Story 2 - Advanced)

**Objective**: Create multiple markers with different categories and filter

```bash
# Create marker in category 3
curl -X POST http://localhost:8080/api/v1/markers \
  -H "Content-Type: application/json" \
  -H "X-User-Id: 1" \
  -d '{
    "name": "Weather Station",
    "category": 3,
    "latitude": 46.0,
    "longitude": -93.5,
    "description": "Automated weather monitoring"
  }'

# Retrieve only category 2 markers
curl -X GET "http://localhost:8080/api/v1/markers?category=2" \
  -H "X-User-Id: 1"
```

**Expected Response**:
- Only markers with `category: 2` returned (ignores category 3 marker)
- Array length is 1 (only the first marker)

**Verification**:
- [ ] Category filter correctly returns only matching markers
- [ ] Other categories are excluded
- [ ] Status code is 200 OK

---

### Test 4: Get Specific Marker (User Story 2)

**Objective**: Retrieve a single marker by ID

```bash
# Use the ID from Test 1
curl -X GET http://localhost:8080/api/v1/markers/GroupCategory1_Marker1_v1.0_1707999893 \
  -H "X-User-Id: 1"
```

**Expected Response** (HTTP 200):
```json
{
  "id": "GroupCategory1_Marker1_v1.0_1707999893",
  "name": "Deer Sighting",
  "category": 2,
  "latitude": 45.5,
  "longitude": -93.2,
  "description": "Large buck spotted near north fence",
  "created_at": "2026-02-15T13:44:53Z",
  "updated_at": "2026-02-15T13:44:53Z"
}
```

**Verification**:
- [ ] Status code is 200 OK
- [ ] Returned marker matches requested ID
- [ ] All fields present and correct

---

### Test 5: Update Marker Details (User Story 3)

**Objective**: Update marker name and coordinates

```bash
curl -X PATCH http://localhost:8080/api/v1/markers/GroupCategory1_Marker1_v1.0_1707999893 \
  -H "Content-Type: application/json" \
  -H "X-User-Id: 1" \
  -d '{
    "name": "Deer Sighting (Updated)",
    "latitude": 45.6,
    "longitude": -93.1
  }'
```

**Expected Response** (HTTP 200):
```json
{
  "id": "GroupCategory1_Marker1_v1.0_1707999893",
  "name": "Deer Sighting (Updated)",
  "category": 2,
  "latitude": 45.6,
  "longitude": -93.1,
  "description": "Large buck spotted near north fence",
  "created_at": "2026-02-15T13:44:53Z",
  "updated_at": "2026-02-15T13:45:00Z"
}
```

**Verification**:
- [ ] Status code is 200 OK
- [ ] Name updated to "(Updated)"
- [ ] Latitude and longitude changed
- [ ] Description unchanged
- [ ] `updated_at` timestamp changed
- [ ] `created_at` unchanged
- [ ] Marker ID unchanged

---

### Test 6: Partial Update with Optional Fields (User Story 3)

**Objective**: Update only description, leaving coordinates unchanged

```bash
curl -X PATCH http://localhost:8080/api/v1/markers/GroupCategory1_Marker1_v1.0_1707999893 \
  -H "Content-Type: application/json" \
  -H "X-User-Id: 1" \
  -d '{
    "description": "Buck sighting confirmed; two additional deer observed"
  }'
```

**Expected Response** (HTTP 200):
```json
{
  "id": "GroupCategory1_Marker1_v1.0_1707999893",
  "name": "Deer Sighting (Updated)",
  "category": 2,
  "latitude": 45.6,
  "longitude": -93.1,
  "description": "Buck sighting confirmed; two additional deer observed",
  "created_at": "2026-02-15T13:44:53Z",
  "updated_at": "2026-02-15T13:45:15Z"
}
```

**Verification**:
- [ ] Only description changed
- [ ] Name, category, coordinates unchanged
- [ ] `updated_at` incremented

---

### Test 7: Delete Marker (User Story 4)

**Objective**: Delete a marker and verify it no longer exists

```bash
curl -X DELETE http://localhost:8080/api/v1/markers/GroupCategory1_Marker1_v1.0_1707999893 \
  -H "X-User-Id: 1"
```

**Expected Response** (HTTP 204 No Content):
- No body
- Status code 204

```bash
# Verify marker is gone
curl -X GET http://localhost:8080/api/v1/markers/GroupCategory1_Marker1_v1.0_1707999893 \
  -H "X-User-Id: 1"
```

**Expected Response** (HTTP 404):
```json
{
  "error": "Marker not found"
}
```

**Verification**:
- [ ] DELETE returns 204 No Content
- [ ] GET after DELETE returns 404 Not Found
- [ ] Error message is descriptive

---

### Test 8: User Isolation (User Story 2 - Security)

**Objective**: Verify that users cannot see other users' markers

```bash
# Create marker for User 1 (from earlier tests)
# Now try to retrieve it as User 2
curl -X GET http://localhost:8080/api/v1/markers \
  -H "X-User-Id: 2"
```

**Expected Response** (HTTP 200):
```json
[]
```

Empty array - User 2 has no markers.

```bash
# Try to access User 1's specific marker as User 2
curl -X GET http://localhost:8080/api/v1/markers/GroupCategory1_Marker1_v1.0_1707999893 \
  -H "X-User-Id: 2"
```

**Expected Response** (HTTP 404):
```json
{
  "error": "Marker not found"
}
```

Not 403 Forbidden - follows REST convention of 404 for "user doesn't have this resource"

**Verification**:
- [ ] User 2's marker list is empty
- [ ] User 2 cannot access User 1's markers (404, not 403)
- [ ] User isolation is enforced at repository layer

---

### Test 9: Swagger Documentation (User Story 5)

**Objective**: Access interactive API documentation

```bash
open http://localhost:8080/swagger/index.html
# or
curl http://localhost:8080/swagger/index.html
```

**Expected Result**:
- [ ] Swagger UI loads successfully
- [ ] All CRUD endpoints visible (POST, GET, PATCH, DELETE /api/v1/markers)
- [ ] Request/response schemas documented
- [ ] X-User-Id header requirement visible

---

### Test 10: Error Handling

**Objective**: Verify error responses are informative

```bash
# Missing required header
curl -X GET http://localhost:8080/api/v1/markers

# Invalid category in filter
curl -X GET "http://localhost:8080/api/v1/markers?category=999" \
  -H "X-User-Id: 1"

# Invalid request body (missing required field)
curl -X POST http://localhost:8080/api/v1/markers \
  -H "Content-Type: application/json" \
  -H "X-User-Id: 1" \
  -d '{
    "name": "Incomplete Marker"
  }'
```

**Expected Response**:
- Missing header: 400 Bad Request with error message
- Invalid category: 400 Bad Request or ignored
- Missing field: 400 Bad Request with validation error

**Verification**:
- [ ] Error responses have HTTP status code (400, 404, 500)
- [ ] Error body contains `error` field with message
- [ ] Error messages are user-friendly

---

## Concurrent Operations Test

**Objective**: Verify thread-safe concurrent marker operations

```bash
# Simulate concurrent creates (run in parallel shells)
for i in {1..5}; do
  curl -X POST http://localhost:8080/api/v1/markers \
    -H "Content-Type: application/json" \
    -H "X-User-Id: 1" \
    -d "{
      \"name\": \"Concurrent Marker $i\",
      \"category\": $((RANDOM % 4 + 1)),
      \"latitude\": $((RANDOM % 180 - 90)),
      \"longitude\": $((RANDOM % 360 - 180))
    }" &
done
wait

# Verify all markers created
curl -X GET http://localhost:8080/api/v1/markers \
  -H "X-User-Id: 1" | jq 'length'
```

**Expected Result**:
- [ ] All 5 markers created successfully
- [ ] No duplicate IDs (timestamps ensure uniqueness)
- [ ] Marker count is 5

---

## Test Checklist (MVP Validation)

Use this checklist to verify MarkerService Go is ready for deployment:

- [ ] **US1 Create Markers**: Unique ID generation, all fields stored correctly
- [ ] **US1 Create Markers**: No description field works (optional)
- [ ] **US2 Retrieve All**: GET /api/v1/markers returns correct list
- [ ] **US2 Retrieve with Filter**: Category filtering works correctly
- [ ] **US2 User Isolation**: Different users see only their own markers
- [ ] **US3 Update Details**: PATCH modifies specified fields only
- [ ] **US3 Partial Updates**: Requests with some optional fields work
- [ ] **US4 Delete**: DELETE returns 204 and marker is gone
- [ ] **US5 Documentation**: Swagger UI accessible at /swagger
- [ ] **Error Handling**: 400/404/500 responses with error messages
- [ ] **Concurrency**: Parallel creates don't cause corruption
- [ ] **Latency**: All operations complete in <100ms

---

## Performance Baseline

Run these commands to establish baseline performance:

```bash
# Single request latency
time curl -X GET http://localhost:8080/api/v1/markers \
  -H "X-User-Id: 1" > /dev/null

# Throughput test (100 sequential requests)
for i in {1..100}; do
  curl -X GET http://localhost:8080/api/v1/markers \
    -H "X-User-Id: 1" > /dev/null
done &
time wait
```

**Expected Performance**:
- Single request: <50ms
- 100 sequential requests: <5 seconds average
- No memory growth or leaks after 1000 operations

---

## Troubleshooting

| Issue | Diagnosis | Solution |
|-------|-----------|----------|
| Port 8080 already in use | `lsof -i :8080` or `netstat -tulpn \| grep 8080` | Change port in main.go or kill process |
| 404 on all endpoints | Server not running or wrong port | Check server logs; verify `go run main.go` working |
| Markers persist between server restarts | Data persisted to disk | Verify in-memory map (no file I/O) in code |
| Missing Swagger UI | swaggo not installed | Run `go install github.com/swaggo/swag/cmd/swag@latest` |
| Concurrent tests failing | Race condition in code | Run `go test -race` to detect data races |

---

## Next Steps

After validation:

1. **Run full test suite**: `go test ./...`
2. **Run race detector**: `go test -race ./...`
3. **Build for deployment**: `go build -o bin/marker-service`
4. **Monitor performance**: Add request logging and timing
5. **Consider future enhancements**:
   - Persistent storage (PostgreSQL)
   - Authentication/authorization
   - Rate limiting
   - Request logging and tracing

