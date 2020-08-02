package s3

// UploadRequest model binding
type UploadRequest struct {
	Bucket      string `json:"bucket" binding:"required"`
	URLToUpload string `json:"urltoupload" binding:"required"`
	Filepath    string `json:"filepath" binding:"required"`
	Filename    string `json:"filename" binding:"required"`
	Region      string `json:"region" binding:"required"`
	ACL         string `json:"acl"`
}
