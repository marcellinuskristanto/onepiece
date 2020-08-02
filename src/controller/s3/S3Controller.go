package s3

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gin-gonic/gin"

	"github.com/marcellinuskristanto/onepiece/src/helper"
	s3Model "github.com/marcellinuskristanto/onepiece/src/model/request/s3"
)

// Upload file to bucket
func Upload(c *gin.Context) {
	res := gin.H{
		"success":  false,
		"message":  "Unknown error",
		"location": "",
	}

	request := s3Model.UploadRequest{
		ACL: "public-read",
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		res["message"] = err.Error()
		c.JSON(http.StatusBadRequest, res)
		return
	}

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(request.Region),
	}))

	uploader := s3manager.NewUploader(sess)

	f, contentType, err := helper.DownloadFileAndReturn(request.URLToUpload)
	if err != nil {
		res["message"] = fmt.Sprintf("failed to download url %q, %v", request.URLToUpload, err)
		c.JSON(500, res)
		return
	}
	defer os.Remove(f.Name())
	defer f.Close()

	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(request.Bucket),
		Key:         aws.String(path.Join(request.Filepath, request.Filename)),
		Body:        f,
		ACL:         &request.ACL,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		res["message"] = fmt.Sprintf("failed to upload file, %v", err)
		c.JSON(200, res)
		return
	}
	res["success"] = true
	res["message"] = "File uploaded"
	res["location"] = aws.StringValue(&result.Location)
	c.JSON(200, res)
}

// GetBucket get bucket
func GetBucket(c *gin.Context) {
	bucket := c.Query("bucket")
	region := c.Query("region")

	res := gin.H{
		"success": false,
		"message": "Unknown error",
	}

	sess := session.Must(session.NewSession())
	ctx := context.Background()

	region, err := s3manager.GetBucketRegion(ctx, sess, bucket, region)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == "NotFound" {
			res["message"] = fmt.Sprintf("unable to find bucket %s's region not found", bucket)
			c.JSON(200, res)
			return
		}
		c.JSON(200, res)
		return
	}
	res["success"] = true
	res["message"] = fmt.Sprintf("Bucket %s is in %s region\n", bucket, region)
	c.JSON(200, res)
	return
}

// CreateBucket post
func CreateBucket(c *gin.Context) {
	name := c.Query("name")
	region := c.Query("region")

	res := gin.H{
		"success": true,
		"message": "Bucket created",
	}

	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String(region),
		},
	})
	if err != nil {
		res["success"] = false
		res["message"] = "Create session failed"
		c.JSON(200, res)
		return
	}

	svc := s3.New(sess)
	input := &s3.CreateBucketInput{
		Bucket: aws.String(name),
	}

	_, err = svc.CreateBucket(input)
	if err != nil {
		res["success"] = false
		res["message"] = "Unknown error"

		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeBucketAlreadyExists:
				res["message"] = s3.ErrCodeBucketAlreadyExists
			case s3.ErrCodeBucketAlreadyOwnedByYou:
				res["message"] = s3.ErrCodeBucketAlreadyOwnedByYou
			default:
				res["message"] = aerr.Error()
			}
		} else {
			res["message"] = err.Error()
		}
	}
	c.JSON(200, res)
}
