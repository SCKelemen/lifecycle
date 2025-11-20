package lifecycle

import (
	"context"
	"strings"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

// OTelIntegration provides OpenTelemetry integration for lifecycle events
type OTelIntegration struct {
	tracer  trace.Tracer
	meter   metric.Meter
	counter map[string]metric.Int64Counter
	histogram map[string]metric.Float64Histogram
}

// NewOTelIntegration creates a new OpenTelemetry integration
func NewOTelIntegration(serviceName string) *OTelIntegration {
	tracer := otel.Tracer("lifecycle")
	meter := otel.Meter("lifecycle")

	return &OTelIntegration{
		tracer:   tracer,
		meter:    meter,
		counter:  make(map[string]metric.Int64Counter),
		histogram: make(map[string]metric.Float64Histogram),
	}
}

// StartSpan starts an OpenTelemetry span for an event
func (o *OTelIntegration) StartSpan(ctx context.Context, eventType string, attrs ...attribute.KeyValue) (context.Context, trace.Span) {
	spanName := o.getSpanName(eventType)
	ctx, span := o.tracer.Start(ctx, spanName, trace.WithAttributes(attrs...))
	return ctx, span
}

// RecordMetric records a metric for an event
func (o *OTelIntegration) RecordMetric(ctx context.Context, eventType string, duration time.Duration, attrs ...attribute.KeyValue) {
	// Record counter
	counterName := o.getCounterName(eventType)
	counter, ok := o.counter[counterName]
	if !ok {
		var err error
		counter, err = o.meter.Int64Counter(counterName, metric.WithDescription("Count of "+eventType+" events"))
		if err == nil {
			o.counter[counterName] = counter
		}
	}
	if counter != nil {
		counter.Add(ctx, 1, metric.WithAttributes(attrs...))
	}

	// Record duration histogram for timed events
	if duration > 0 {
		histogramName := o.getHistogramName(eventType)
		histogram, ok := o.histogram[histogramName]
		if !ok {
			var err error
			histogram, err = o.meter.Float64Histogram(histogramName, metric.WithDescription("Duration of "+eventType+" events"))
			if err == nil {
				o.histogram[histogramName] = histogram
			}
		}
		if histogram != nil {
			histogram.Record(ctx, duration.Seconds(), metric.WithAttributes(attrs...))
		}
	}
}

// RecordValue records a value metric (for gauges or histograms)
func (o *OTelIntegration) RecordValue(ctx context.Context, metricName string, value float64, attrs ...attribute.KeyValue) {
	histogram, ok := o.histogram[metricName]
	if !ok {
		var err error
		histogram, err = o.meter.Float64Histogram(metricName)
		if err == nil {
			o.histogram[metricName] = histogram
		}
	}
	if histogram != nil {
		histogram.Record(ctx, value, metric.WithAttributes(attrs...))
	}
}

// getSpanName converts event type to span name
func (o *OTelIntegration) getSpanName(eventType string) string {
	// Convert event type to span name
	// e.g., "api.request.received" -> "api.request"
	parts := splitEventType(eventType)
	if len(parts) >= 2 {
		return parts[0] + "." + parts[1]
	}
	return eventType
}

// getCounterName converts event type to counter name
func (o *OTelIntegration) getCounterName(eventType string) string {
	return eventType + ".count"
}

// getHistogramName converts event type to histogram name
func (o *OTelIntegration) getHistogramName(eventType string) string {
	return eventType + ".duration"
}

// splitEventType splits an event type into parts
func splitEventType(eventType string) []string {
	return strings.Split(eventType, ".")
}

// EventAttributes converts event data to OpenTelemetry attributes
func EventAttributes(event Event) []attribute.KeyValue {
	attrs := []attribute.KeyValue{
		attribute.String("event.type", event.GetEventType()),
		attribute.String("service.name", event.GetService()),
		attribute.String("service.instance.id", event.GetHost()),
	}

	if correlationID := event.GetCorrelationID(); correlationID != "" {
		attrs = append(attrs, attribute.String("correlation.id", correlationID))
	}

	return attrs
}

