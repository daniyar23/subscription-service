package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/daniyar23/subscribe-service/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SubscriptionHandler struct {
	service SubscriptionService
}

func NewSubscriptionHandler(service SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{
		service: service,
	}
}

func (h *SubscriptionHandler) Create(c *gin.Context) {

	log.Println("Create subscription request")

	var sub model.Subscription
	if err := c.ShouldBindJSON(&sub); err != nil {
		log.Println("Create subscription: invalid request body:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.Create(c.Request.Context(), sub)
	if err != nil {
		log.Println("Create subscription: service error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	log.Println("Subscription created:", result.ID)

	c.JSON(http.StatusCreated, result)
}

func (h *SubscriptionHandler) GetByID(c *gin.Context) {

	idParam := c.Param("id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		log.Println("GetByID: invalid uuid:", idParam)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid uuid",
		})
		return
	}

	log.Println("Get subscription by id:", id)

	sub, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		log.Println("GetByID: service error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, sub)
}

func (h *SubscriptionHandler) GetAll(c *gin.Context) {

	log.Println("Get all subscriptions")

	subs, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		log.Println("GetAll: service error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, subs)
}

func (h *SubscriptionHandler) GetByUserID(c *gin.Context) {

	userIDParam := c.Param("user_id")

	if userIDParam == "" {
		log.Println("GetByUserID: user_id is empty")
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		log.Println("GetByUserID: invalid uuid:", userIDParam)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid uuid"})
		return
	}

	log.Println("Get subscriptions for user:", userID)

	subs, err := h.service.GetByUserID(c.Request.Context(), userID)
	if err != nil {
		log.Println("GetByUserID: service error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, subs)
}

func (h *SubscriptionHandler) Update(c *gin.Context) {

	log.Println("Update subscription request")

	var sub model.Subscription

	if err := c.ShouldBindJSON(&sub); err != nil {
		log.Println("Update: invalid body:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := h.service.Update(c.Request.Context(), sub)
	if err != nil {
		log.Println("Update: service error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	log.Println("Subscription updated:", sub.ID)

	c.JSON(http.StatusOK, gin.H{
		"status": "updated",
	})
}

func (h *SubscriptionHandler) Delete(c *gin.Context) {

	idParam := c.Param("id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		log.Println("Delete: invalid uuid:", idParam)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid uuid",
		})
		return
	}

	log.Println("Delete subscription:", id)

	err = h.service.Delete(c.Request.Context(), id)
	if err != nil {
		log.Println("Delete: service error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	log.Println("Subscription deleted:", id)

	c.JSON(http.StatusOK, gin.H{
		"status": "deleted",
	})
}

func (h *SubscriptionHandler) SumByFilter(c *gin.Context) {

	userIDParam := c.Query("user_id")
	serviceName := c.Query("service")
	fromParam := c.Query("from")
	toParam := c.Query("to")

	log.Println("SumByFilter request:", userIDParam, serviceName, fromParam, toParam)

	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		log.Println("SumByFilter: invalid user_id:", userIDParam)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	from, err := time.Parse("2006-01-02", fromParam)
	if err != nil {
		log.Println("SumByFilter: invalid from date:", fromParam)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid from date"})
		return
	}

	to, err := time.Parse("2006-01-02", toParam)
	if err != nil {
		log.Println("SumByFilter: invalid to date:", toParam)
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
		log.Println("SumByFilter: service error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	log.Println("SumByFilter result:", sum)

	c.JSON(http.StatusOK, gin.H{
		"sum": sum,
	})
}
