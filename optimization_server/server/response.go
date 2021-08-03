package server

import (
	"encoding/json"
	"net/http"

	"github.com/cxrdevelop/optimization_engine/pkg/logger"
)

func writeResponse(writer http.ResponseWriter, data interface{}, statusCode int, log *logger.Logger) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		log.Errorf("failed to marshal response (%s)", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)
	writer.Write(jsonBytes)
}
