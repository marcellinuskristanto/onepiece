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
	config, err := configuration.LoadConfigurations()
	if err != nil {
		log.Fatalf("An error occurred while loading the configurations: %v", err)
	}

	listenAddr := fmt.Sprintf(":%d", config.App.Port)
	// set environment
	if config.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
		listenAddr = fmt.Sprintf("0.0.0.0:%d", config.App.Port)
	}
	// set log
	if config.Logger.Path != "" {
		f, _ := os.Create(config.Logger.Path)
		gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
		log.SetOutput(gin.DefaultWriter)
	}

	router := gin.Default()

	// jwt middleware
	authMiddleware := middleware.JWTMiddleware(config.Auth)

	// Load all route
	router.GET("/", func(c *gin.Context) {
		c.String(200, "OK")
	})
	router.POST("/login", authMiddleware.LoginHandler)
	router.Use(authMiddleware.MiddlewareFunc())
	{
		router.GET("/refresh_token", authMiddleware.RefreshHandler)
		route.LoadRoute(router)
	}

	router.Run(listenAddr)
}
