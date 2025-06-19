package telemetry

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"log"
)

type Jaeger struct {
	provider *tracesdk.TracerProvider
	tracer   trace.Tracer
}

func NewJaeger(ctx context.Context, serviceName string) (*Jaeger, error) {
	var tp *tracesdk.TracerProvider
	tp, err := createJaegerTraceProvider(ctx, serviceName)
	if err != nil {
		return nil, err
	}
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	tracer := tp.Tracer(serviceName)

	return &Jaeger{
		provider: tp,
		tracer:   tracer,
	}, nil
}

func (ot *Jaeger) Start(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, Span) {
	if len(opts) == 0 {
		return ot.tracer.Start(ctx, name)
	}
	return ot.tracer.Start(ctx, name, opts...)
}

func (ot *Jaeger) Shutdown(ctx context.Context) error {
	return ot.provider.Shutdown(ctx)
}

func createJaegerTraceProvider(ctx context.Context, serviceName string) (*tracesdk.TracerProvider, error) {
	res, err := resource.New(ctx)
	if err != nil {
		return nil, err
	}

	// Configure Jaeger exporter
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(
		jaeger.WithEndpoint("http://localhost:14268/api/traces"),
	))
	if err != nil {
		log.Fatalf("Failed to create Jaeger exporter: %v", err)
	}

	// Configure trace provider
	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exporter),
		tracesdk.WithResource(res),
		tracesdk.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(serviceName)),
		),
	)

	return tp, nil
}
