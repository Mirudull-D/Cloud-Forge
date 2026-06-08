package api

import (
	"CloudHub/internal/deployments"
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Application app struct def
type Application struct {
	addr string
	db   *sql.DB
}

func NewApplication(addr string, db *sql.DB) *Application {
	return &Application{addr, db}
}

func (app *Application) Mount() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.ClientIPFromRemoteAddr)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	deploymentStore := deployments.NewStore(app.db)
	deploymentHandler := deployments.NewHandler(deploymentStore)
	deploymentHandler.RegisterRoutes(r)

	return r
}

func (app *Application) Start(h http.Handler) error {
	server := &http.Server{
		Addr:         app.addr,
		Handler:      h,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	log.Printf("Starting server on %s", app.addr)
	return server.ListenAndServe()
}
