package fileUpload

type FileUploadStorage struct {
	s3Bucket string
}

func NewFileStorage(s3Bucket string) *FileUploadStorage {
	return &FileUploadStorage{
		s3Bucket: s3Bucket,
	}
}
