package api

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	db "github.com/cyriljohn147/carpooling/db/sqlc"
)

// TripHandler handles trip-related HTTP requests
type TripHandler struct {
	tripService *TripService
}

// NewTripHandler creates a new trip handler
func NewTripHandler(store *db.Queries, database *sql.DB) *TripHandler {
	return &TripHandler{
		tripService: NewTripService(store, database),
	}
}

// CreateTrip handles POST /trips
func (h *TripHandler) CreateTrip(c *gin.Context) {
	var req TripRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request format",
			Validations: []ValidationError{
				{
					Field:   "request",
					Message: err.Error(),
				},
			},
		})
		return
	}

	// Validate required fields
	if err := validateTripRequest(req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:       "Validation failed",
			Validations: err,
		})
		return
	}

	// Create trip
	response, err := h.tripService.CreateTrip(c.Request.Context(), req)
	if err != nil {
		// Check if it's a user not found error
		if containsString(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to create trip: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetTrip handles GET /trips/:id
func (h *TripHandler) GetTrip(c *gin.Context) {
	tripIDStr := c.Param("id")
	tripID, err := strconv.ParseInt(tripIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid trip ID",
		})
		return
	}

	response, err := h.tripService.GetTrip(c.Request.Context(), int32(tripID))
	if err != nil {
		if containsString(err.Error(), "no rows") {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "Trip not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to get trip: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// validateTripRequest validates the simple trip request
func validateTripRequest(req TripRequest) []ValidationError {
	var errors []ValidationError

	if req.DriverName == "" {
		errors = append(errors, ValidationError{
			Field:   "driver_name",
			Message: "Driver name is required",
		})
	}

	if req.TotalCost <= 0 {
		errors = append(errors, ValidationError{
			Field:   "total_cost",
			Message: "Total cost must be greater than 0",
		})
	}

	if len(req.Participants) == 0 {
		errors = append(errors, ValidationError{
			Field:   "participants",
			Message: "At least one participant is required",
		})
	}

	// Check if driver is in participants list
	driverInParticipants := false
	for _, participant := range req.Participants {
		if participant == req.DriverName {
			driverInParticipants = true
			break
		}
	}

	if !driverInParticipants {
		errors = append(errors, ValidationError{
			Field:   "participants",
			Message: "Driver must be included in participants list",
		})
	}

	// Check for duplicate participants
	participantSet := make(map[string]bool)
	for _, participant := range req.Participants {
		if participant == "" {
			errors = append(errors, ValidationError{
				Field:   "participants",
				Message: "Participant names cannot be empty",
			})
			continue
		}

		if participantSet[participant] {
			errors = append(errors, ValidationError{
				Field:   "participants",
				Message: "Duplicate participant: " + participant,
			})
		}
		participantSet[participant] = true
	}

	return errors
}

// containsString checks if a string contains a substring
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) &&
			(s[:len(substr)] == substr ||
				s[len(s)-len(substr):] == substr ||
				containsSubstring(s, substr))))
}

// containsSubstring is a helper function for string matching
func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// RegisterTripRoutes registers all trip-related routes
func RegisterTripRoutes(router *gin.Engine, store *db.Queries, database *sql.DB) {
	handler := NewTripHandler(store, database)

	trips := router.Group("/trips")
	{
		// Trip endpoints
		trips.POST("/", handler.CreateTrip)
		trips.GET("/:id", handler.GetTrip)
	}
}
