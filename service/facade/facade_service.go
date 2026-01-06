package facade

import (
	"context"
	"github.com/Out-Of-India-Theory/helper-service/config"
	"github.com/Out-Of-India-Theory/helper-service/service/image_generator"
	"github.com/Out-Of-India-Theory/helper-service/service/image_uploader"
	"github.com/Out-Of-India-Theory/helper-service/service/supply"
	"github.com/Out-Of-India-Theory/oit-go-commons/logging"
	"go.uber.org/zap"
)

type FacadeService struct {
	logger                *zap.Logger
	configuration         *config.Configuration
	imageGeneratorService image_generator.Service
	imageUploaderService  image_uploader.Service
	supplyService         supply.Service
}

func InitFacadeService(ctx context.Context, imageService image_generator.Service, imageUploaderService image_uploader.Service,
	supplyService supply.Service) *FacadeService {
	return &FacadeService{
		logger:                logging.WithContext(ctx),
		imageGeneratorService: imageService,
		imageUploaderService:  imageUploaderService,
		supplyService:         supplyService,
	}
}

func (s *FacadeService) ImageGeneratorService() image_generator.Service {
	return s.imageGeneratorService
}

func (s *FacadeService) ImageUploaderService() image_uploader.Service {
	return s.imageUploaderService
}

func (s *FacadeService) SupplyService() supply.Service {
	return s.supplyService
}
