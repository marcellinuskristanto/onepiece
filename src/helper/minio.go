package helper

import (
	"fmt"

	minio "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var minioClient *minio.Client

func InitMinio(endpoint string, user string, secret string) error {
	// Initialize minio client object.
	var err error
	minioClient, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(user, secret, ""),
		Secure: true,
	})
	return err
}

func GetMinioInstance() *minio.Client {
	return minioClient
}

func BuildLocation(domain string, bucket string, key string) string {
	return fmt.Sprintf("https://%s/%s/%s", domain, bucket, key)
}
