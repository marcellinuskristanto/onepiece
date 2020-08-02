package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()

	// Load all route
	router.GET("/", func(c *gin.Context) {
		c.String(200, "OK")
	})

	router.Run(fmt.Sprintf(":%d", 80))
}
