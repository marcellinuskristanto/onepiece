package route

import (
	"github.com/gin-gonic/gin"
	"github.com/marcellinuskristanto/onepiece/src/controller/minio"
	"github.com/marcellinuskristanto/onepiece/src/controller/s3"
)

// LoadRoute load all routing
func LoadRoute(r *gin.Engine) {
	api := r.Group("/api")
	v1 := api.Group("/v1")
	s3route := v1.Group("/s3")
	{
		s3route.POST("/upload", s3.Upload)
		s3route.POST("/uploadFile", s3.UploadFile)
		s3route.GET("/bucket", s3.GetBucket)
		s3route.POST("/bucket", s3.CreateBucket)
	}
	minioroute := v1.Group("/minio")
	{
		minioroute.POST("/upload", minio.Upload)
		minioroute.POST("/uploadFile", minio.UploadFile)
	}
}
