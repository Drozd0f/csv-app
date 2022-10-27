package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) registerCsvHandlers(g *gin.RouterGroup) {
	csvG := g.Group("/csv-file")
	{
		csvG.POST("", s.uploadCsvFile)
	}
}

func (s *Server) uploadCsvFile(c *gin.Context) {
	filePtr, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": err.Error(),
		})
		return
	}

	if err = s.service.UploadCsvFile(c.Request.Context(), filePtr); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "File was uploaded to database",
	})
}
