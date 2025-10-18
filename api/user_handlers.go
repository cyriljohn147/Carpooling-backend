package api

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	db "github.com/cyriljohn147/carpooling/db/sqlc"
)

// UserRequest represents the request body for creating/updating users
type UserRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Phone string `json:"phone" binding:"required"`
}

// UserResponse represents the response for user operations
type UserResponse struct {
	ID        int32   `json:"id"`
	Name      string  `json:"name"`
	Email     string  `json:"email"`
	Phone     string  `json:"phone"`
	Balance   float64 `json:"balance"`
	CreatedAt string  `json:"created_at"`
}

// UsersListResponse represents the response for listing users
type UsersListResponse struct {
	Users []UserResponse `json:"users"`
	Total int32          `json:"total"`
	Page  int32          `json:"page"`
	Limit int32          `json:"limit"`
}

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	store *db.Queries
}

// NewUserHandler creates a new user handler
func NewUserHandler(store *db.Queries) *UserHandler {
	return &UserHandler{store: store}
}

// CreateUser handles POST /users
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req UserRequest

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
	if validationErrors := validateUserRequest(req); len(validationErrors) > 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:       "Validation failed",
			Validations: validationErrors,
		})
		return
	}

	// Check if user with this email already exists
	_, err := h.store.GetUserByEmail(c.Request.Context(), sql.NullString{String: req.Email, Valid: true})
	if err == nil {
		// User exists
		c.JSON(http.StatusConflict, ErrorResponse{
			Error: "User with this email already exists",
		})
		return
	}

	// Create user
	user, err := h.store.CreateUser(c.Request.Context(), db.CreateUserParams{
		Name:  sql.NullString{String: req.Name, Valid: true},
		Email: sql.NullString{String: req.Email, Valid: true},
		Phone: sql.NullString{String: req.Phone, Valid: true},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to create user: " + err.Error(),
		})
		return
	}

	response := UserResponse{
		ID:        user.ID,
		Name:      nullStringToString(user.Name),
		Email:     nullStringToString(user.Email),
		Phone:     nullStringToString(user.Phone),
		CreatedAt: nullTimeToString(user.CreatedAt),
	}

	c.JSON(http.StatusCreated, response)
}

// GetUser handles GET /users/:id
func (h *UserHandler) GetUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid user ID",
		})
		return
	}

	user, err := h.store.GetUser(c.Request.Context(), int32(userID))
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "User not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to get user: " + err.Error(),
		})
		return
	}

	response := UserResponse{
		ID:        user.ID,
		Name:      nullStringToString(user.Name),
		Email:     nullStringToString(user.Email),
		Phone:     nullStringToString(user.Phone),
		Balance:   nullStringToFloat64(user.Balance),
		CreatedAt: nullTimeToString(user.CreatedAt),
	}

	c.JSON(http.StatusOK, response)
}

func nullStringToFloat64(ns sql.NullString) float64 {
	if !ns.Valid {
		return 0.0
	}
	f, err := strconv.ParseFloat(ns.String, 64)
	if err != nil {
		return 0.0 // fallback if invalid
	}
	return f
}

// ListUsers handles GET /users
func (h *UserHandler) ListUsers(c *gin.Context) {
	// Parse query parameters
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.ParseInt(pageStr, 10, 32)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.ParseInt(limitStr, 10, 32)
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	users, err := h.store.ListUsers(c.Request.Context(), db.ListUsersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to list users: " + err.Error(),
		})
		return
	}

	var userResponses []UserResponse
	for _, user := range users {
		userResponses = append(userResponses, UserResponse{
			ID:        user.ID,
			Name:      nullStringToString(user.Name),
			Email:     nullStringToString(user.Email),
			Phone:     nullStringToString(user.Phone),
			Balance:   nullStringToFloat64(user.Balance),
			CreatedAt: nullTimeToString(user.CreatedAt),
		})
	}

	response := UsersListResponse{
		Users: userResponses,
		Total: int32(len(userResponses)), // In a real app, you'd get the total count separately
		Page:  int32(page),
		Limit: int32(limit),
	}

	c.JSON(http.StatusOK, response)
}

// UpdateUser handles PUT /users/:id
func (h *UserHandler) UpdateUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid user ID",
		})
		return
	}

	var req UserRequest
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
	if validationErrors := validateUserRequest(req); len(validationErrors) > 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:       "Validation failed",
			Validations: validationErrors,
		})
		return
	}

	// Check if user exists
	_, err = h.store.GetUser(c.Request.Context(), int32(userID))
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "User not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to check user existence: " + err.Error(),
		})
		return
	}

	// Update user
	user, err := h.store.UpdateUser(c.Request.Context(), db.UpdateUserParams{
		ID:    int32(userID),
		Name:  sql.NullString{String: req.Name, Valid: true},
		Email: sql.NullString{String: req.Email, Valid: true},
		Phone: sql.NullString{String: req.Phone, Valid: true},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to update user: " + err.Error(),
		})
		return
	}

	response := UserResponse{
		ID:        user.ID,
		Name:      nullStringToString(user.Name),
		Email:     nullStringToString(user.Email),
		Phone:     nullStringToString(user.Phone),
		Balance:   nullStringToFloat64(user.Balance),
		CreatedAt: nullTimeToString(user.CreatedAt),
	}

	c.JSON(http.StatusOK, response)
}

// DeleteUser handles DELETE /users/:id
func (h *UserHandler) DeleteUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid user ID",
		})
		return
	}

	// Check if user exists
	_, err = h.store.GetUser(c.Request.Context(), int32(userID))
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "User not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to check user existence: " + err.Error(),
		})
		return
	}

	// Delete user
	err = h.store.DeleteUser(c.Request.Context(), int32(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to delete user: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
	})
}

// GetUserByName handles GET /users/by-name/:name
func (h *UserHandler) GetUserByName(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Name parameter is required",
		})
		return
	}

	user, err := h.store.GetUserByName(c.Request.Context(), sql.NullString{String: name, Valid: true})
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "User not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to get user: " + err.Error(),
		})
		return
	}

	response := UserResponse{
		ID:        user.ID,
		Name:      nullStringToString(user.Name),
		Email:     nullStringToString(user.Email),
		Phone:     nullStringToString(user.Phone),
		Balance:   nullStringToFloat64(user.Balance),
		CreatedAt: nullTimeToString(user.CreatedAt),
	}

	c.JSON(http.StatusOK, response)
}

// validateUserRequest validates the user request
func validateUserRequest(req UserRequest) []ValidationError {
	var errors []ValidationError

	if req.Name == "" {
		errors = append(errors, ValidationError{
			Field:   "name",
			Message: "Name is required",
		})
	}

	if req.Email == "" {
		errors = append(errors, ValidationError{
			Field:   "email",
			Message: "Email is required",
		})
	}

	if req.Phone == "" {
		errors = append(errors, ValidationError{
			Field:   "phone",
			Message: "Phone is required",
		})
	}

	return errors
}

// RegisterUserRoutes registers all user-related routes
func RegisterUserRoutes(router *gin.Engine, store *db.Queries) {
	handler := NewUserHandler(store)

	users := router.Group("/users")
	{
		users.POST("/", handler.CreateUser)
		users.GET("/", handler.ListUsers)
		users.GET("/:id", handler.GetUser)
		users.PUT("/:id", handler.UpdateUser)
		users.DELETE("/:id", handler.DeleteUser)
		users.GET("/by-name/:name", handler.GetUserByName)
	}
}
