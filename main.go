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

// @title           Wikipedia Clone API
// @version         1.0
// @description     This API is used to get the short description of a given wikipedia article.

// @contact.name   Youssef Sobhy
// @contact.email  youssefsobhy22@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.noAuth

func main() {
	wikipediaURL := os.Getenv("WIKIPEDIA_API_URL")

	r := gin.New()

	r.Use(func(c *gin.Context) {
		reqID := uuid.New()
		c.Writer.Header().Set("X-Request-Id", reqID)
		c.Set("reqID", reqID)
		c.Next()
	})

	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/api")
	})

	r.GET("/api", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "operational",
		})
	})

	r.GET("/wikipedia/search", func(c *gin.Context) {
		query := c.Query("query")
		if query == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":            "error",
				"short_description": "query parameter is required",
			})

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
			c.JSON(http.StatusOK, gin.H{
				"status":            "success",
				"short_description": "no short description found in english wikipedia",
			})

			return
		}

		content := response.Query.Pages[0].Revisions[0].Content
		re := regexp.MustCompile(`(?mi){{short description\|(.*?)}}`)
		shortDescription := re.FindStringSubmatch(content)

		if len(shortDescription) == 0 {
			c.JSON(http.StatusOK, gin.H{
				"status":            "success",
				"short_description": "no short description found in english wikipedia",
			})

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":            "success",
			"short_description": shortDescription[1],
		})
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	r.Run(":" + port)
}

func internalServerError(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"status":            "error",
		"short_description": "Something went wrong, please contact the developer at youssefsobhy22@gmail.com and provide the following request ID -> " + c.GetString("reqID"),
	})

	log.Printf("Request ID: %s, Error: %s", c.GetString("reqID"), err.Error())
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
