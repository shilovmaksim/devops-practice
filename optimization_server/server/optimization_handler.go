package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/cxrdevelop/optimization_engine/optimization_server/internal/models"
	"github.com/cxrdevelop/optimization_engine/optimization_server/internal/optimizer"
	"github.com/cxrdevelop/optimization_engine/pkg/logger"
)

const (
	ErrMsgEmptyRequest = "empty request body"
	ErrMsgJsonParse    = "error parsing json body"
	ErMsgJsonValidate  = "error validating json body"
	ErrMsgInternal     = "internal error"
	ErrMsgDownload     = "failed to download files"
	ErrMsgScript       = "script error"
	ErrMsgUpload       = "failed to upload the result"
)

type OptimizationHandler struct {
	optimizer optimizer.Optimizer
	log       *logger.Logger
}

func NewOptimizationHandler(optimizer optimizer.Optimizer, log *logger.Logger) *OptimizationHandler {
	return &OptimizationHandler{
		optimizer: optimizer,
		log:       log,
	}
}

// Optimize
// @Summary Run optimization script
// @Description Download files from the bucket and run optimization script
// @ID optimization-handler
// @Accept  json
// @Produce  json
// @Param message body models.OptimizationRequest true "Path to file in the bucket"
// @Success 200 {object} models.OptimizationResponse "Successful run result"
// @Failure 500 {object} models.ErrorResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /v1/optimize [post]
func (h *OptimizationHandler) ServeHTTP(writer http.ResponseWriter, r *http.Request) {
	req := models.OptimizationRequest{}
	if r.Body == nil {
		h.log.Errorf("empty request body")
		writeResponse(writer, models.NewErrorResponse(ErrMsgEmptyRequest), http.StatusBadRequest, h.log)
		return
	}
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Errorf("failed to parse json body, err: %s", err)
		writeResponse(writer, models.NewErrorResponse(ErrMsgJsonParse), http.StatusBadRequest, h.log)
		return
	}

	if validErrs := req.Validate(); len(validErrs) > 0 {
		h.log.Errorf("failed to validate input json body, err: %s", validErrs)
		writeResponse(writer, models.NewErrorResponse(ErMsgJsonValidate), http.StatusBadRequest, h.log)
		return
	}

	res, err := h.optimizer.Execute(req.Filename)

	switch {
	case err == nil:
		writeResponse(writer, models.NewOptimizationResponse(
			res.Location,
			res.Filename,
			res.ETag,
			int64(res.ExecutionTime.Milliseconds())), http.StatusOK, h.log)

	case errors.Is(err, optimizer.ErrDownload):
		writeResponse(writer, models.NewErrorResponse(ErrMsgDownload), http.StatusInternalServerError, h.log)

	case errors.Is(err, optimizer.ErrOptimize):
		writeResponse(writer, models.NewErrorResponse(ErrMsgScript), http.StatusInternalServerError, h.log)

	case errors.Is(err, optimizer.ErrUpload):
		writeResponse(writer, models.NewErrorResponse(ErrMsgUpload), http.StatusInternalServerError, h.log)

	case errors.Is(err, optimizer.ErrEnvCreate):
		fallthrough
	default:
		writeResponse(writer, models.NewErrorResponse(ErrMsgInternal), http.StatusInternalServerError, h.log)
	}
}
