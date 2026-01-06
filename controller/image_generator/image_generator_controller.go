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
	ctx := c.Request.Context()
	supplyIdStr := c.Param("supply_id")
	supplyId, err := strconv.Atoi(supplyIdStr)
	err = con.service.ImageGeneratorService().GenerateImage(ctx, supplyId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    nil,
		"status":  http.StatusOK,
		"message": "Successful",
	})
}
