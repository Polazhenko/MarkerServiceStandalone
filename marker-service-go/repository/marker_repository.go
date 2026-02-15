package repository

import (
	"sync"

	"github.com/marker-service-go/models"
)

// IMarkerRepository defines the interface for marker storage operations
type IMarkerRepository interface {
	CreateAsync(marker *models.Marker) (*models.Marker, error)
	GetByIdAsync(id string, userId int64) (*models.Marker, error)
	GetByCategoryAsync(userId int64, category *models.MarkerCategory) ([]*models.Marker, error)
	UpdateAsync(marker *models.Marker) (*models.Marker, error)
	DeleteAsync(id string, userId int64) (bool, error)
}

// InMemoryMarkerRepository implements IMarkerRepository with in-memory storage
type InMemoryMarkerRepository struct {
	mu      sync.RWMutex
	markers map[int64]map[string]*models.Marker
}

// NewInMemoryMarkerRepository creates a new in-memory repository
func NewInMemoryMarkerRepository() IMarkerRepository {
	return &InMemoryMarkerRepository{
		markers: make(map[int64]map[string]*models.Marker),
	}
}

// CreateAsync stores a new marker
func (r *InMemoryMarkerRepository) CreateAsync(marker *models.Marker) (*models.Marker, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Create user entry if needed
	if _, exists := r.markers[marker.UserID]; !exists {
		r.markers[marker.UserID] = make(map[string]*models.Marker)
	}

	// Store marker with ID as key
	r.markers[marker.UserID][marker.ID] = marker

	return marker, nil
}

// GetByIdAsync retrieves a marker by ID for a specific user
func (r *InMemoryMarkerRepository) GetByIdAsync(id string, userId int64) (*models.Marker, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	userMarkers, exists := r.markers[userId]
	if !exists {
		return nil, nil
	}

	marker, found := userMarkers[id]
	if !found {
		return nil, nil
	}

	return marker, nil
}

// GetByCategoryAsync retrieves markers for a user, optionally filtered by category
func (r *InMemoryMarkerRepository) GetByCategoryAsync(userId int64, category *models.MarkerCategory) ([]*models.Marker, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	userMarkers, exists := r.markers[userId]
	if !exists {
		return []*models.Marker{}, nil
	}

	var result []*models.Marker

	if category == nil {
		// Return all markers for user
		for _, marker := range userMarkers {
			result = append(result, marker)
		}
	} else {
		// Return only markers matching category
		for _, marker := range userMarkers {
			if marker.Category == *category {
				result = append(result, marker)
			}
		}
	}

	return result, nil
}

// UpdateAsync updates an existing marker
func (r *InMemoryMarkerRepository) UpdateAsync(marker *models.Marker) (*models.Marker, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	userMarkers, exists := r.markers[marker.UserID]
	if !exists {
		return nil, nil
	}

	existing, found := userMarkers[marker.ID]
	if !found {
		return nil, nil
	}

	// Update fields
	existing.Name = marker.Name
	existing.Latitude = marker.Latitude
	existing.Longitude = marker.Longitude
	existing.Description = marker.Description
	existing.UpdatedAt = marker.UpdatedAt

	return existing, nil
}

// DeleteAsync removes a marker
func (r *InMemoryMarkerRepository) DeleteAsync(id string, userId int64) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	userMarkers, exists := r.markers[userId]
	if !exists {
		return false, nil
	}

	if _, found := userMarkers[id]; !found {
		return false, nil
	}

	delete(userMarkers, id)
	return true, nil
}
