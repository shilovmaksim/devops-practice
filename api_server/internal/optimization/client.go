package optimization

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/cxrdevelop/optimization_engine/api_server/internal/models"
	"github.com/cxrdevelop/optimization_engine/pkg/logger"
	"github.com/cxrdevelop/optimization_engine/pkg/storage"
)

const (
	optUrl = "http://%s:%s/api/v1/optimize"
)

type Client struct {
	storage storage.Storage
	optUrl  string
	log     *logger.Logger
}

type Response struct {
	Filepath      string
	Location      string
	ETag          string
	ExecutionTime time.Duration
}

func New(storage storage.Storage, endpoint string, port string, log *logger.Logger) *Client {
	return &Client{
		storage: storage,
		optUrl:  fmt.Sprintf(optUrl, endpoint, port),
		log:     log,
	}
}

func (c *Client) PostOptimize(filename string) (*Response, error) {
	var requestBody bytes.Buffer
	optimizationRequest := models.NewOptimizationRequest(filename)

	if err := json.NewEncoder(&requestBody).Encode(&optimizationRequest); err != nil {
		return nil, err
	}

	request, err := http.NewRequestWithContext(context.Background(), "GET", c.optUrl, &requestBody)
	if err != nil {
		return nil, err
	}

	httpClient := &http.Client{}
	response, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("optimization server error, '%s", response.Body)
	}

	optimizationResponse := models.OptimizationResponse{}
	if err := json.NewDecoder(response.Body).Decode(&optimizationResponse); err != nil {
		return nil, err
	}
	c.log.Debugf("optimization response received, filename: '%s', location: '%s'", optimizationResponse.BucketFilename, optimizationResponse.BucketLocation)

	return &Response{
		Filepath:      optimizationResponse.BucketFilename,
		Location:      optimizationResponse.BucketLocation,
		ETag:          optimizationResponse.BucketETag,
		ExecutionTime: time.Duration(optimizationResponse.ExecutionTime),
	}, nil
}
