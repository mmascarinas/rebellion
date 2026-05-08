package validators

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ValidateBody[T any]() gin.HandlerFunc {
	return func(c *gin.Context) {
		var body T
		if err := c.ShouldBindJSON(&body); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error":   "validation failed",
				"details": err.Error(),
			})
			return
		}

		c.Set("payload", body)
		c.Next()
	}
}
