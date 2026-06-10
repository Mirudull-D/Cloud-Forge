package api

import (
	"CloudHub/cmd/worker"
	"CloudHub/internal/deployments"
	"CloudHub/internal/docker"
	"CloudHub/internal/queue"
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/redis/go-redis/v9"
)

// Application app struct def
type Application struct {
	addr string
	db   *sql.DB
	rdb  *redis.Client
}

func NewApplication(addr string, db *sql.DB, rdb *redis.Client) *Application {
	return &Application{addr, db, rdb}
}

func (app *Application) Mount() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.ClientIPFromRemoteAddr)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.StripSlashes)

	r.Use(middleware.Timeout(60 * time.Second))

	deploymentStore := deployments.NewStore(app.db)
	redisStore := queue.NewRedisStore(app.rdb)
	deploymentHandler := deployments.NewHandler(deploymentStore, redisStore)
	deploymentHandler.RegisterRoutes(r)

	cli, _ := docker.NewDockerClient()

	Worker := worker.NewWorker(redisStore, deploymentStore, cli)
	go Worker.Run(context.Background())

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
