package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Drozd0f/csv-app/schemes"
	"github.com/gin-gonic/gin"
)

func (s *Server) registerCsvHandlers(g *gin.RouterGroup) {
	csvG := g.Group("/csv-file")
	{
		csvG.GET("", s.getTransactions)
		csvG.GET("/download", s.downloadCsvFile)
		csvG.POST("/upload", s.uploadCsvFile)
	}
}

// uploadCsvFile godoc
// @Summary upload csv file to database
// @Tags    Transactions
// @Produce json
// @Param   file formData file             true "file to upload"
// @Success 200  {object} schemes.Response "File is uploaded"
// @Failure 422  {object} schemes.Response "invalid content type provided"
// @Failure 400  {object} schemes.Response "onep file error"
// @Failure 400  {object} schemes.Response "invalid file signature"
// @Failure 400  {object} schemes.Response "transaction already exist"
// @Failure 500  {object} schemes.Response
// @Router  /csv-file/upload [post]
func (s *Server) uploadCsvFile(c *gin.Context) {
	filePtr, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, schemes.Response{
			Message: err.Error(),
		})
		return
	}

	if err = s.service.UploadCsvFile(c.Request.Context(), filePtr); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, schemes.Response{
		Message: "OK",
	})
}

// getTransactions godoc
// @Summary show slice transactions
// @Tags    Transactions
// @Produce json
// @Param   transaction_id    query    int                 false "Search by transaction_id"
// @Param   terminal_id       query    []int               false "Search by terminal_id (possible to specify several ids at the same time)"
// @Param   status            query    string              false "Search by status"       Enums(accepted, declined)
// @Param   payment_type      query    string              false "Search by payment_type" Enums(card, cash)
// @Param   from              query    string              false "From date inclusive"    Format(date) example(2022-08-12)
// @Param   to                query    string              false "To date not inclusive"  Format(date) example(2022-09-01)
// @Param   payment_narrative query    string              false "Search on the partially specified payment_narrative"
// @Success 200               {array}  schemes.Transaction "Show slice transactions"
// @Failure 500               {object} schemes.Response
// @Router  /csv-file [get]
func (s *Server) getTransactions(c *gin.Context) {
	var rft schemes.RawTransactionFilter
	_ = c.ShouldBindQuery(&rft)

	sliceT, err := s.service.GetFilteredTransactions(c.Request.Context(), rft)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, sliceT)
}

// downloadCsvFile godoc
// @Summary download csv file to database
// @Tags    Transactions
// @Produce text/csv
// @Success 200 {file}   binary "return csv file"
// @Failure 500 {object} schemes.Response
// @Router  /csv-file/download [get]
func (s *Server) downloadCsvFile(c *gin.Context) {
	c.Writer.Header().Set("content-type", "text/csv")
	c.Writer.Header().Set("content-disposition",
		fmt.Sprintf("attachment; filename=\"dump_%s.csv\"", time.Now().Format("2006-01-02 15:04:05")))

	if err := s.service.DownloadCsvFile(c.Request.Context(), c.Writer); err != nil {
		c.Error(err)
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
}
