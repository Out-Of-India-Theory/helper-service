package server

import (
	"context"
	"fmt"
	"github.com/Out-Of-India-Theory/oit-go-commons/app"
	"github.com/Out-Of-India-Theory/supply-pn-image-generator/config"
	"github.com/Out-Of-India-Theory/supply-pn-image-generator/service/facade"
	"github.com/Out-Of-India-Theory/supply-pn-image-generator/service/image_generator"
	"github.com/Out-Of-India-Theory/supply-pn-image-generator/service/image_uploader"
	"github.com/Out-Of-India-Theory/supply-pn-image-generator/service/supply"
	"github.com/gin-gonic/gin"
	"github.com/newrelic/go-agent/v3/integrations/nrgin"
	"github.com/newrelic/go-agent/v3/newrelic"
)

func InitServer(ctx context.Context, app *app.App, configuration *config.Configuration) {
	supplyService := supply.InitSupplyService(ctx, configuration)
	imageUploadService := image_uploader.InitImageUploader(ctx, configuration)
	imageService := image_generator.InitImageGeneratorService(ctx, supplyService, imageUploadService)
	facadeService := facade.InitFacadeService(ctx, imageService, imageUploadService, supplyService)
	registerMiddleware(app, configuration)
	registerRoutes(ctx, app, facadeService, configuration)

	app.StartHttpServer()
	err := app.StartMetricsServer()
	if err != nil {
		panic("Error while initializing http client")
	}

	<-make(chan int)
}

func registerMiddleware(app *app.App, configuration *config.Configuration) {
	newrelicApp, err := newrelic.NewApplication(
		newrelic.ConfigAppName(fmt.Sprintf("%s-%s", app.Config.AppName, app.Config.Env)),
		newrelic.ConfigLicense(app.Config.NewRelicLicense),
		newrelic.ConfigAppLogForwardingEnabled(true),
	)
	if err != nil {
		fmt.Println("Error while initializing new relic app")
		return
	}
	app.Engine.Use(nrgin.Middleware(newrelicApp))
	app.Engine.Use(newrelicTransactionMiddleware(newrelicApp))
	app.Engine.Use(CORSMiddleware())
}

func newrelicTransactionMiddleware(newRelicApp *newrelic.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, "newRelicTransaction", newrelic.FromContext(c))
		c.Request = c.Request.Clone(ctx)
		c.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, source")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
		c.Writer.Header().Set("source", "*")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
