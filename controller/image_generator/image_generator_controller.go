package image_generator

import (
	"context"
	"github.com/Out-Of-India-Theory/helper-service/config"
	"github.com/Out-Of-India-Theory/helper-service/service/facade"
	"github.com/Out-Of-India-Theory/oit-go-commons/logging"
	"github.com/Out-Of-India-Theory/oit-go-commons/util"
	"github.com/gin-gonic/gin"
	"github.com/newrelic/go-agent/v3/newrelic"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Controller struct {
	logger  *zap.Logger
	service facade.Service
	config  *config.Configuration
	// jobs tracks background image generations so shutdown can wait for them.
	jobs *sync.WaitGroup
}

func InitImageGeneratorController(ctx context.Context, service facade.Service, config *config.Configuration, jobs *sync.WaitGroup) *Controller {
	return &Controller{
		logger:  logging.WithContext(ctx),
		service: service,
		config:  config,
		jobs:    jobs,
	}
}

func (con *Controller) GeneratePNImage(c *gin.Context) {
	supplyIdStr := c.Param("supply_id")
	supplyId, err := strconv.Atoi(supplyIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid supply_id"})
		return
	}
	// The request transaction ends when this handler returns, so the job gets its
	// own background transaction for its HTTP segments to attach to.
	nrApp := newrelic.FromContext(c).Application()
	con.jobs.Add(1)
	util.SafeGo(c.Request.Context(), "generate-pn-image", func() {
		defer con.jobs.Done()
		bgCtx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()
		var txn *newrelic.Transaction
		if nrApp != nil {
			txn = nrApp.StartTransaction("GeneratePNImage")
			defer txn.End()
			bgCtx = newrelic.NewContext(bgCtx, txn)
		}
		if err := con.service.ImageGeneratorService().GenerateImage(bgCtx, supplyId); err != nil {
			txn.NoticeError(err)
			con.logger.Error("image generation failed", zap.Int("supply_id", supplyId), zap.Error(err))
		}
	})
	c.JSON(http.StatusAccepted, gin.H{
		"status":  http.StatusAccepted,
		"message": "Image generation started",
	})
}
