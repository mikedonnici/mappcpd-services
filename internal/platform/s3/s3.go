package s3

import (
	"os"
	"time"

	"github.com/34South/envr"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func init() {
	envr.New("mappcpd-attachments", []string{
		"AWS_ACCESS_KEY_ID",
		"AWS_SECRET_ACCESS_KEY",
		"AWS_REGION",
	}).Auto()
}

// PutRequest issues a signed URL that allows for a PUT to an Amazon S3 bucket. It receives the
// key (full path to file including file name', and the name of the bucket.
// The aws package ASSUMES the presence of AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY env vars so they have been
// added to init() above. AWS_REGION was added by me.
func PutRequest(key, bucket string) (string, error) {

	sess := session.Must(session.NewSession())
	svc := s3.New(sess, aws.NewConfig().WithRegion(os.Getenv("AWS_REGION")))
	req, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		//Body:   strings.NewReader("EXPECTED CONTENTS"),
	})

	return req.Presign(15 * time.Minute)
}
