package rest

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var validOperators = map[string]bool{
	">": true, ">=": true, "<": true, "<=": true, "==": true, "!=": true,
}

type ruleInput struct {
	Name   *string `json:"name"`
	Metric *string `json:"metric"`
	Op     *string `json:"operator"`
}

func ValidateRule() gin.HandlerFunc {
	return func(c *gin.Context) {
		var in ruleInput
		if err := c.ShouldBindJSON(&in); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
			return
		}
		if in.Name == nil || strings.TrimSpace(*in.Name) == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "name is required"})
			return
		}
		if in.Metric == nil || strings.TrimSpace(*in.Metric) == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "metric is required"})
			return
		}
		if in.Op != nil && !validOperators[*in.Op] {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "operator must be one of: >, >=, <, <=, ==, !=",
			})
			return
		}
	}
}
