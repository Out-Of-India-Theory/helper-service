package facade

import (
	"github.com/Out-Of-India-Theory/supply-pn-image-generator/service/image_generator"
	"github.com/Out-Of-India-Theory/supply-pn-image-generator/service/image_uploader"
	"github.com/Out-Of-India-Theory/supply-pn-image-generator/service/supply"
)

type Service interface {
	ImageGeneratorService() image_generator.Service
	ImageUploaderService() image_uploader.Service
	SupplyService() supply.Service
}
