package supply

import (
	"context"
	"fmt"
	"github.com/Out-Of-India-Theory/helper-service/config"
	"github.com/Out-Of-India-Theory/oit-go-commons/client/http_client"
	"github.com/Out-Of-India-Theory/oit-go-commons/logging"
	"go.uber.org/zap"
	"net/http"
)

type SupplyService struct {
	logger     *zap.Logger
	config     *config.Configuration
	httpClient *http_client.HttpBaseClientV2
}

func InitSupplyService(ctx context.Context, config *config.Configuration) *SupplyService {
	return &SupplyService{
		logger:     logging.WithContext(ctx),
		config:     config,
		httpClient: http_client.NewHttpClientV2("", 0, config.OMSClientConfig.Timeout),
	}
}

func (s *SupplyService) GetSupplyDetails(ctx context.Context, supplyId int) (*SupplyResponse, error) {
	url := fmt.Sprintf("%s/internal/supply/supplies/%d", s.config.OMSClientConfig.Address, supplyId)

	header := http.Header{}
	header.Set("Accept", "application/json")

	var result SupplyResponse
	ex, err := s.httpClient.Call(ctx, url, nil, header, http.MethodGet, &result)
	if err != nil {
		s.logger.Error("supply service request failed", zap.Int("supply_id", supplyId),
			zap.Int("status", ex.StatusCode), zap.String("body", ex.ResponseBody), zap.Error(err))
		return nil, fmt.Errorf("GetSupplyDetails: %w", err)
	}

	return &result, nil
}
