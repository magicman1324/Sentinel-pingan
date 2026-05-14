package rest

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/pingan/monitor-backend/internal/service"
	"github.com/redis/go-redis/v9"
)

func RegisterRoutes(r *gin.Engine, svc *service.Service, db *sql.DB, rdb *redis.Client) {
	h := &handler{svc: svc, db: db, rdb: rdb}

	api := r.Group("/api/v1")
	{
		rules := api.Group("/rules")
		{
			rules.GET("", h.ListRules)
			rules.POST("", ValidateRule(), h.CreateRule)
			rules.PUT("/:id", ValidateRule(), h.UpdateRule)
			rules.DELETE("/:id", h.DeleteRule)
		}

		alerts := api.Group("/alerts")
		{
			alerts.GET("", h.ListAlerts)
			alerts.POST("/:id/resolve", h.ResolveAlert)
		}

		api.GET("/health", h.Health)
	}
}
