package image_uploader

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Out-Of-India-Theory/helper-service/config"
	"github.com/Out-Of-India-Theory/oit-go-commons/logging"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
)

type UploadApiResponse struct {
	Data    string `json:"data"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

type ImageUploader struct {
	logger        *zap.Logger
	configuration *config.Configuration
}

func InitImageUploader(ctx context.Context, configuration *config.Configuration) *ImageUploader {
	return &ImageUploader{
		logger:        logging.WithContext(ctx),
		configuration: configuration,
	}
}

func (s *ImageUploader) UploadToS3(ctx context.Context, fileName string, fileStream []byte) (string, error) {
	var apiURL string
	if s.configuration.ServerConfig.Env == "PROD" || s.configuration.ServerConfig.Env == "PRODUCTION" {
		apiURL = fmt.Sprintf("%s/platform/document/v1/upload/prod_supply_pn_images?file_name=%s", s.configuration.SupplyClientConfig.Address, fileName)
	} else {
		apiURL = fmt.Sprintf("%s/platform/document/v1/upload/jyotisha_pn_image?file_name=%s", s.configuration.SupplyClientConfig.Address, fileName)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, bytes.NewReader(fileStream))
	if err != nil {
		s.logger.Error("failed to create upload request", zap.Error(err))
		return "", err
	}

	req.Header.Set("Content-Type", "application/octet-stream")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := ioutil.ReadAll(resp.Body)
		return "", fmt.Errorf("upload failed with status code %d: %s", resp.StatusCode, string(body))
	}

	var Response UploadApiResponse
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read API response: %w", err)
	}
	err = json.Unmarshal(body, &Response)
	if err != nil {
		return "", fmt.Errorf("failed to parse API response: %w", err)
	}
	return Response.Data, nil
}
