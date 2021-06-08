package models

import "net/url"

type (
	// HealthResponse - a model for health check api
	HealthResponse struct {
		Health bool `json:"health"`
	}

	// ErrorResponse - a model used for general error responses
	ErrorResponse struct {
		Text string `json:"text"`
	}

	OptimizationRequest struct {
		Filename string `json:"filename"`
	}
	OptimizationResponse struct {
		BucketFilename string `json:"filename"`
		BucketLocation string `json:"location"`
		BucketETag     string `json:"etag"`
		ExecutionTime  int64  `json:"executionTime"`
	}
)

func (o *OptimizationRequest) Validate() url.Values {
	errs := url.Values{}

	if len(o.Filename) == 0 {
		errs.Add("filename", "not set")
	}
	// TODO: add regexp to check for extension

	return errs
}

func NewHealthResponse() HealthResponse {
	return HealthResponse{
		Health: true,
	}
}

func NewErrorResponse(errMsg string) ErrorResponse {
	return ErrorResponse{
		Text: errMsg,
	}
}

func NewOptimizationResponse(bucketLocation, bucketFilename, bucketEtag string, execTime int64) OptimizationResponse {
	return OptimizationResponse{
		BucketLocation: bucketLocation,
		BucketFilename: bucketFilename,
		BucketETag:     bucketEtag,
		ExecutionTime:  execTime,
	}
}
