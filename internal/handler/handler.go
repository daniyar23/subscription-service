package handler

import (
	"net/http"
	"time"

	"github.com/daniyar23/subscribe-service/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// SubscriptionHandler handles HTTP requests for subscription operations.
type SubscriptionHandler struct {
	service SubscriptionService
	logger  *zap.Logger
}

// NewSubscriptionHandler creates a new SubscriptionHandler with the given service and logger.
func NewSubscriptionHandler(service SubscriptionService, logger *zap.Logger) *SubscriptionHandler {
	return &SubscriptionHandler{
		service: service,
		logger:  logger,
	}
}

// Create handles the HTTP request to create a new subscription.
func (h *SubscriptionHandler) Create(c *gin.Context) {

	h.logger.Info("create subscription request")

	var sub model.Subscription
	if err := c.ShouldBindJSON(&sub); err != nil {
		h.logger.Error("invalid request body",
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.Create(c.Request.Context(), sub)
	if err != nil {
		h.logger.Error("create subscription: service error",
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	h.logger.Info("subscription created",
		zap.String("id", result.ID.String()),
	)

	c.JSON(http.StatusCreated, result)
}

// GetByID handles the HTTP request to get a subscription by its ID.
func (h *SubscriptionHandler) GetByID(c *gin.Context) {

	idParam := c.Param("id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		h.logger.Error("invalid uuid",
			zap.String("id", idParam),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid uuid",
		})
		return
	}

	h.logger.Info("get subscription by id",
		zap.String("id", id.String()),
	)

	sub, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("getByID service error",
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, sub)
}

// GetAll handles the HTTP request to get all subscriptions.
func (h *SubscriptionHandler) GetAll(c *gin.Context) {

	h.logger.Info("get all subscriptions")

	subs, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		h.logger.Error("getAll service error",
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, subs)
}

// GetByUserID handles the HTTP request to get all subscriptions for a given user ID.
func (h *SubscriptionHandler) GetByUserID(c *gin.Context) {

	userIDParam := c.Param("user_id")

	if userIDParam == "" {
		h.logger.Error("user_id is empty")
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		h.logger.Error("invalid uuid",
			zap.String("user_id", userIDParam),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid uuid"})
		return
	}

	h.logger.Info("get subscriptions by user",
		zap.String("user_id", userID.String()),
	)

	subs, err := h.service.GetByUserID(c.Request.Context(), userID)
	if err != nil {
		h.logger.Error("getByUserID service error",
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, subs)
}

// Update handles the HTTP request to update a subscription.
func (h *SubscriptionHandler) Update(c *gin.Context) {

	h.logger.Info("update subscription request")

	var sub model.Subscription

	if err := c.ShouldBindJSON(&sub); err != nil {
		h.logger.Error("invalid body",
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := h.service.Update(c.Request.Context(), sub)
	if err != nil {
		h.logger.Error("update service error",
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	h.logger.Info("subscription updated",
		zap.String("id", sub.ID.String()),
	)

	c.JSON(http.StatusOK, gin.H{
		"status": "updated",
	})
}

// Delete handles the HTTP request to delete a subscription.
func (h *SubscriptionHandler) Delete(c *gin.Context) {

	idParam := c.Param("id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		h.logger.Error("invalid uuid",
			zap.String("id", idParam),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid uuid",
		})
		return
	}

	h.logger.Info("delete subscription",
		zap.String("id", id.String()),
	)

	err = h.service.Delete(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("delete service error",
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	h.logger.Info("subscription deleted",
		zap.String("id", id.String()),
	)

	c.Status(http.StatusNoContent)
}

// SumByFilter handles the HTTP request to sum subscriptions by filter criteria.
func (h *SubscriptionHandler) SumByFilter(c *gin.Context) {

	userIDParam := c.Query("user_id")
	serviceName := c.Query("service")
	fromParam := c.Query("from")
	toParam := c.Query("to")

	h.logger.Info("sum by filter request",
		zap.String("user_id", userIDParam),
		zap.String("service", serviceName),
		zap.String("from", fromParam),
		zap.String("to", toParam),
	)

	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		h.logger.Error("invalid user_id",
			zap.String("user_id", userIDParam),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	from, err := time.Parse("2006-01-02", fromParam)
	if err != nil {
		h.logger.Error("invalid from date",
			zap.String("from", fromParam),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid from date"})
		return
	}

	to, err := time.Parse("2006-01-02", toParam)
	if err != nil {
		h.logger.Error("invalid to date",
			zap.String("to", toParam),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid to date"})
		return
	}

	sum, err := h.service.SumByFilter(
		c.Request.Context(),
		userID,
		serviceName,
		from,
		to,
	)

	if err != nil {
		h.logger.Error("sumByFilter service error",
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	h.logger.Info("sum by filter result",
		zap.Int("sum", sum),
	)

	c.JSON(http.StatusOK, gin.H{
		"sum": sum,
	})
}
