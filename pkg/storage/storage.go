package storage

type UploadResult struct {
	Filename string
	Location string
	ETag     string
}

type Storage interface {
	DownloadFiles(dir string, paths ...string) error
	UploadFiles(dir string, paths ...string) ([]UploadResult, error)
}
