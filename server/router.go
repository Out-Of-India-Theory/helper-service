package server

import (
	"context"
	"github.com/Out-Of-India-Theory/helper-service/config"
	"github.com/Out-Of-India-Theory/helper-service/controller/image_generator"
	"github.com/Out-Of-India-Theory/helper-service/service/facade"
	"github.com/Out-Of-India-Theory/oit-go-commons/app"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
)

func registerRoutes(ctx context.Context, app *app.App, service facade.Service, configuration *config.Configuration, imageJobs *sync.WaitGroup) {
	basepath := app.Engine.Group("helper-service")
	app.Engine.GET("/health-check", HealthCheck)
	basepath.GET("/health-check", HealthCheck)

	//pn-image_generator-geenrator
	{
		imageController := image_generator.InitImageGeneratorController(ctx, service, configuration, imageJobs)
		basepath.POST("/pn-image/:supply_id", imageController.GeneratePNImage)
	}
}

func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "health!",
	})
}
