package lib

import (
	"context"
	"log/slog"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

func InitTracer() (*sdktrace.TracerProvider, error) {
	ctx := context.Background()
	exporterEnabled := os.Getenv("OTEL_EXPORTER_ENABLED")

	var exporter sdktrace.SpanExporter
	var err error

	if exporterEnabled != "" {
		exporter, err = otlptracegrpc.New(ctx, otlptracegrpc.WithInsecure())
		if err != nil {
			return nil, err
		}
	}

	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			"",
			semconv.ServiceName("server"),
		),
	)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to create resource", "error", err)
		return nil, err
	}

	opts := []sdktrace.TracerProviderOption{
		sdktrace.WithResource(r),
	}

	if exporter != nil {
		opts = append(opts, sdktrace.WithBatcher(exporter))
	}

	tp := sdktrace.NewTracerProvider(opts...)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return tp, nil
}
