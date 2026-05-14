package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/pingan/monitor-backend/internal/service"
)

func RegisterRoutes(r *gin.Engine, svc *service.Service) {
	h := &handler{svc: svc}

	api := r.Group("/api/v1")
	{
		api.GET("/rules", h.ListRules)
		api.POST("/rules", h.CreateRule)
		api.PUT("/rules/:id", h.UpdateRule)
		api.DELETE("/rules/:id", h.DeleteRule)

		api.GET("/alerts", h.ListAlerts)
		api.POST("/alerts/:id/resolve", h.ResolveAlert)

		api.GET("/health", h.Health)
	}
}
