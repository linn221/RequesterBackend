package api

import (
	"fmt"
	"net/http"

	"github.com/linn221/RequesterBackend/middlewares"
	"gorm.io/gorm"
)

type App struct {
	DB               *gorm.DB
	SecretMiddleware func(http.Handler) http.Handler
}

// Start starts the HTTP server
func (app *App) Start(host, port string) error {
	mux := app.RegisterRoutes()
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", host, port),
		Handler: app.SecretMiddleware(middlewares.Recovery(middlewares.LoggingMiddleware(mux))),
	}

	fmt.Printf("Server starting on http://%s:%s\n", host, port)
	return server.ListenAndServe()
}
