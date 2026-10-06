package image_uploader

import (
	"context"
	"fmt"
	"github.com/Out-Of-India-Theory/helper-service/config"
	"github.com/Out-Of-India-Theory/oit-go-commons/client/http_client"
	"github.com/Out-Of-India-Theory/oit-go-commons/logging"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"time"
)

// uploadTimeout covers pushing a generated PNG to the platform upload API.
const uploadTimeout = time.Minute

type UploadApiResponse struct {
	Data    string `json:"data"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

type ImageUploader struct {
	logger        *zap.Logger
	configuration *config.Configuration
	httpClient    *http_client.HttpBaseClientV2
}

func InitImageUploader(ctx context.Context, configuration *config.Configuration) *ImageUploader {
	return &ImageUploader{
		logger:        logging.WithContext(ctx),
		configuration: configuration,
		httpClient:    http_client.NewHttpClientV2("", 0, uploadTimeout),
	}
}

func (s *ImageUploader) UploadToS3(ctx context.Context, fileName string, fileStream []byte) (string, error) {
	var apiURL string
	if s.configuration.ServerConfig.Env == "PROD" || s.configuration.ServerConfig.Env == "PRODUCTION" || strings.Contains(s.configuration.OMSClientConfig.Address, "production") {
		apiURL = fmt.Sprintf("%s/platform/document/v1/upload/prod_supply_pn_images?file_name=%s", s.configuration.OMSClientConfig.Address, fileName)
	} else {
		apiURL = fmt.Sprintf("%s/platform/document/v1/upload/jyotisha_pn_image?file_name=%s", s.configuration.OMSClientConfig.Address, fileName)
	}
	header := http.Header{}
	header.Set("Content-Type", "application/octet-stream")

	var response UploadApiResponse
	ex, err := s.httpClient.Call(ctx, apiURL, fileStream, header, http.MethodPost, &response)
	if err != nil {
		s.logger.Error("image upload failed", zap.String("file_name", fileName),
			zap.Int("status", ex.StatusCode), zap.String("body", ex.ResponseBody), zap.Error(err))
		return "", fmt.Errorf("UploadToS3: %w", err)
	}
	return response.Data, nil
}
