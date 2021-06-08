package server

import (
	"net/http"

	"github.com/cxrdevelop/optimization_engine/optimization_server/internal/models"
	"github.com/cxrdevelop/optimization_engine/pkg/logger"
)

type HealthHandler struct {
	log *logger.Logger
}

func NewHealthHandler(log *logger.Logger) *HealthHandler {
	return &HealthHandler{
		log: log,
	}
}

func (h *HealthHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	writeResponse(writer, models.NewHealthResponse(), http.StatusOK, h.log)
}
