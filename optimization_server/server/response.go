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
		jsonBytes = []byte(`{}`)
	}

	writer.WriteHeader(statusCode)
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(jsonBytes)
}
