package models

import (
	"encoding/json"
	"fmt"
	"time"
)

// MarkerCategory enum for marker types
type MarkerCategory int

const (
	CategoryGeneral MarkerCategory = iota + 1
	CategoryGroupCategory1Marker1
	CategoryGroupCategory1Marker2
	CategoryGroupCategory2Marker1
)

// String returns the string representation of MarkerCategory
func (mc MarkerCategory) String() string {
	switch mc {
	case CategoryGeneral:
		return "General"
	case CategoryGroupCategory1Marker1:
		return "GroupCategory1_Marker1"
	case CategoryGroupCategory1Marker2:
		return "GroupCategory1_Marker2"
	case CategoryGroupCategory2Marker1:
		return "GroupCategory2_Marker1"
	default:
		return "Unknown"
	}
}

// MarkerCategoryFromString converts string to MarkerCategory
func MarkerCategoryFromString(s string) (MarkerCategory, error) {
	switch s {
	case "General":
		return CategoryGeneral, nil
	case "GroupCategory1_Marker1":
		return CategoryGroupCategory1Marker1, nil
	case "GroupCategory1_Marker2":
		return CategoryGroupCategory1Marker2, nil
	case "GroupCategory2_Marker1":
		return CategoryGroupCategory2Marker1, nil
	default:
		return 0, fmt.Errorf("invalid marker category: %s", s)
	}
}

// IsValidCategory checks if category is valid
func IsValidCategory(category MarkerCategory) bool {
	return category >= CategoryGeneral && category <= CategoryGroupCategory2Marker1
}

// Marker represents a geotagged marker
type Marker struct {
	ID          string    `json:"id"`
	UserID      int64     `json:"user_id,omitempty"`
	Name        string    `json:"name"`
	Category    MarkerCategory `json:"category"`
	Latitude    float32   `json:"latitude"`
	Longitude   float32   `json:"longitude"`
	Description *string   `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateMarkerRequest is the request body for creating a marker
type CreateMarkerRequest struct {
	Name        string    `json:"name" binding:"required"`
	Category    MarkerCategory `json:"category" binding:"required"`
	Latitude    float32   `json:"latitude" binding:"required"`
	Longitude   float32   `json:"longitude" binding:"required"`
	Description *string   `json:"description"`
}

// UpdateMarkerRequest is the request body for updating a marker
type UpdateMarkerRequest struct {
	Name        *string   `json:"name"`
	Latitude    *float32  `json:"latitude"`
	Longitude   *float32  `json:"longitude"`
	Description *string   `json:"description"`
}

// MarkerResponse is the response body for marker endpoints
type MarkerResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Category    MarkerCategory `json:"category"`
	Latitude    float32   `json:"latitude"`
	Longitude   float32   `json:"longitude"`
	Description *string   `json:"description"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at"`
}

// ToMarkerResponse converts Marker to MarkerResponse
func (m *Marker) ToMarkerResponse() *MarkerResponse {
	return &MarkerResponse{
		ID:          m.ID,
		Name:        m.Name,
		Category:    m.Category,
		Latitude:    m.Latitude,
		Longitude:   m.Longitude,
		Description: m.Description,
		CreatedAt:   m.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   m.UpdatedAt.Format(time.RFC3339),
	}
}

// MarshalJSON implements custom JSON marshaling for MarkerResponse
func (mr *MarkerResponse) MarshalJSON() ([]byte, error) {
	type Alias MarkerResponse
	return json.Marshal(&struct {
		Category int `json:"category"`
		*Alias
	}{
		Category: int(mr.Category),
		Alias:    (*Alias)(mr),
	})
}

// UnmarshalJSON implements custom JSON unmarshaling for MarkerCategory
func (mc *MarkerCategory) UnmarshalJSON(data []byte) error {
	var v int
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*mc = MarkerCategory(v)
	return nil
}

// MarshalJSON implements custom JSON marshaling for MarkerCategory
func (mc MarkerCategory) MarshalJSON() ([]byte, error) {
	return json.Marshal(int(mc))
}
