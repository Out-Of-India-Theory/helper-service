package image_uploader

import "context"

type Service interface {
	UploadToS3(ctx context.Context, fileName string, fileStream []byte) (string, error)
}
