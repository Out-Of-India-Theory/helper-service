package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Out-Of-India-Theory/oit-go-commons/app"
	"github.com/Out-Of-India-Theory/oit-go-commons/logging"
	"go.uber.org/zap"
)

const (
	// httpDrainTimeout plus newRelicFlushTimeout stays inside the default 30s
	// container stop grace period.
	httpDrainTimeout     = 15 * time.Second
	newRelicFlushTimeout = 10 * time.Second
)

// waitForShutdown blocks until SIGINT/SIGTERM, then drains in-flight HTTP requests,
// stops the metrics server and flushes the final New Relic harvest. Background image
// generation jobs are not waited for; one in progress is cut off.
func waitForShutdown(ctx context.Context, app *app.App) {
	logger := logging.WithContext(ctx)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	sig := <-stop
	logger.Info("shutdown signal received", zap.String("signal", sig.String()))

	drainCtx, cancel := context.WithTimeout(context.Background(), httpDrainTimeout)
	defer cancel()
	if err := app.StopHttpServer(drainCtx); err != nil {
		logger.Error("http server shutdown incomplete", zap.Error(err))
	}
	if app.Config.PrometheusEnabled {
		if err := app.StopMetricsServer(drainCtx); err != nil {
			logger.Error("metrics server shutdown incomplete", zap.Error(err))
		}
	}

	app.StopNewRelic(newRelicFlushTimeout)
	logger.Info("shutdown complete")
}
