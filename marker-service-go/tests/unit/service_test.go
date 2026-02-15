package unit

import (
	"context"
	"testing"
	"time"

	"github.com/marker-service-go/models"
	"github.com/marker-service-go/repository"
	"github.com/marker-service-go/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateMarkerID(t *testing.T) {
	repo := repository.NewInMemoryMarkerRepository()
	svc := service.NewMarkerService(repo)

	timestamp := time.Unix(1707999893, 0)

	id := svc.(*service.MarkerService).GenerateMarkerID(models.CategoryGeneral, timestamp)
	assert.Equal(t, "General_v1.0_1707999893", id)

	id2 := svc.(*service.MarkerService).GenerateMarkerID(models.CategoryGroupCategory1Marker1, timestamp)
	assert.Equal(t, "GroupCategory1_Marker1_v1.0_1707999893", id2)
}

func TestCreateMarkerAsync(t *testing.T) {
	repo := repository.NewInMemoryMarkerRepository()
	svc := service.NewMarkerService(repo)
	ctx := context.Background()

	request := &models.CreateMarkerRequest{
		Name:        "Test Marker",
		Category:    models.CategoryGeneral,
		Latitude:    45.5,
		Longitude:   -93.2,
		Description: strPtr("A test marker"),
	}

	response, err := svc.CreateMarkerAsync(ctx, 1, request)
	require.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "Test Marker", response.Name)
	assert.Equal(t, models.CategoryGeneral, response.Category)
	assert.Equal(t, float32(45.5), response.Latitude)
	assert.Equal(t, float32(-93.2), response.Longitude)
	assert.NotEmpty(t, response.ID)
}

func TestCreateMarkerAsync_Validation(t *testing.T) {
	repo := repository.NewInMemoryMarkerRepository()
	svc := service.NewMarkerService(repo)
	ctx := context.Background()

	// Missing name
	request := &models.CreateMarkerRequest{
		Name:     "",
		Category: models.CategoryGeneral,
		Latitude: 45.5,
		Longitude: -93.2,
	}

	_, err := svc.CreateMarkerAsync(ctx, 1, request)
	assert.Error(t, err)

	// Invalid latitude
	request.Name = "Test"
	request.Latitude = 95.0
	_, err = svc.CreateMarkerAsync(ctx, 1, request)
	assert.Error(t, err)

	// Invalid longitude
	request.Latitude = 45.5
	request.Longitude = 200.0
	_, err = svc.CreateMarkerAsync(ctx, 1, request)
	assert.Error(t, err)
}

func TestGetMarkerAsync(t *testing.T) {
	repo := repository.NewInMemoryMarkerRepository()
	svc := service.NewMarkerService(repo)
	ctx := context.Background()

	request := &models.CreateMarkerRequest{
		Name:     "Test Marker",
		Category: models.CategoryGeneral,
		Latitude: 45.5,
		Longitude: -93.2,
	}

	created, _ := svc.CreateMarkerAsync(ctx, 1, request)

	retrieved, err := svc.GetMarkerAsync(ctx, 1, created.ID)
	require.NoError(t, err)
	assert.NotNil(t, retrieved)
	assert.Equal(t, created.ID, retrieved.ID)
	assert.Equal(t, "Test Marker", retrieved.Name)
}

func TestGetAllMarkersAsync(t *testing.T) {
	repo := repository.NewInMemoryMarkerRepository()
	svc := service.NewMarkerService(repo)
	ctx := context.Background()

	// Create multiple markers
	req1 := &models.CreateMarkerRequest{
		Name:     "Marker 1",
		Category: models.CategoryGeneral,
		Latitude: 45.5,
		Longitude: -93.2,
	}

	req2 := &models.CreateMarkerRequest{
		Name:     "Marker 2",
		Category: models.CategoryGroupCategory1Marker1,
		Latitude: 46.0,
		Longitude: -93.0,
	}

	svc.CreateMarkerAsync(ctx, 1, req1)
	svc.CreateMarkerAsync(ctx, 1, req2)

	// Get all markers
	markers, err := svc.GetAllMarkersAsync(ctx, 1, nil)
	require.NoError(t, err)
	assert.Len(t, markers, 2)

	// Get markers by category
	category := models.CategoryGeneral
	categoryMarkers, err := svc.GetAllMarkersAsync(ctx, 1, &category)
	require.NoError(t, err)
	assert.Len(t, categoryMarkers, 1)
	assert.Equal(t, "Marker 1", categoryMarkers[0].Name)
}

func TestUpdateMarkerAsync(t *testing.T) {
	repo := repository.NewInMemoryMarkerRepository()
	svc := service.NewMarkerService(repo)
	ctx := context.Background()

	request := &models.CreateMarkerRequest{
		Name:     "Original Name",
		Category: models.CategoryGeneral,
		Latitude: 45.5,
		Longitude: -93.2,
	}

	created, _ := svc.CreateMarkerAsync(ctx, 1, request)

	updateRequest := &models.UpdateMarkerRequest{
		Name: strPtr("Updated Name"),
		Latitude: float32Ptr(46.0),
	}

	updated, err := svc.UpdateMarkerAsync(ctx, 1, created.ID, updateRequest)
	require.NoError(t, err)
	assert.NotNil(t, updated)
	assert.Equal(t, "Updated Name", updated.Name)
	assert.Equal(t, float32(46.0), updated.Latitude)
	assert.Equal(t, float32(-93.2), updated.Longitude) // Unchanged
}

func TestDeleteMarkerAsync(t *testing.T) {
	repo := repository.NewInMemoryMarkerRepository()
	svc := service.NewMarkerService(repo)
	ctx := context.Background()

	request := &models.CreateMarkerRequest{
		Name:     "Test Marker",
		Category: models.CategoryGeneral,
		Latitude: 45.5,
		Longitude: -93.2,
	}

	created, _ := svc.CreateMarkerAsync(ctx, 1, request)

	deleted, err := svc.DeleteMarkerAsync(ctx, 1, created.ID)
	require.NoError(t, err)
	assert.True(t, deleted)

	// Verify it's gone
	retrieved, _ := svc.GetMarkerAsync(ctx, 1, created.ID)
	assert.Nil(t, retrieved)
}

// Helper functions
func strPtr(s string) *string {
	return &s
}

func float32Ptr(f float32) *float32 {
	return &f
}
