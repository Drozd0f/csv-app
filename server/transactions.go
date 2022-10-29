package server

import (
	"errors"
	"log"
	"net/http"

	"github.com/Drozd0f/csv-app/schemes"
	"github.com/Drozd0f/csv-app/service"
	"github.com/gin-gonic/gin"
)

func (s *Server) registerCsvHandlers(g *gin.RouterGroup) {
	csvG := g.Group("/csv-file")
	{
		csvG.GET("", s.getTransactions)
		csvG.POST("/upload", s.uploadCsvFile)
	}
}

// uploadCsvFile godoc
// @Summary upload csv file to database
// @Tags    Transactions
// @Produce json
// @Param   file formData file     true "file to upload"
// @Success 200  {object} Response "File is uploaded"
// @Failure 422  {object} Response "invalid content type provided"
// @Failure 400  {object} Response "onep file error"
// @Failure 400  {object} Response "invalid file signature"
// @Failure 400  {object} Response "transaction already exist"
// @Failure 500  {object} Response
// @Router  /csv-file/upload [post]
func (s *Server) uploadCsvFile(c *gin.Context) {
	filePtr, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, Response{
			Message: err.Error(),
		})
		return
	}

	if err = s.service.UploadCsvFile(c.Request.Context(), filePtr); err != nil {
		switch {
		case errors.Is(err, service.ErrOpenFile),
			errors.Is(err, service.ErrParsing),
			errors.Is(err, service.ErrTransactionExist):
			c.JSON(http.StatusBadRequest, Response{
				Message: err.Error(),
			})
			return
		default:
			log.Println(err)
			c.JSON(http.StatusInternalServerError, Response{
				Message: http.StatusText(http.StatusInternalServerError),
			})
			return
		}
	}

	c.JSON(http.StatusOK, Response{
		Message: "OK",
	})
}

// getTransactions godoc
// @Summary show slice transactions
// @Tags    Transactions
// @Produce json
// @Param   transaction_id    query    int                         false "Search by transaction_id"
// @Param   terminal_id       query    []int                       false "Search by terminal_id (possible to specify several ids at the same time)"
// @Param   status            query    string                      false "Search by status"       Enums(accepted, declined)
// @Param   payment_type      query    string                      false "Search by payment_type" Enums(card, cash)
// @Param   from              query    string                      false "From date inclusive"    Format(date) example(2022-08-12)
// @Param   to                query    string                      false "To date not inclusive"  Format(date) example(2022-09-01)
// @Param   payment_narrative query    string                      false "Search on the partially specified payment_narrative"
// @Success 200               {array}  schemes.ResponseTransaction "Show slice transactions"
// @Failure 500               {object} Response
// @Router  /csv-file [get]
func (s *Server) getTransactions(c *gin.Context) {
	var rrt schemes.RawRequestTransaction
	_ = c.ShouldBindQuery(&rrt)

	sliceT, err := s.service.GetSliceTransactions(c.Request.Context(), rrt)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, Response{
			Message: http.StatusText(http.StatusInternalServerError),
		})
		return
	}

	c.JSON(http.StatusOK, sliceT)
}
