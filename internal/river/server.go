package river

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
)

const (
	metricsPort = "9100"
)

// setupMetricsExporter sets up the OpenTelemetry Prometheus metrics exporter
func setupMetricsExporter() error {
	exporter, err := prometheus.New()
	if err != nil {
		return err
	}

	mp := metric.NewMeterProvider(
		metric.WithReader(exporter),
	)

	otel.SetMeterProvider(mp)

	return nil
}

// registerMetricsServer starts an HTTP server to expose Prometheus metrics
func registerMetricsServer(ctx context.Context) error {
	r := chi.NewRouter()

	// Expose metrics
	r.Get("/metrics", promhttp.Handler().ServeHTTP)

	// Setup HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", metricsPort),
		Handler:      r,
		ReadTimeout:  15 * time.Second, // nolint: mnd
		WriteTimeout: 15 * time.Second, // nolint: mnd
	}

	go func() {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // nolint: mnd

		defer cancel()

		if err := srv.Shutdown(seedContext(shutdownCtx)); err != nil {
			log.Error().Err(err).Msg("failed to shutdown metrics server")
		}
	}()

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

// SeedContext ensures the provided context carries a logger, returning a derived context when necessary.
func seedContext(ctx context.Context) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	logger := zerolog.Ctx(ctx)
	if logger == nil || logger.GetLevel() == zerolog.Disabled {
		logger = &log.Logger
	}

	return logger.WithContext(ctx)
}
