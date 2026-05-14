package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/pingan/monitor-backend/internal/repository"
	"github.com/pingan/monitor-backend/internal/service"
	"github.com/redis/go-redis/v9"
)

func RegisterRoutes(r *gin.Engine, svc *service.Service, db *sqlx.DB, rdb *redis.Client) {
	h := &handler{svc: svc, db: db.DB, rdb: rdb}
	silenceRepo := repository.NewSilenceRepository(db)
	sh := &silenceHandler{repo: silenceRepo}

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

		silences := api.Group("/silences")
		{
			silences.GET("", sh.List)
			silences.POST("", sh.Create)
			silences.DELETE("/:id", sh.Delete)
		}

		api.GET("/health", h.Health)
	}
}
