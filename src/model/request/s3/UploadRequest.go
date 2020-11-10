package s3

// UploadRequest model binding
type UploadRequest struct {
	Bucket      string `json:"bucket"`
	URLToUpload string `json:"urltoupload" binding:"required"`
	Filepath    string `json:"filepath" binding:"required"`
	Filename    string `json:"filename" binding:"required"`
	Region      string `json:"region"`
	ACL         string `json:"acl"`
	Referer     string `json:"referer"`
}
