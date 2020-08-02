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

	// set environment
	if config.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
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
	router.POST("/login", authMiddleware.LoginHandler)
	router.Use(authMiddleware.MiddlewareFunc())
	{
		router.GET("/refresh_token", authMiddleware.RefreshHandler)
		route.LoadRoute(router)
	}

	router.Run(fmt.Sprintf(":%d", config.App.Listen))
}
