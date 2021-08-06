package server

import (
	"net/http"

	"github.com/cxrdevelop/optimization_engine/api_server/internal/models"
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

// HealthHandler
// @Summary Response with a health status
// @Description Get health status from a service
// @ID health-handler
// @Accept plain
// @Produce  json
// @Success 200 {object} models.HealthResponse true
// @Failure 404
// @Router /v1/health [get]
func (h *HealthHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	writeResponse(writer, models.NewHealthResponse(), http.StatusOK, h.log)
}
