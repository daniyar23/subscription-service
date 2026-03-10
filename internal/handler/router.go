package handler

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine, h *SubscriptionHandler) {
	api := r.Group("/api")
	v1 := api.Group("/v1")
	v1.POST("/subscriptions", h.Create)
	v1.GET("/subscriptions", h.GetAll)

	v1.GET("/subscriptions/:id", h.GetByID)

	v1.GET("/users/:user_id/subscriptions", h.GetByUserID)

	v1.PUT("/subscriptions", h.Update)

	v1.DELETE("/subscriptions/:id", h.Delete)

	v1.GET("/subscriptions/sum", h.SumByFilter)

	v1.GET("", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Subscription API v1",
		})
	})

}
