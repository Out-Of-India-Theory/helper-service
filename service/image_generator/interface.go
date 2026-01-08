package image_generator

import "context"

type Service interface {
	GenerateImage(ctx context.Context, supplyId int) error
}
