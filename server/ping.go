package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Message string `json:"message"`
}

// ping godoc
// @Summary show pong
// @Tags    Healthcheck
// @Produce json
// @Success 200 {object} Response "Server is alive"
// @Router  /ping [get]
func (s *Server) ping(c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Message: "pong",
	})
}
