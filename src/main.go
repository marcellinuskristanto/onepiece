package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/marcellinuskristanto/onepiece/src/configuration"
	"github.com/marcellinuskristanto/onepiece/src/helper"
	"github.com/marcellinuskristanto/onepiece/src/middleware"
	"github.com/marcellinuskristanto/onepiece/src/route"
)

func main() {
	var err error
	err = configuration.LoadConfigurations()
	if err != nil {
		log.Fatalf("An error occurred while loading the configurations: %v", err)
	}

	config := configuration.GetConfig()

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

	// init minio
	err = helper.InitMinio(config.App.MinioUrl, config.App.MinioUser, config.App.MinioSecret)
	if err != nil {
		log.Fatalf("An error occurred while init minio: %v", err)
	}

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
