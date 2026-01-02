package server

import (
	"context"
	"github.com/Out-Of-India-Theory/oit-go-commons/app"
	"github.com/Out-Of-India-Theory/supply-pn-image-generator/config"
	"github.com/Out-Of-India-Theory/supply-pn-image-generator/controller/image_generator"
	"github.com/Out-Of-India-Theory/supply-pn-image-generator/service/facade"
	"github.com/gin-gonic/gin"
	"net/http"
)

func registerRoutes(ctx context.Context, app *app.App, service facade.Service, configuration *config.Configuration) {
	basepath := app.Engine.Group("image_generator")
	app.Engine.GET("/health-check", HealthCheck)

	//pn-image_generator-geenrator
	{
		imageController := image_generator.InitImageGeneratorController(ctx, service, configuration)
		basepath.POST("/pn-image/:supply_id", imageController.GeneratePNImage)
	}
}

func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "health!",
	})
}
