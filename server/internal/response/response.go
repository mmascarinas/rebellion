package response

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	ErrorID string `json:"error_id,omitempty"`
}

func Success(c *gin.Context, status int, message string, data any) {
	if status == 0 {
		status = 200
	}

	if message == "" {
		message = "success"
	}

	r := &Response{Status: status, Message: message, Data: data}
	c.JSON(status, r)
}

func Error(c *gin.Context, status int, message string, err any) {
	if status == 0 {
		status = 500
	}

	if message == "" {
		message = "internal server error"
	}

	errorID := uuid.New().String()

	slog.Error("request error",
		"error_id", errorID,
		"status", status,
		"message", message,
		"error", err,
	)

	r := &Response{Status: status, Message: message, ErrorID: errorID}
	c.JSON(status, r)
}
