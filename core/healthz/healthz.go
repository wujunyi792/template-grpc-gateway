package healthz

import (
	"context"
	"google.golang.org/grpc/codes"
	health "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
	"pinnacle-primary-be/core/syncx"
)

var Health = syncx.ForAtomicBool(false)

type S struct{}

func (s S) Check(ctx context.Context, in *health.HealthCheckRequest) (*health.HealthCheckResponse, error) {
	if Health.True() {
		return &health.HealthCheckResponse{Status: health.HealthCheckResponse_SERVING}, nil
	}
	return &health.HealthCheckResponse{Status: health.HealthCheckResponse_NOT_SERVING}, nil
}

func (s S) Watch(in *health.HealthCheckRequest, _ health.Health_WatchServer) error {
	// Example of how to register both methods but only implement the Check method.
	return status.Error(codes.Unimplemented, "unimplemented")
}
