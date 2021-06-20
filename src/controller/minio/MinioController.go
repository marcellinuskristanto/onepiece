package minio

import (
	"context"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path"
	fp "path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"

	"github.com/marcellinuskristanto/onepiece/src/configuration"
	"github.com/marcellinuskristanto/onepiece/src/helper"
	s3Model "github.com/marcellinuskristanto/onepiece/src/model/request/s3"
)

// Upload file to bucket
func Upload(c *gin.Context) {
	ctx := context.Background()
	minioClient := helper.GetMinioInstance()
	config := configuration.GetConfig()

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

	f, contentType, err := helper.DownloadFileAndReturn(request.URLToUpload, request.Referer)
	if err != nil {
		res["message"] = fmt.Sprintf("failed to download url %q, %v", request.URLToUpload, err)
		c.JSON(500, res)
		return
	}
	fi, err := f.Stat()
	if err != nil {
		res["message"] = fmt.Sprintf("failed to get fileinfo %q, %v", request.URLToUpload, err)
		c.JSON(500, res)
		return
	}
	defer os.Remove(f.Name())
	defer f.Close()

	result, err := minioClient.PutObject(ctx, request.Bucket, path.Join(request.Filepath, request.Filename), f, fi.Size(), minio.PutObjectOptions{
		ContentType:  contentType,
		StorageClass: "REDUCED_REDUNDANCY",
	})

	if err != nil {
		res["message"] = fmt.Sprintf("failed to upload file, %v", err)
		c.JSON(500, res)
		return
	}
	res["success"] = true
	res["message"] = "File uploaded"
	res["location"] = helper.BuildLocation(config.App.MinioUrl, result.Bucket, result.Key)
	c.JSON(200, res)
}

// UploadFile file to bucket
func UploadFile(c *gin.Context) {
	ctx := context.Background()
	minioClient := helper.GetMinioInstance()
	config := configuration.GetConfig()

	res := gin.H{
		"success":  false,
		"message":  "Unknown error",
		"location": "",
	}
	form, _ := c.MultipartForm()
	files := form.File["file"]
	filename := c.PostForm("filename")
	filepath := c.PostForm("filepath")
	bucket := c.PostForm("bucket")
	// acl := "public-read"

	if bucket == "" {
		res["message"] = "Bucket required"
		c.JSON(http.StatusBadRequest, res)
		return
	}
	if len(files) <= 0 {
		res["message"] = "File required"
		c.JSON(http.StatusBadRequest, res)
		return
	}

	file := files[0]
	if filename == "" {
		filename = file.Filename
	}

	f, err := file.Open()
	if err != nil {
		res["message"] = fmt.Sprintf("failed reading file, %v", err)
		c.JSON(500, res)
		return
	}
	contentType, err := helper.GetFileContentType(f)
	if err != nil {
		ext := fp.Ext(filename)
		contentType = mime.TypeByExtension(ext)
	} else {
		if _, err = f.Seek(0, io.SeekStart); err != nil {
			res["message"] = fmt.Sprintf("failed rewind file pointer, %v", err)
			c.JSON(500, res)
			return
		}
	}

	result, err := minioClient.PutObject(ctx, bucket, path.Join(filepath, filename), f, file.Size, minio.PutObjectOptions{
		ContentType:  contentType,
		StorageClass: "REDUCED_REDUNDANCY",
	})
	if err != nil {
		res["message"] = fmt.Sprintf("failed to upload file, %v", err)
		c.JSON(500, res)
		return
	}

	res["success"] = true
	res["message"] = "File uploaded"
	res["location"] = helper.BuildLocation(config.App.MinioUrl, result.Bucket, result.Key)
	c.JSON(http.StatusOK, res)
}
