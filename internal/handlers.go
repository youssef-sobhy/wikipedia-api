package internal

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HttpSuccessHandler(c *gin.Context, shortDescription string) {
	c.JSON(http.StatusOK, SuccessResponse{
		Status: "success",
		Data:   Data{ShortDescription: shortDescription},
	})
}

func HttpMissingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, MissingResponse{
		Status:  "success",
		Message: "No wikipedia article found.",
		Missing: true,
	})
}

func HttpNoDescriptionHandler(c *gin.Context) {
	c.JSON(http.StatusOK, NoDescriptionResponse{
		Status:  "success",
		Message: "No short description found for this article.",
		Missing: false,
	})
}

func HttpErrorHandler(c *gin.Context, code int, message string) {
	c.JSON(code, ErrorResponse{
		Status: "error",
		Errors: []HTTPError{
			{
				Code:      code,
				RequestID: c.GetString("reqID"),
				Detail:    message,
			},
		},
	})
}

func BadRequestErrorHandler(c *gin.Context, message string) {
	HttpErrorHandler(c, http.StatusBadRequest, message)
}

func WikipediaApiErrorHandler(c *gin.Context, httpStatusCode int) {
	HttpErrorHandler(
		c,
		http.StatusInternalServerError,
		fmt.Sprintf("An error occurred while communicating with the wikipedia API with http code %v. Please find more information at https://en.wikipedia.org/w/api.php.", httpStatusCode),
	)
}

func InternalServerErrorHandler(c *gin.Context, err error) {
	HttpErrorHandler(
		c,
		http.StatusInternalServerError,
		"An internal server error occurred. Please contact the developer at youssefsobhy22@gmail.com and provide the request ID.",
	)

	log.Printf("Request ID: %s, Error: %s", c.GetString("reqID"), err.Error())
}
