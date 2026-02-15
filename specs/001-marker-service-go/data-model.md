# Data Model: MarkerService Go Port

**Date**: 2026-02-15  
**Feature**: MarkerService Go Port  
**Phase**: 1 (Design & Contracts)

## Overview

The MarkerService uses a simple three-entity model:
1. **User**: Identified by numeric ID (from X-User-Id header); owns multiple markers
2. **MarkerCategory**: Enumeration of valid marker types; used for classification and ID generation
3. **Marker**: Geotagged data point with ownership, category, coordinates, and timestamps

Storage structure is hierarchical: `UserID → MarkerID → Marker`

---

## Entity: User

**Responsibility**: User identity and data isolation boundary

**Fields**:

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| `id` | `int64` | Yes | Extracted from `X-User-Id` HTTP header; no persistent user record; implicit entity |

**Relationships**:
- One User owns many Markers (1:N)
- Markers are scoped to UserID; users cannot access other users' markers

**Validation Rules**:
- User ID MUST be a positive integer
- User ID MUST be present in all API requests via `X-User-Id` header

**Notes**:
- No authentication/authorization beyond header presence check
- User is implicit (created on-demand when first marker created)
- No user profile data stored; ID only

---

## Entity: MarkerCategory

**Responsibility**: Enumeration of valid marker types; used for filtering and ID generation

**Fields** (Enumeration):

| Value | Numeric Code | Description |
|-------|--------------|-------------|
| `General` | 1 | Generic/uncategorized marker |
| `GroupCategory1_Marker1` | 2 | Category group 1, marker type 1 |
| `GroupCategory1_Marker2` | 3 | Category group 1, marker type 2 |
| `GroupCategory2_Marker1` | 4 | Category group 2, marker type 1 |

**Validation Rules**:
- Only these four enumerated values are valid
- Category MUST be specified when creating a marker
- Category CANNOT be changed after marker creation (immutable)
- If category change is needed, DELETE old marker and CREATE new one with different category

**Notes**:
- String representation of enum MUST be serializable in JSON (e.g., "General")
- String representations SHOULD be case-insensitive in request parsing (normalize to proper case)
- Numeric codes enable compact storage and fast lookups

---

## Entity: Marker

**Responsibility**: Geotagged point of interest with temporal metadata and user association

### Fields

| Field | Type | Required | Constraints | Notes |
|-------|------|----------|-------------|-------|
| `id` | `string` | Yes (generated) | Format: `{Category}_v1.0_{UnixTimestamp}` | Generated at creation; immutable; used as primary key |
| `userId` | `int64` | Yes | Positive integer | Set from X-User-Id header; denormalizes user ownership |
| `name` | `string` | Yes | 1-1000 characters | Human-readable marker name |
| `category` | `MarkerCategory` | Yes | Must be valid enum value | Determines ID prefix; immutable after creation |
| `latitude` | `float32` | Yes | Range: -90.0 to +90.0 | Geographic north-south coordinate |
| `longitude` | `float32` | Yes | Range: -180.0 to +180.0 | Geographic east-west coordinate |
| `description` | `string` | No | 0-5000 characters (if provided) | Optional detailed information; nullable |
| `createdAt` | `time.Time` | Yes (generated) | ISO 8601 UTC | Timestamp of creation; used in ID generation |
| `updatedAt` | `time.Time` | Yes (generated) | ISO 8601 UTC | Timestamp of last update; updated on any field change |

### Primary Key

**Composite Key**: `(UserId, Id)` where `Id` is unique per user.

**Generation Logic** (CreatedAt):
```
MarkerID = format("{MarkerCategory}_v1.0_{UnixTimestamp}", category, now.Unix())
```

Example: `GroupCategory1_Marker1_v1.0_1707999893`

### Relationships

```
User (1)
  ↓ owns
  ↓ (many)
Marker
```

One-to-many relationship; no direct Marker-to-Marker relationships.

### Validation Rules

| Rule | Enforcement | Impact |
|------|-------------|--------|
| ID must match format `{Category}_v1.0_{UnixTimestamp}` | Repository creation | Reject malformed IDs |
| UserId must be positive | Handler/Middleware | Return 400 Bad Request |
| Name must not be empty | Service layer | Return 400 Bad Request |
| Category must be valid enum value | Service layer | Return 400 Bad Request |
| Latitude must be ±90.0 | Service layer | Return 400 Bad Request |
| Longitude must be ±180.0 | Service layer | Return 400 Bad Request |
| Description max 5000 chars | Service layer | Truncate or reject (TBD in implementation) |
| Category cannot change after creation | Service layer | Treat as DELETE old + CREATE new |
| Only marker owner can read/update/delete | Repository layer (userID scoping) | Return 404 Not Found (not 403 Forbidden, per REST convention) |

### State Transitions

```
Created (initial state via POST)
  ↓
Updated (via PATCH with name/location/description changes)
  ↓
Deleted (via DELETE endpoint)
```

