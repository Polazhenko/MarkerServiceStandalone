package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/marker-service-go/handlers"
	"github.com/marker-service-go/models"
	"github.com/marker-service-go/repository"
	"github.com/marker-service-go/service"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	router := gin.Default()

	repo := repository.NewInMemoryMarkerRepository()
	svc := service.NewMarkerService(repo)
	markerHandler := handlers.NewMarkerHandler(svc)
	swaggerHandler := handlers.NewSwaggerHandler()

	router.Use(handlers.ExtractUserID())

	router.GET("/", swaggerHandler.HealthCheck)

	api := router.Group("/api/v1")
	{
		api.POST("/markers", markerHandler.CreateMarker)
		api.GET("/markers", markerHandler.GetAllMarkers)
		api.GET("/markers/:markerId", markerHandler.GetMarker)
		api.PATCH("/markers/:markerId", markerHandler.UpdateMarker)
		api.DELETE("/markers/:markerId", markerHandler.DeleteMarker)
	}

	return router
}

func TestCreateMarkerIntegration(t *testing.T) {
	router := setupRouter()

	body := models.CreateMarkerRequest{
		Name:        "Deer Sighting",
		Category:    models.CategoryGroupCategory1Marker1,
		Latitude:    45.5,
		Longitude:   -93.2,
		Description: strPtr("Large buck spotted"),
	}

	bodyBytes, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/api/v1/markers", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-Id", "1")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.MarkerResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Deer Sighting", response.Name)
	assert.NotEmpty(t, response.ID)
}

func TestGetMarkerIntegration(t *testing.T) {
	router := setupRouter()

	// Create a marker first
	createBody := models.CreateMarkerRequest{
		Name:      "Test Marker",
		Category:  models.CategoryGeneral,
		Latitude:  45.5,
		Longitude: -93.2,
	}

	bodyBytes, _ := json.Marshal(createBody)
	createReq := httptest.NewRequest("POST", "/api/v1/markers", bytes.NewReader(bodyBytes))
	createReq.Header.Set("Content-Type", "application/json")
	createReq.Header.Set("X-User-Id", "1")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, createReq)

	var createdMarker models.MarkerResponse
	json.Unmarshal(w.Body.Bytes(), &createdMarker)

	// Now get the marker
	getReq := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/markers/%s", createdMarker.ID), nil)
	getReq.Header.Set("X-User-Id", "1")

	w = httptest.NewRecorder()
	router.ServeHTTP(w, getReq)

	assert.Equal(t, http.StatusOK, w.Code)

	var retrievedMarker models.MarkerResponse
	json.Unmarshal(w.Body.Bytes(), &retrievedMarker)
	assert.Equal(t, createdMarker.ID, retrievedMarker.ID)
	assert.Equal(t, "Test Marker", retrievedMarker.Name)
}

func TestGetAllMarkersIntegration(t *testing.T) {
	router := setupRouter()

	// Create multiple markers
	for i := 1; i <= 3; i++ {
		body := models.CreateMarkerRequest{
			Name:      fmt.Sprintf("Marker %d", i),
			Category:  models.CategoryGeneral,
			Latitude:  45.5,
			Longitude: -93.2,
		}

		bodyBytes, _ := json.Marshal(body)
		req := httptest.NewRequest("POST", "/api/v1/markers", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-User-Id", "1")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}

	// Get all markers
	getReq := httptest.NewRequest("GET", "/api/v1/markers", nil)
	getReq.Header.Set("X-User-Id", "1")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, getReq)

	assert.Equal(t, http.StatusOK, w.Code)

	var markers []models.MarkerResponse
	json.Unmarshal(w.Body.Bytes(), &markers)
	assert.Len(t, markers, 3)
}

func TestUpdateMarkerIntegration(t *testing.T) {
	router := setupRouter()

	// Create a marker
	createBody := models.CreateMarkerRequest{
		Name:      "Original Name",
		Category:  models.CategoryGeneral,
		Latitude:  45.5,
		Longitude: -93.2,
	}

	bodyBytes, _ := json.Marshal(createBody)
	createReq := httptest.NewRequest("POST", "/api/v1/markers", bytes.NewReader(bodyBytes))
	createReq.Header.Set("Content-Type", "application/json")
	createReq.Header.Set("X-User-Id", "1")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, createReq)

	var createdMarker models.MarkerResponse
	json.Unmarshal(w.Body.Bytes(), &createdMarker)

	// Update the marker
	updateBody := models.UpdateMarkerRequest{
		Name: strPtr("Updated Name"),
	}

	updateBytes, _ := json.Marshal(updateBody)
	updateReq := httptest.NewRequest("PATCH", fmt.Sprintf("/api/v1/markers/%s", createdMarker.ID), bytes.NewReader(updateBytes))
	updateReq.Header.Set("Content-Type", "application/json")
	updateReq.Header.Set("X-User-Id", "1")

	w = httptest.NewRecorder()
	router.ServeHTTP(w, updateReq)

	assert.Equal(t, http.StatusOK, w.Code)

	var updatedMarker models.MarkerResponse
	json.Unmarshal(w.Body.Bytes(), &updatedMarker)
	assert.Equal(t, "Updated Name", updatedMarker.Name)
}

func TestDeleteMarkerIntegration(t *testing.T) {
	router := setupRouter()

	// Create a marker
	createBody := models.CreateMarkerRequest{
		Name:      "Test Marker",
		Category:  models.CategoryGeneral,
		Latitude:  45.5,
		Longitude: -93.2,
	}

	bodyBytes, _ := json.Marshal(createBody)
	createReq := httptest.NewRequest("POST", "/api/v1/markers", bytes.NewReader(bodyBytes))
	createReq.Header.Set("Content-Type", "application/json")
	createReq.Header.Set("X-User-Id", "1")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, createReq)

	var createdMarker models.MarkerResponse
	json.Unmarshal(w.Body.Bytes(), &createdMarker)

	// Delete the marker
	deleteReq := httptest.NewRequest("DELETE", fmt.Sprintf("/api/v1/markers/%s", createdMarker.ID), nil)
	deleteReq.Header.Set("X-User-Id", "1")

	w = httptest.NewRecorder()
	router.ServeHTTP(w, deleteReq)

	assert.Equal(t, http.StatusNoContent, w.Code)

	// Verify it's gone
	getReq := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/markers/%s", createdMarker.ID), nil)
	getReq.Header.Set("X-User-Id", "1")

	w = httptest.NewRecorder()
	router.ServeHTTP(w, getReq)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUserIsolationIntegration(t *testing.T) {
	router := setupRouter()

	// Create marker for user 1
	body := models.CreateMarkerRequest{
		Name:      "User 1 Marker",
		Category:  models.CategoryGeneral,
		Latitude:  45.5,
		Longitude: -93.2,
	}

	bodyBytes, _ := json.Marshal(body)
	createReq := httptest.NewRequest("POST", "/api/v1/markers", bytes.NewReader(bodyBytes))
	createReq.Header.Set("Content-Type", "application/json")
	createReq.Header.Set("X-User-Id", "1")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, createReq)

	var createdMarker models.MarkerResponse
	json.Unmarshal(w.Body.Bytes(), &createdMarker)

	// Try to get marker as user 2
	getReq := httptest.NewRequest("GET", fmt.Sprintf("/api/v1/markers/%s", createdMarker.ID), nil)
	getReq.Header.Set("X-User-Id", "2")

	w = httptest.NewRecorder()
	router.ServeHTTP(w, getReq)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// Helper function
func strPtr(s string) *string {
	return &s
}
