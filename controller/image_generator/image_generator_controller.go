package image_generator

import (
	"context"
	"github.com/Out-Of-India-Theory/helper-service/config"
	"github.com/Out-Of-India-Theory/helper-service/service/facade"
	"github.com/Out-Of-India-Theory/oit-go-commons/logging"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"
)

type Controller struct {
	logger  *zap.Logger
	service facade.Service
	config  *config.Configuration
}

func InitImageGeneratorController(ctx context.Context, service facade.Service, config *config.Configuration) *Controller {
	return &Controller{
		logger:  logging.WithContext(ctx),
		service: service,
		config:  config,
	}
}

func (con *Controller) GeneratePNImage(c *gin.Context) {
	supplyIdStr := c.Param("supply_id")
	supplyId, err := strconv.Atoi(supplyIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid supply_id"})
		return
	}
	go func() {
		bgCtx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()
		if err := con.service.ImageGeneratorService().GenerateImage(bgCtx, supplyId); err != nil {
			con.logger.Error("image generation failed", zap.Error(err))
		}
	}()
	c.JSON(http.StatusAccepted, gin.H{
		"status":  http.StatusAccepted,
		"message": "Image generation started",
	})
}
