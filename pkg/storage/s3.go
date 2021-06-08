package storage

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/cxrdevelop/optimization_engine/pkg/logger"
)

type S3Storage struct {
	svc        s3iface.S3API
	downloader *s3manager.Downloader
	uploader   *s3manager.Uploader
	bucket     string
	log        *logger.Logger
}

var _ Storage = (*S3Storage)(nil)

// New S3Client creates s3 service and downloader or panics on error
func NewS3Storage(region string, bucket string, log *logger.Logger) *S3Storage {
	// verify aws auth and panic on error
	if err := verifyEnv("AWS_ACCESS_KEY_ID", "AWS_SECRET_ACCESS_KEY"); err != nil {
		panic(err)
	}

	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
	if err != nil {
		panic(fmt.Errorf("failed to create aws session for region '%s' with '%w'", region, err))
	}
	rawClient := s3.New(awsSession)

	return &S3Storage{
		svc:        rawClient,
		downloader: s3manager.NewDownloaderWithClient(rawClient),
		uploader:   s3manager.NewUploaderWithClient(rawClient),
		bucket:     bucket,
		log:        log,
	}
}

// DownloadFiles function takes S3 bucket name and a remote filename.
// It creates a new file with the same name in the dir. If download fails the file will be empty.
func (s *S3Storage) DownloadFiles(dir string, paths ...string) error {
	for _, path := range paths {
		if err := s.download(dir, path, ""); err != nil {
			return err
		}
	}
	return nil
}

// UploadFiles function takes local dir name and a local file name.
// It takes a file from the dir and uploads it to the bucket with the same name.
func (s *S3Storage) UploadFiles(dir string, paths ...string) ([]UploadResult, error) {
	res := make([]UploadResult, 0, len(paths))
	for _, path := range paths {
		if r, err := s.upload(dir, path, ""); err != nil {
			return nil, err
		} else {
			res = append(res, *r)
		}
	}
	return res, nil
}

func (s *S3Storage) download(dir string, remoteFilename string, localFilename string) error {
	if localFilename == "" {
		localFilename = remoteFilename
	}
	path := getFilePath(dir, localFilename)
	s.log.Debugf("Creating tmp file '%s'...", path)
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create tmp file '%s', error: `%w`", path, err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			s.log.Warnf("error closing file: %s", err)
		}
	}()

	s.log.Debugf("Downloading file '%s' from s3...", remoteFilename)

	fileSize, err := s.downloader.Download(file, &s3.GetObjectInput{
		Bucket: &s.bucket,
		Key:    &remoteFilename,
	})
	s.log.Debugf("Downloaded file size: '%d'", fileSize)

	if err != nil {
		return fmt.Errorf("unable to download '%s' from bucket '%s' with '%w'", localFilename, s.bucket, err)
	}

	return nil
}

func (s *S3Storage) upload(dir string, localFilename string, remoteFilename string) (*UploadResult, error) {
	if remoteFilename == "" {
		remoteFilename = localFilename
	}
	path := getFilePath(dir, localFilename)
	s.log.Debugf("Opening a file '%s'...", path)

	file, err := os.Open(path)
	if err != nil {
		s.log.Fatalf("failed to create file '%s'", path)
		return nil, err
	}
	defer func() {
		if err := file.Close(); err != nil {
			s.log.Warnf("error closing file: %s", err)
		}
	}()

	s.log.Debugf("Uploading file '%s' to s3...", remoteFilename)

	out, err := s.uploader.Upload(&s3manager.UploadInput{
		Bucket: &s.bucket,
		Key:    &remoteFilename,
		Body:   file,
	})

	if err != nil {
		return nil, fmt.Errorf("unable to upload local file '%s' to bucket '%s', remote filename: '%s', error: '%w'", localFilename, s.bucket, remoteFilename, err)
	}
	s.log.Debugf("Upload successful, location: %s, ETag: %s", out.Location, *out.ETag)

	return &UploadResult{
		Filename: remoteFilename,
		Location: out.Location,
		ETag:     *out.ETag,
	}, nil
}

func getFilePath(dir, fileName string) string {
	return fmt.Sprintf("%s/%s", dir, fileName)
}

func verifyEnv(keys ...string) error {
	for _, currKey := range keys {
		if val := os.Getenv(currKey); len(val) == 0 {
			return fmt.Errorf("environment var %s not set", currKey)
		}
	}
	return nil
}
