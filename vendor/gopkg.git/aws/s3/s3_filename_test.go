package s3

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func TestParseS3filename(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		expected *s3.GetObjectInput
	}{
		{
			name:     "HTTPS scheme",
			filename: "https://alpha-us-east-1.s3.amazonaws.com/temp-data/presal/9dd8bd54.csv.gz",
			expected: &s3.GetObjectInput{
				Bucket: aws.String("alpha-us-east-1"),
				Key:    aws.String("temp-data/presal/9dd8bd54.csv.gz"),
			},
		},
		{
			name:     "S3 scheme",
			filename: "s3://alpha-us-east-1/temp-data/presal/9dd8bd54.csv.gz",
			expected: &s3.GetObjectInput{
				Bucket: aws.String("alpha-us-east-1"),
				Key:    aws.String("temp-data/presal/9dd8bd54.csv.gz"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseS3filename(tt.filename)
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("parseS3filename() = %v, want %v", got, tt.expected)
			}
		})
	}
}
