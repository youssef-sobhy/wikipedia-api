package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/pborman/uuid"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/youssef1337/wikipedia-api/docs"
)

//	@title			Wikipedia Clone API
//	@version		1.0
//	@description	This API is used to get the short description of a given wikipedia article.

//	@contact.name	Youssef Sobhy
//	@contact.email	youssefsobhy22@gmail.com

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8080
//	@BasePath	/api

//	@securityDefinitions.noAuth

func main() {
	r := gin.New()

	r.Use(func(c *gin.Context) {
		reqID := uuid.New()
		c.Writer.Header().Set("X-Request-Id", reqID)
		c.Set("reqID", reqID)
		c.Next()
	})

	r.Use(gin.Logger())
	r.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		log.Printf("PANIC: %v", recovered)
		internalServerError(c, fmt.Errorf("%v", recovered))
	}))

	r.GET("/panic", func(c *gin.Context) {
		panic("PANIC")
	})

	api := r.Group("/api")
	{
		api.GET("/", health)
		api.GET("/search", search)
		api.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		api.GET("/docs", func(c *gin.Context) {
			c.Redirect(http.StatusMovedPermanently, "/api/docs/index.html")
		})
	}

	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/api/docs/index.html")
	})

	port := os.Getenv("PORT")

	if port == "" {
		port = "3000"
	}

	r.Run(":" + port)
}

// health godoc
//
//	@Summary		Check if the API is operational.
//	@Description	Check if the API is operational.
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}  SuccessResponse
//	@Failure		500	{object}  InternalServerError
//	@Router			/api [get]
func health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "operational",
	})
}

// search godoc
//
//	@Summary		Search for a short description of a person, place, or thing.
//	@Description	Search for a short description of a person, place, or thing.
//	@Accept			json
//	@Produce		json
//	@Param			query	query		string	true	"The name of the person, place, or thing you want to search for."
//	@Success		200		{object}	SuccessResponse
//	@Success    200   {object}  MissingResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		500		{object}	InternalServerError
//	@Router			/api/search [get]
func search(c *gin.Context) {
	wikipediaURL := os.Getenv("WIKIPEDIA_API_URL")

	query := c.Query("query")
	if query == "" {
		badRequestError(c, "Query parameter is required.")

		return
	}

	url := fmt.Sprintf("%s?action=query&prop=revisions&titles=%s&rvlimit=1&formatversion=2&format=json&rvprop=content", wikipediaURL, query)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		internalServerError(c, err)

		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		internalServerError(c, err)

		return
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		internalServerError(c, err)

		return
	}

	var response wikipediaResponse

	err = json.Unmarshal(body, &response)
	if err != nil {
		internalServerError(c, err)

		return
	}

	if response.Query.Pages[0].Missing {
		httpMissingHandler(c, "No wikipedia article found.")

		return
	}

	content := response.Query.Pages[0].Revisions[0].Content
	re := regexp.MustCompile(`(?mi){{short description\|(.*?)}}`)
	shortDescription := re.FindStringSubmatch(content)

	if len(shortDescription) == 0 {
		httpMissingHandler(c, "No short description found in english wikipedia.")

		return
	}

	httpSuccessHandler(c, shortDescription[1])
}

func httpSuccessHandler(c *gin.Context, shortDescription string) {
	c.JSON(http.StatusOK, SuccessResponse{
		Status: "success",
		Data:   Data{ShortDescription: shortDescription},
	})
}

func httpMissingHandler(c *gin.Context, message string) {
	c.JSON(http.StatusOK, MissingResponse{
		Status:  "success",
		Message: message,
	})
}

func httpErrorHandler(c *gin.Context, code int, message string) {
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

func badRequestError(c *gin.Context, message string) {
	httpErrorHandler(c, http.StatusBadRequest, message)
}

func internalServerError(c *gin.Context, err error) {
	httpErrorHandler(
		c,
		http.StatusInternalServerError,
		"An internal server error occurred. Please contact the developer at youssefsobhy22@gmail.com and provide the request ID.",
	)

	log.Printf("Request ID: %s, Error: %s", c.GetString("reqID"), err.Error())
}

type SuccessResponse struct {
	Status string `json:"status" example:"success"`
	Data   Data   `json:"data"`
}

type MissingResponse struct {
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"No wikipedia article found."`
	Missing bool   `json:"missing"`
}

type ErrorResponse struct {
	Status string      `json:"status" example:"error"`
	Errors []HTTPError `json:"errors"`
}

type Data struct {
	ShortDescription string `json:"short_description" example:"A short description of the person, place, or thing you searched for."`
}

type HTTPError struct {
	Code      int    `json:"code" example:"400"`
	RequestID string `json:"request_id" example:"f7a4c0c0-5b5e-4b4c-9c1f-1b5c1b5c1b5c"`
	Detail    string `json:"detail" example:"Query parameter is required"`
}

type InternalServerError struct {
	Code      int    `json:"code" example:"500"`
	RequestID string `json:"request_id" example:"f7a4c0c0-5b5e-4b4c-9c1f-1b5c1b5c1b5c"`
	Detail    string `json:"detail" example:"An internal server error occurred. Please contact the developer at youssefsobhy22@gmail.com and provide the request ID."`
}

type wikipediaResponse struct {
	Query Query `json:"query"`
}

type Query struct {
	Pages []Page `json:"pages"`
}

type Page struct {
	PageID    int        `json:"pageid"`
	Ns        int        `json:"ns"`
	Title     string     `json:"title"`
	Revisions []Revision `json:"revisions"`
	Missing   bool       `json:"missing"`
}

type Revision struct {
	Content string `json:"content"`
}
