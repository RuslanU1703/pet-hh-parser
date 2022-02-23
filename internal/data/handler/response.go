package handler

import "github.com/gin-gonic/gin"

type errorResponse struct {
	Message string
}

func newErrorResponse(c *gin.Context, statusCode int, message string) {
	c.AbortWithStatusJSON(statusCode, errorResponse{Message: message})
}
