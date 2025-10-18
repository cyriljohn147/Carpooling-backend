package api

import (
	"database/sql"
)

// TripRequest represents the incoming JSON for creating a basic trip
type TripRequest struct {
	DriverName   string    `json:"driver_name" binding:"required"`
	TotalCost    float64   `json:"total_cost" binding:"required"`
	Participants []string  `json:"participants" binding:"required,min=1"`
}


// TripResponse represents the response after creating a basic trip
type TripResponse struct {
	TripID       int32             `json:"trip_id"`
	DriverID     int32             `json:"driver_id"`
	DriverName   string            `json:"driver_name"`
	Date         string            `json:"date"`
	Time         string            `json:"time"`
	TotalCost    float64           `json:"total_cost"`
	Participants []ParticipantInfo `json:"participants"`
	Message      string            `json:"message"`
	CreatedAt    string            `json:"created_at"`
}


// ParticipantInfo represents participant details in the response
type ParticipantInfo struct {
	UserID      int32   `json:"user_id"`
	Name        string  `json:"name"`
	ShareAmount float64 `json:"share_amount"`
}


// ValidationError represents validation errors for the request
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ErrorResponse represents error responses
type ErrorResponse struct {
	Error       string            `json:"error"`
	Validations []ValidationError `json:"validations,omitempty"`
}

// Helper function to convert SQL NullString to string
func nullStringToString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

// Helper function to convert SQL NullTime to string
func nullTimeToString(nt sql.NullTime) string {
	if nt.Valid {
		return nt.Time.Format("2006-01-02")
	}
	return ""
}

// Helper function to convert SQL NullTime to time string
func nullTimeToTimeString(nt sql.NullTime) string {
	if nt.Valid {
		return nt.Time.Format("15:04:05")
	}
	return ""
}
