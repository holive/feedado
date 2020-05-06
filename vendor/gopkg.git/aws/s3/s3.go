package s3

import (
	"context"
	"io"

	"github.com/pkg/errors"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// S3Client to interact with a S3 gz files saving much memory as possible
type S3Client struct {
	session *session.Session
}

// Open S3 files stream, filename can be written as follows:
// - "https://${bucket}/${key}",
// - "s3://${bucket}/${key}"
func (sc *S3Client) Open(ctx context.Context, filename string) (io.ReadCloser, error) {
	goi, err := parseS3filename(filename)
	if err != nil {
		return nil, err
	}
	s3Client := s3.New(sc.session)

	output, err := s3Client.GetObjectWithContext(ctx, goi)
	if err != nil {
		return nil, err
	}

	return output.Body, nil
}

// OpenDownload a file from S3.
func (sc *S3Client) OpenDownload(ctx context.Context, filename string) (io.ReadCloser, error) {
	goi, err := parseS3filename(filename)
	if err != nil {
		return nil, err
	}

	buf, err := NewTempFileBuffer()
	if err != nil {
		return nil, errors.Wrap(err, "could not create buffer")
	}

	downloader := s3manager.NewDownloader(sc.session)
	_, err = downloader.DownloadWithContext(ctx, buf, goi)
	if err != nil {
		buf.Close()
		return nil, err
	}

	return buf, nil
}

// NewS3Client returns a new S3Client.
func NewS3Client(sess *session.Session) *S3Client {
	return &S3Client{session: sess}

}
