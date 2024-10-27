package metadata

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func InjectMetadataIntoContext(ctx context.Context, metadata map[string]string) context.Context {
	propagator := otel.GetTextMapPropagator()

	return propagator.Extract(
		ctx,
		propagation.MapCarrier(metadata),
	)
}

func ExtractMetadataFromContext(ctx context.Context) map[string]string {
	propagator := otel.GetTextMapPropagator()

	metadata := map[string]string{}
	propagator.Inject(
		ctx,
		propagation.MapCarrier(metadata),
	)

	return metadata
}
