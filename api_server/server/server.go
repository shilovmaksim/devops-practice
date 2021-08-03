package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/cxrdevelop/optimization_engine/api_server/config"
	"github.com/cxrdevelop/optimization_engine/api_server/internal/optimization"
	"github.com/cxrdevelop/optimization_engine/pkg/logger"
	"github.com/cxrdevelop/optimization_engine/pkg/metrics"
	"github.com/cxrdevelop/optimization_engine/pkg/storage"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const (
	gracefulShutdownTimeoutMs = 5000
)

type Server struct {
	config  *config.Config
	storage storage.Storage
	client  *optimization.Client
	logger  *logger.Logger
}

func New(config *config.Config) *Server {
	return &Server{config: config}
}

func (s *Server) Start() {
	var wait time.Duration = gracefulShutdownTimeoutMs * time.Millisecond

	s.SetDefaults()

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", s.config.Application.Port),
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 30,
		IdleTimeout:  time.Second * 60,
		Handler:      s.SetupRoutes(),
	}
	s.logger.Infof("api server is running at port %s", s.config.Application.Port)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatalf("error occurred while running http server: %s\n", err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		s.logger.Errorf("api server server shutdown failed: %s", err)
	}

	s.logger.Warn("api server shutting down")
}

func (s *Server) SetDefaults() {
	if s.logger == nil {
		level, err := logrus.ParseLevel(s.config.Application.LogLevel)
		if err != nil {
			level = logrus.DebugLevel
		}
		s.logger = logger.New(s.config.Application.LogPath, "api_service", level)
	}

	if s.storage == nil {
		switch strings.ToLower(s.config.Storage.Type) {
		case config.S3:
			s.storage = storage.NewS3Storage(s.config.Storage.Region, s.config.Storage.Bucket, s.logger)
		case config.Local:
			fallthrough
		default:
			s.storage = storage.NewFSStorage(s.config.Storage.Bucket, s.logger)
		}
	}

	if s.client == nil {
		s.client = optimization.New(s.storage, s.config.OptSrv.Endpoint, s.config.OptSrv.Port, s.logger)
	}
}

func (s *Server) SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	r.Use(metrics.PrometheusMiddleware)
	r.Handle("/metrics", metrics.Handler())

	apiPrefix := r.PathPrefix("/api/v1").Subrouter()

	wrappedHealthHandler := handlers.LoggingHandler(s.logger.Writer(),
		NewHealthHandler(s.logger),
	)
	apiPrefix.Handle("/health", wrappedHealthHandler).Methods(http.MethodGet, http.MethodOptions)

	wrappedUploadHandler := handlers.LoggingHandler(s.logger.Writer(),
		NewUploadHandler(s.storage, s.client, s.logger),
	)
	apiPrefix.Handle("/upload", wrappedUploadHandler).Methods(http.MethodPost, http.MethodOptions)

	return r
}
