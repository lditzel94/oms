package discovery

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"math/rand"
)

func ServiceConnection(ctx context.Context, serviceName string, registry Registry) (*grpc.ClientConn, error) {
	addresses, err := registry.Discover(ctx, serviceName)
	if err != nil {
		return nil, err
	}

	log.Printf("Discovered %d instances of %s", len(addresses), serviceName)

	// Randomly select an instance
	return grpc.Dial(
		addresses[rand.Intn(len(addresses))],
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		// Add OpenTelemetry interceptors
		//grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		//grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
	)
}
