package api

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	db "github.com/cyriljohn147/carpooling/db/sqlc"
)

// BalanceHandler handles balance-related HTTP requests
type BalanceHandler struct {
	store *db.Queries
	db    *sql.DB
}

// NewBalanceHandler creates a new balance handler
func NewBalanceHandler(store *db.Queries, database *sql.DB) *BalanceHandler {
	return &BalanceHandler{
		store: store,
		db:    database,
	}
}

// UserBalance represents a user's current balance
type UserBalance struct {
	UserID    int32   `json:"user_id"`
	Name      string  `json:"name,omitempty"`
	TotalOwed float64 `json:"total_owed"`
}

// PaymentRequest represents a payment request
type PaymentRequest struct {
	UserID int32   `json:"user_id" binding:"required"`
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

// PaymentResponse represents a payment response
type PaymentResponse struct {
	UserID        int32   `json:"user_id"`
	Name          string  `json:"name"`
	AmountPaid    float64 `json:"amount_paid"`
	RemainingOwed float64 `json:"remaining_owed"`
	Message       string  `json:"message"`
}

// GetAllUsersBalances handles GET /balances
func (h *BalanceHandler) GetAllUsersBalances(c *gin.Context) {
	balances, err := h.store.GetAllUsersCurrentBalances(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to get user balances: " + err.Error(),
		})
		return
	}

	var userBalances []UserBalance
	for _, balance := range balances {
		totalOwed, _ := strconv.ParseFloat(balance.TotalOwed, 64)
		userBalances = append(userBalances, UserBalance{
			UserID:    balance.UserID,
			Name:      nullStringToString(balance.Name),
			TotalOwed: totalOwed,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"balances": userBalances,
		"message":  "User balances retrieved successfully",
	})
}

// GetUserBalance handles GET /users/:id/balance
func (h *BalanceHandler) GetUserBalance(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid user ID",
		})
		return
	}

	// First get user details
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

	// Get user's current balance
	balance, err := h.store.GetCurrentUserBalance(c.Request.Context(), int32(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to get user balance: " + err.Error(),
		})
		return
	}

	totalOwed, _ := strconv.ParseFloat(balance.TotalOwed, 64)
	c.JSON(http.StatusOK, UserBalance{
		UserID:    balance.UserID,
		Name:      nullStringToString(user.Name),
		TotalOwed: totalOwed,
	})
}

// PayOffUserBalance handles POST /users/:id/pay-off
func (h *BalanceHandler) PayOffUserBalance(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid user ID",
		})
		return
	}

	// Get user details
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

	// Get current balance before paying off
	currentBalance, err := h.store.GetCurrentUserBalance(c.Request.Context(), int32(userID))
	var totalOwedBefore float64 = 0
	if err == nil {
		totalOwedBefore, _ = strconv.ParseFloat(currentBalance.TotalOwed, 64)
	}

	// Pay off all balances
	err = h.store.PayOffUserBalance(c.Request.Context(), int32(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to pay off balance: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, PaymentResponse{
		UserID:        int32(userID),
		Name:          nullStringToString(user.Name),
		AmountPaid:    totalOwedBefore,
		RemainingOwed: 0.0,
		Message:       "All balances paid off successfully",
	})
}

// MakePartialPayment handles POST /users/:id/pay
func (h *BalanceHandler) MakePartialPayment(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid user ID",
		})
		return
	}

	var req PaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request format: " + err.Error(),
		})
		return
	}

	// Validate that the user ID in URL matches the request body
	if req.UserID != int32(userID) {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "User ID in URL must match user ID in request body",
		})
		return
	}

	// Get user details
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

	// Get current balance before payment
	currentBalance, err := h.store.GetCurrentUserBalance(c.Request.Context(), int32(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to get current balance: " + err.Error(),
		})
		return
	}

	totalOwedBefore, _ := strconv.ParseFloat(currentBalance.TotalOwed, 64)

	// Check if user has outstanding balance
	if totalOwedBefore <= 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "User has no outstanding balance to pay",
		})
		return
	}

	// Make partial payment
	err = h.store.PayPartialUserBalance(c.Request.Context(), db.PayPartialUserBalanceParams{
		ID:      int32(userID),
		Balance: sql.NullString{String: strconv.FormatFloat(req.Amount, 'f', 2, 64), Valid: true},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to process payment: " + err.Error(),
		})
		return
	}

	// Get updated balance
	updatedBalance, err := h.store.GetCurrentUserBalance(c.Request.Context(), int32(userID))
	var remainingOwed float64 = 0
	if err == nil {
		remainingOwed, _ = strconv.ParseFloat(updatedBalance.TotalOwed, 64)
	}

	actualPayment := totalOwedBefore - remainingOwed
	message := "Payment processed successfully"
	if remainingOwed == 0 {
		message = "All balances paid off successfully"
	}

	c.JSON(http.StatusOK, PaymentResponse{
		UserID:        int32(userID),
		Name:          nullStringToString(user.Name),
		AmountPaid:    actualPayment,
		RemainingOwed: remainingOwed,
		Message:       message,
	})
}

// RegisterBalanceRoutes registers all balance-related routes
func RegisterBalanceRoutes(router *gin.Engine, store *db.Queries, database *sql.DB) {
	handler := NewBalanceHandler(store, database)

	// Balance endpoints
	router.GET("/balances", handler.GetAllUsersBalances)

	// User balance endpoints - register directly on router to avoid conflicts
	router.GET("/users/:id/balance", handler.GetUserBalance)
	router.POST("/users/:id/pay-off", handler.PayOffUserBalance)
	router.POST("/users/:id/pay", handler.MakePartialPayment)
}
