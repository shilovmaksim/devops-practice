package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"

	"github.com/cxrdevelop/optimization_engine/optimization_server/config"
	"github.com/cxrdevelop/optimization_engine/optimization_server/internal/optimizer"
	"github.com/cxrdevelop/optimization_engine/optimization_server/internal/python"
	"github.com/cxrdevelop/optimization_engine/pkg/logger"
	"github.com/cxrdevelop/optimization_engine/pkg/storage"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const (
	gracefulShutdownTimeoutMs = 5000
	defaultScriptTimeoutMs    = 5000
)

type Server struct {
	config        *config.Config
	wrapper       *python.Wrapper
	optimizer     optimizer.Optimizer
	storage       storage.Storage
	logger        *logger.Logger
	signalChannel chan os.Signal
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
	s.logger.Infof("optimization server is running at port %s", s.config.Application.Port)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatalf("error occurred while running http server: %s\n", err)
		}
	}()

	s.signalChannel = make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(s.signalChannel, os.Interrupt)
	<-s.signalChannel

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		s.logger.Errorf("optimization server server shutdown failed: %s", err)
	}

	s.logger.Warn("optimization server shutting down")
}

func (s *Server) SetDefaults() {
	if s.logger == nil {
		level, err := logrus.ParseLevel(s.config.Application.LogLevel)
		if err != nil {
			level = logrus.DebugLevel
		}
		s.logger = logger.New(s.config.Application.LogPath, "optimization_service", level)
	}

	if s.wrapper == nil {
		execTimeout := s.config.Script.Timeout
		if execTimeout <= 0 {
			execTimeout = defaultScriptTimeoutMs * time.Millisecond
			s.logger.Errorf("incorrect timeout value %s, switching to default value of %s", s.config.Script.Timeout, execTimeout)
		}
		scriptPath, err := filepath.Abs(s.config.Script.Path)
		if err != nil {
			s.logger.Errorf("error creating absolute path for script: %s", err)
		}
		s.wrapper = python.NewWrapper(scriptPath, execTimeout, s.logger)
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
	if s.optimizer == nil {
		s.optimizer = optimizer.NewRackOptimizer(s.wrapper, s.storage, ".", "tmp_prefix", s.logger)
	}
}

func (s *Server) SetupRoutes() *mux.Router {
	r := mux.NewRouter()
	apiPrefix := r.PathPrefix("/api/v1").Subrouter()

	wrappedHealthHandler := handlers.LoggingHandler(s.logger.Writer(),
		NewHealthHandler(s.logger),
	)
	apiPrefix.Handle("/health", wrappedHealthHandler).Methods("GET")

	wrappedOptimizationHandler := handlers.LoggingHandler(s.logger.Writer(),
		NewOptimizationHandler(s.optimizer, s.logger),
	)
	apiPrefix.Handle("/optimize", wrappedOptimizationHandler).Methods("GET")

	return r
}
