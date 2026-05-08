package handler

import (
	"rebellion/internal/response"

	"github.com/gin-gonic/gin"
)

func HealthCheck(c *gin.Context) {
	response.Success(c, 0, "Server is healthy . . .", nil)
}
