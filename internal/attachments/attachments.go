package attachments

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// S3PutRequest issues a signed URL that allows for a PUT to an Amazon S3 bucket. It receives the
// key (full path to file including file name', and the name of the bucket.
func S3PutRequest(key, bucket string) (string, error) {

	sess := session.Must(session.NewSession())
	svc := s3.New(sess, aws.NewConfig().WithRegion("ap-southeast-2"))
	req, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(key),
		Key:    aws.String(bucket),
		//Body:   strings.NewReader("EXPECTED CONTENTS"),
	})

	return req.Presign(15 * time.Minute)
}
