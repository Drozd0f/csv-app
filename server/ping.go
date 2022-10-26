package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
