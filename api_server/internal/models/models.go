package models

type (
	// HealthResponse - a model for health check api
	HealthResponse struct {
		Health bool `json:"health"`
	}

	// ErrorResponse - a model used for general error responses
	ErrorResponse struct {
		Text string `json:"text"`
	}

	// UploadResponse - a model used to respond to the upload API request
	UploadResponse struct {
		BucketFileName string `json:"filename"`
		BucketLocation string `json:"location"`
		BucketETag     string `json:"etag"`
	}

	// OptimizationRequest - a model used to form a request to the optimization service
	OptimizationRequest struct {
		BucketFilename string `json:"filename"`
	}

	// OptimizationResponse - a model representing optimization service response
	OptimizationResponse struct {
		BucketFilename string `json:"filename"`
		BucketLocation string `json:"location"`
		BucketETag     string `json:"etag"`
		ExecutionTime  int64  `json:"executionTime"`
	}
)

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

func NewUploadResponse(bucketLocation, bucketFilename, bucketEtag string) UploadResponse {
	return UploadResponse{
		BucketLocation: bucketLocation,
		BucketFileName: bucketFilename,
		BucketETag:     bucketEtag,
	}
}

func NewOptimizationRequest(filename string) OptimizationRequest {
	return OptimizationRequest{
		BucketFilename: filename,
	}
}

func NewOptimizationResponse() OptimizationResponse {
	return OptimizationResponse{}
}
