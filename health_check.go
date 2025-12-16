package health_check

import (
	"connectrpc.com/connect"
	"context"
	"errors"
	"fmt"
	healthv1 "rahulxf.com/general-connectrpc-demo/internal/gen/go/health/v1"
	"sync"
)

type CheckRequest struct {
	Service string
}
type Status uint8

type CheckResponse struct {
	Status Status
}
type Checker interface {
	Check(context.Context, *CheckRequest) (*CheckResponse, error)
}

const (
	// StatusUnknown indicates that the service's health state is indeterminate.
	StatusUnknown Status = 0

	// StatusServing indicates that the service is ready to accept requests.
	StatusServing Status = 1

	// StatusNotServing indicates that the process is healthy but the service is
	// not accepting requests. For example, StatusNotServing is often appropriate
	// when your primary database is down or unreachable.
	StatusNotServing Status = 2
)

// String representation of the status.
func (s Status) String() string {
	switch s {
	case StatusUnknown:
		return "unknown"
	case StatusServing:
		return "serving"
	case StatusNotServing:
		return "not_serving"
	}

	return fmt.Sprintf("status_%d", s)
}

func NewHandler(checker Checker, options ...connect.HandlerOption) {
	const serviceName = "/health.v1.Health/"
	connect.NewUnaryHandler(serviceName+"Check",
		func(ctx context.Context, req *connect.Request[healthv1.HealthCheckRequest]) (*connect.Response[healthv1.HealthCheckResponse], error) {
			var checkRequest CheckRequest
			if req.Msg != nil {
				checkRequest.Service = req.Msg.GetService()
			}
			checkResponse, err := checker.Check(ctx, &checkRequest)
			if err != nil {
				return nil, err
			}

			return connect.NewResponse(&healthv1.HealthCheckResponse{
				Status: healthv1.HealthCheckResponse(),
			}), nil
		})
	watch := connect.NewServerStreamHandler(
		serviceName+"Watch",
		func(
			_ context.Context,
			_ *connect.Request[healthv1.HealthCheckRequest],
			_ *connect.ServerStream[healthv1.HealthCheckResponse],
		) error {
			return connect.NewError(
				connect.CodeUnimplemented,
				errors.New("connect doesn't support watching health state"),
			)
		},
		options...,
	)
	mux.Handle(serviceName+"Watch", watch)
	return serviceName, mux
}

type StaticChecker struct {
	mu       sync.RWMutex
	statuses map[string]Status
}

func (c *StaticChecker) Check(_ context.Context, req *CheckRequest) (*CheckResponse, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if status, registered := c.statuses[req.Service]; registered {
		return &CheckResponse{Status: status}, nil
	}
	if req.Service == "" {
		return &CheckResponse{Status: StatusServing}, nil
	}
	return nil, connect.NewError(
		connect.CodeNotFound,
		fmt.Errorf("unknown service %s", req.Service),
	)
}
