package supply

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Out-Of-India-Theory/oit-go-commons/logging"
	"github.com/Out-Of-India-Theory/supply-pn-image-generator/config"
	"go.uber.org/zap"
	"io"
	"net/http"
	"time"
)

type SupplyService struct {
	logger *zap.Logger
	config *config.Configuration
}

func InitSupplyService(ctx context.Context, config *config.Configuration) *SupplyService {
	return &SupplyService{
		logger: logging.WithContext(ctx),
		config: config,
	}
}

func (s *SupplyService) GetSupplyDetails(ctx context.Context, supplyId int) (*SupplyResponse, error) {
	url := fmt.Sprintf("%s/internal/supply/supplies/%d", s.config.SupplyClientConfig.Address, supplyId)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		s.logger.Error("failed to create supply request", zap.Error(err))
		return nil, err
	}

	req.Header.Set("Accept", "application/json")

	client := &http.Client{
		Timeout: 10 * time.Minute,
	}

	resp, err := client.Do(req)
	if err != nil {
		s.logger.Error("supply service request failed", zap.Error(err), zap.Int("supply_id", supplyId))
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Error("failed to read supply response", zap.Error(err))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		s.logger.Error("unexpected supply service status", zap.Int("status", resp.StatusCode), zap.ByteString("body", body))
		return nil, fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	var result SupplyResponse
	if err := json.Unmarshal(body, &result); err != nil {
		s.logger.Error("failed to unmarshal supply response", zap.Error(err), zap.ByteString("body", body))
		return nil, err
	}

	return &result, nil
}
