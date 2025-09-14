package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/linn221/RequesterBackend/api"
	"github.com/linn221/RequesterBackend/config"
	"github.com/linn221/RequesterBackend/middlewares"
	"github.com/linn221/RequesterBackend/utils"
)

func main() {
	// Create and configure the application
	db := config.ConnectDB()

	// Get server configuration from environment
	port := utils.GetEnv("PORT", "8081")
	host := utils.GetEnv("HOST", "localhost")
	secretConfig := middlewares.SecretConfig{
		Host:        "http://localhost:" + port,
		SecretPath:  "start-session",
		RedirectUrl: "http://localhost:5173",
		SecretFunc: func() string {
			return utils.GenerateRandomString(20)
		},
	}
	app := api.App{
		DB:               db,
		SecretMiddleware: secretConfig.Middleware(),
	}

	// Start server in a goroutine
	go func() {
		if err := app.Start(host, port); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Server shutting down...")

	// Graceful shutdown would go here if needed
	// ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	// defer cancel()
	// if err := server.Shutdown(ctx); err != nil {
	//     log.Fatal("Server forced to shutdown:", err)
	// }
}
