package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
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

// @host wikipedia.youssefsobhy.com

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

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://wikipedia.youssefsobhy.com"}
	config.AllowMethods = []string{"GET"}

	r.Use(cors.New(config))

	v1 := r.Group("/api/v1")
	{
		v1.GET("", internal.Health)
		v1.GET("/search", internal.Search)
		v1.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		v1.GET("/docs", func(c *gin.Context) {
			c.Redirect(http.StatusMovedPermanently, "/api/v1/docs/index.html")
		})
	}

	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/api/v1/docs/index.html")
	})

	port := os.Getenv("PORT")

	if port == "" {
		port = "3000"
	}

	r.Run(":" + port)
}
