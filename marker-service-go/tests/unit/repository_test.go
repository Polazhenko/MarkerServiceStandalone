package unit

import (
	"testing"
	"time"

	"github.com/marker-service-go/models"
	"github.com/marker-service-go/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateMarker(t *testing.T) {
	repo := repository.NewInMemoryMarkerRepository()

	marker := &models.Marker{
		ID:        "General_v1.0_1234567890",
		UserID:    1,
		Name:      "Test Marker",
		Category:  models.CategoryGeneral,
		Latitude:  45.5,
		Longitude: -93.2,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	created, err := repo.CreateAsync(marker)
	require.NoError(t, err)
	assert.Equal(t, marker.ID, created.ID)
	assert.Equal(t, marker.Name, created.Name)
}

func TestGetByIdAsync(t *testing.T) {
	repo := repository.NewInMemoryMarkerRepository()

	marker := &models.Marker{
		ID:        "General_v1.0_1234567890",
		UserID:    1,
		Name:      "Test Marker",
		Category:  models.CategoryGeneral,
		Latitude:  45.5,
		Longitude: -93.2,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	repo.CreateAsync(marker)

	retrieved, err := repo.GetByIdAsync(marker.ID, 1)
	require.NoError(t, err)
	assert.NotNil(t, retrieved)
	assert.Equal(t, marker.ID, retrieved.ID)
}

func TestGetByIdAsync_NotFound(t *testing.T) {
	repo := repository.NewInMemoryMarkerRepository()

	retrieved, err := repo.GetByIdAsync("nonexistent", 1)
	require.NoError(t, err)
	assert.Nil(t, retrieved)
}

func TestUserIsolation(t *testing.T) {
	repo := repository.NewInMemoryMarkerRepository()

	marker1 := &models.Marker{
		ID:       "General_v1.0_1234567890",
		UserID:   1,
		Name:     "User 1 Marker",
		Category: models.CategoryGeneral,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	marker2 := &models.Marker{
		ID:       "General_v1.0_1234567891",
		UserID:   2,
		Name:     "User 2 Marker",
		Category: models.CategoryGeneral,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	repo.CreateAsync(marker1)
	repo.CreateAsync(marker2)

	// User 1 should not be able to see user 2's marker
	retrieved, err := repo.GetByIdAsync(marker2.ID, 1)
	require.NoError(t, err)
	assert.Nil(t, retrieved)

	// User 2 should be able to see their own marker
	retrieved, err = repo.GetByIdAsync(marker2.ID, 2)
	require.NoError(t, err)
	assert.NotNil(t, retrieved)
}

func TestGetByCategoryAsync(t *testing.T) {
	repo := repository.NewInMemoryMarkerRepository()

	marker1 := &models.Marker{
		ID:       "General_v1.0_1234567890",
		UserID:   1,
		Name:     "General Marker",
		Category: models.CategoryGeneral,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	marker2 := &models.Marker{
		ID:       "GroupCategory1_Marker1_v1.0_1234567891",
		UserID:   1,
		Name:     "Category 1 Marker",
		Category: models.CategoryGroupCategory1Marker1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	repo.CreateAsync(marker1)
	repo.CreateAsync(marker2)

	// Get all markers
	allMarkers, err := repo.GetByCategoryAsync(1, nil)
	require.NoError(t, err)
	assert.Len(t, allMarkers, 2)

	// Get markers by category
	category := models.CategoryGeneral
	categoryMarkers, err := repo.GetByCategoryAsync(1, &category)
	require.NoError(t, err)
	assert.Len(t, categoryMarkers, 1)
	assert.Equal(t, marker1.ID, categoryMarkers[0].ID)
}

func TestUpdateAsync(t *testing.T) {
	repo := repository.NewInMemoryMarkerRepository()

	marker := &models.Marker{
		ID:        "General_v1.0_1234567890",
		UserID:    1,
		Name:      "Original Name",
		Category:  models.CategoryGeneral,
		Latitude:  45.5,
		Longitude: -93.2,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	repo.CreateAsync(marker)

	// Update marker
	marker.Name = "Updated Name"
	marker.Latitude = 46.0
	updated, err := repo.UpdateAsync(marker)

	require.NoError(t, err)
	assert.NotNil(t, updated)
	assert.Equal(t, "Updated Name", updated.Name)
	assert.Equal(t, float32(46.0), updated.Latitude)
}

func TestDeleteAsync(t *testing.T) {
	repo := repository.NewInMemoryMarkerRepository()

	marker := &models.Marker{
		ID:        "General_v1.0_1234567890",
		UserID:    1,
		Name:      "Test Marker",
		Category:  models.CategoryGeneral,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	repo.CreateAsync(marker)

	// Delete marker
	deleted, err := repo.DeleteAsync(marker.ID, 1)
	require.NoError(t, err)
	assert.True(t, deleted)

	// Verify it's gone
	retrieved, err := repo.GetByIdAsync(marker.ID, 1)
	require.NoError(t, err)
	assert.Nil(t, retrieved)
}
