package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/pborman/uuid"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/youssef1337/wikipedia-api/docs"

	"github.com/youssef1337/wikipedia-api/internal"
)

// @title Wikipedia API
// @description A simple API to get the short description of a wikipedia article.
// @version 1.0.0
// @BasePath /

// @contact.name Youssef Sobhy
// @contact.email youssefsobhy22@gmail.com

// @host wikipedia-api.youssefsobhy.com

// @schemes https http

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

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
		internal.InternalServerErrorHandler(c, fmt.Errorf("%v", recovered))
	}))

	api := r.Group("/api")
	{
		api.GET("", internal.Health)
		api.GET("/search", internal.Search)
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
