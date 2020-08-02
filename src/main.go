package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/marcellinuskristanto/onepiece/src/configuration"
	"github.com/marcellinuskristanto/onepiece/src/middleware"
	"github.com/marcellinuskristanto/onepiece/src/route"
)

func main() {
	log.Println(fmt.Sprintf("1:%s", "LOAD CONFIG"))
	config, err := configuration.LoadConfigurations()
	if err != nil {
		log.Fatalf("An error occurred while loading the configurations: %v", err)
	}
	log.Println(fmt.Sprintf("1:%s", "LOAD CONFIG DONE"))

	// set environment
	if config.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	log.Println(fmt.Sprintf("2:%s", "SET PROD DONE"))
	// set log
	if config.Logger.Path != "" {
		f, _ := os.Create(config.Logger.Path)
		gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
		log.SetOutput(gin.DefaultWriter)
	}
	log.Println(fmt.Sprintf("3:%s", "SET LOG DONE"))

	router := gin.Default()

	log.Println(fmt.Sprintf("3:%s", "JWT MIDDLEWARE START"))
	// jwt middleware
	authMiddleware := middleware.JWTMiddleware(config.Auth)
	log.Println(fmt.Sprintf("3:%s", "JWT MIDDLEWARE DONE"))

	// Load all route
	log.Println(fmt.Sprintf("3:%s", "LOAD ROUTE START"))
	router.GET("/", func(c *gin.Context) {
		c.String(200, "OK")
	})
	router.POST("/login", authMiddleware.LoginHandler)
	router.Use(authMiddleware.MiddlewareFunc())
	{
		router.GET("/refresh_token", authMiddleware.RefreshHandler)
		route.LoadRoute(router)
	}
	log.Println(fmt.Sprintf("3:%s", "LOAD ROUTE END"))

	log.Println(fmt.Sprintf("3:%d", config.App.Listen))
	router.Run(fmt.Sprintf(":%d", config.App.Listen))
}
