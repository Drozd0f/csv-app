package server

import (
	"net/http"

	"github.com/Drozd0f/csv-app/schemes"
	"github.com/gin-gonic/gin"
)

// ping godoc
// @Summary show pong
// @Tags    Healthcheck
// @Produce json
// @Success 200 {object} schemes.Response "Server is alive"
// @Router  /ping [get]
func (s *Server) ping(c *gin.Context) {
	c.JSON(http.StatusOK, schemes.Response{
		Message: "pong",
	})
}
