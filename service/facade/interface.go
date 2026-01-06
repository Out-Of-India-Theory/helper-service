package facade

import (
	"github.com/Out-Of-India-Theory/helper-service/service/image_generator"
	"github.com/Out-Of-India-Theory/helper-service/service/image_uploader"
	"github.com/Out-Of-India-Theory/helper-service/service/supply"
)

type Service interface {
	ImageGeneratorService() image_generator.Service
	ImageUploaderService() image_uploader.Service
	SupplyService() supply.Service
}
