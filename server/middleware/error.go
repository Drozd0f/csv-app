package middleware

import (
	"errors"
	"log"
	"net/http"

	errs "github.com/Drozd0f/csv-app/errors"
	"github.com/Drozd0f/csv-app/schemes"
	"github.com/Drozd0f/csv-app/service"
	"github.com/gin-gonic/gin"
)

var handleErrors map[error]int = map[error]int{
	service.ErrOpenFile:         http.StatusBadRequest,
	service.ErrParsing:          http.StatusBadRequest,
	service.ErrTransactionExist: http.StatusBadRequest,
}

func ErrorHandler(c *gin.Context) {
	c.Next()

	var mesErr *errs.ErrorWithMessage

	for _, err := range c.Errors {
		if errors.As(err, &mesErr) {
			status, ok := handleErrors[mesErr.Err]
			if ok {
				c.JSON(status, schemes.Response{Message: err.Error()})
				return
			}
		}

		log.Println("internal server error", err.Error())
		c.JSON(http.StatusInternalServerError, schemes.Response{
			Message: http.StatusText(http.StatusInternalServerError),
		})
	}
}