**Special Case - Category Change** (PATCH with category field):
```
Existing Marker (category=A)
  ↓ PATCH with category=B
  ↓ DELETE marker with category A
  ↓ CREATE new marker with category B
  ↓ Return new Marker (with new ID)
```

Note: Per Constitution Principle IV and .NET implementation review, category changes are rare; implement as explicit delete + create.

---

## Storage Structure (In-Memory)

### Data Structure Definition

Go struct representation:

```go
// markers is a nested map: UserID → (MarkerID → Marker)
markers map[int64]map[string]Marker

// Protected by RWMutex for concurrent access
type MarkerStore struct {
    mu      sync.RWMutex
    markers map[int64]map[string]Marker // UserID → MarkerID → Marker
}
```

### Access Patterns

| Operation | Lock Type | Typical Path |
|-----------|-----------|--------------|
| Get marker | Read lock | markers[userID][markerID] |
| List user's markers | Read lock | for marker := range markers[userID] |
| Create marker | Write lock | markers[userID][newID] = marker |
| Update marker | Write lock | markers[userID][id] = updatedMarker |
| Delete marker | Write lock | delete(markers[userID], id) |

### Concurrency Guarantees

- **Reader safety**: Multiple goroutines can read simultaneously (RWMutex read lock)
- **Writer safety**: Write operations exclusive; no concurrent writes (RWMutex write lock)
- **User isolation**: Each UserID partition is independent; operations on different users don't block each other
- **Marker isolation**: Operations on different markers within a user can interleave (handled by RWMutex)

---

## Serialization

### JSON Request Format (CreateMarkerRequest)

```json
{
  "name": "Deer Sighting",
  "category": 2,
  "latitude": 45.5,
  "longitude": -93.2,
  "description": "Large buck spotted near north fence"
}
```

### JSON Response Format (MarkerResponse)

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

### Field Mapping

| Go Field | JSON Field | Type | Example |
|----------|-----------|------|---------|
| `Id` | `id` | string | "GroupCategory1_Marker1_v1.0_1707999893" |
| `UserId` | (not in response) | int64 | (internal; from header) |
| `Name` | `name` | string | "Deer Sighting" |
| `Category` | `category` | int | 2 |
| `Latitude` | `latitude` | float32 | 45.5 |
| `Longitude` | `longitude` | float32 | -93.2 |
| `Description` | `description` | string or null | "Large buck..." or null |
| `CreatedAt` | `created_at` | ISO 8601 | "2026-02-15T13:44:53Z" |
| `UpdatedAt` | `updated_at` | ISO 8601 | "2026-02-15T13:44:53Z" |

**Special handling**:
- `UserId` not included in API responses (internal only)
- `Description` omitted from response if null (use `omitempty` struct tag)
- Timestamps formatted as RFC 3339 (Go stdlib default for `time.Time.MarshalJSON()`)

---

## Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                        API Requests                          │
│                   (HTTP with X-User-Id)                      │
└────────────────────────────┬────────────────────────────────┘
                             │
                    ┌────────▼─────────┐
                    │    Middleware    │
                    │ Extract UserID   │
                    │ from X-User-Id   │
                    └────────┬─────────┘
                             │
      ┌──────────────────────▼──────────────────────┐
      │          HTTP Handlers                       │
      │  (POST/GET/PATCH/DELETE /api/v1/markers)   │
      └──────────────────────┬──────────────────────┘
                             │
      ┌──────────────────────▼──────────────────────┐
      │         Service Layer                        │
      │  (ID generation, filtering, validation)     │
      └──────────────────────┬──────────────────────┘
                             │
      ┌──────────────────────▼──────────────────────┐
      │      Repository (In-Memory)                 │
      │   map[UserID]map[MarkerID]Marker            │
      │   Protected by RWMutex                      │
      └──────────────────────┬──────────────────────┘
                             │
             ┌───────────────┴───────────────┐
             │                               │
        ┌────▼────────┐          ┌───────────▼─────┐
        │  UserID 1   │          │   UserID 2      │
        │  ┌────────┐ │          │   ┌──────────┐  │
        │  │Marker1 │ │          │   │ Marker42 │  │
        │  │Marker2 │ │          │   │ Marker43 │  │
        │  │Marker3 │ │          │   │ Marker44 │  │
        │  └────────┘ │          │   └──────────┘  │
        └─────────────┘          └─────────────────┘
```

---

## Constraints Summary

### Hard Constraints (Must Enforce)

1. MarkerID format MUST match regex: `^[a-zA-Z0-9_]+_v[0-9]+\.[0-9]+_[0-9]+$`
2. Latitude MUST be in range [-90.0, 90.0]
3. Longitude MUST be in range [-180.0, 180.0]
4. Category MUST be one of four enumerated values
5. UserId MUST be positive integer
6. Name MUST be non-empty
7. CreatedAt and UpdatedAt MUST use UTC timezone

### Soft Constraints (Should Enforce)

1. Description SHOULD be <= 5000 characters (truncate or reject per implementation)
2. Name SHOULD be <= 1000 characters (reasonable limit for human-readable field)
3. MarkerID generation SHOULD be deterministic (same inputs → same ID)

