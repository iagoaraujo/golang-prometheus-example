package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vtex/golang-server/metrics"
)

func main() {
	r := gin.Default()
	router := r.Group("/hello")
	{
		router.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "Hello World!!",
			})
		})

		router.GET("/myself", func(c *gin.Context) {
			_, err := http.Get("http://localhost:8080/hello")
			if err != nil {
				c.JSON(500, err)
			}

			c.JSON(200, gin.H{
				"message": "Inception!!",
			})
		})
	}

	r.GET("/metrics", metrics.PrometheusHandler())

	r.Run() // listen and serve on 0.0.0.0:8080
}
