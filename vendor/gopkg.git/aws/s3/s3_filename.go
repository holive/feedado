package s3

import (
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pkg/errors"
)

func parseS3filename(v string) (*s3.GetObjectInput, error) {
	u, err := url.Parse(v)
	if err != nil {
		return nil, err
	}

	var bucket string

	switch u.Scheme {
	case "https":
		bucket = strings.Split(u.Host, ".")[0]
	case "s3":
		bucket = u.Host
	default:
		return nil, errors.New("unknown scheme")
	}

	return &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(u.Path[1:]),
	}, nil
}
