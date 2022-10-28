package server

import (
	"log"
	"net/http"

	"github.com/Drozd0f/csv-app/schemes"
	"github.com/gin-gonic/gin"
)

func (s *Server) registerCsvHandlers(g *gin.RouterGroup) {
	csvG := g.Group("/csv-file")
	{
		csvG.GET("", s.getTransactions)
		csvG.POST("/upload", s.uploadCsvFile)
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
		"message": "OK",
	})
}

func (s *Server) getTransactions(c *gin.Context) {
	var rrt schemes.RawRequestTransaction
	if err := c.ShouldBindQuery(&rrt); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": http.StatusText(http.StatusInternalServerError),
		})
		return
	}

	sliceT, err := s.service.GetSliceTransactions(c.Request.Context(), rrt)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": http.StatusText(http.StatusInternalServerError),
		})
		return
	}

	c.JSON(http.StatusOK, sliceT)
}
