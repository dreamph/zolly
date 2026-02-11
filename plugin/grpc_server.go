package plugin

import (
	"context"

	pb "github.com/dreamph/zolly/plugin/proto"
)

// GRPCServer wraps a MiddlewarePlugin implementation as a gRPC server.
// This is used inside plugin binaries.
type GRPCServer struct {
	pb.UnimplementedMiddlewarePluginServer
	Impl MiddlewarePlugin
}

func (s *GRPCServer) HandleRequest(ctx context.Context, req *pb.HttpRequest) (*pb.HttpResponse, error) {
	resp, err := s.Impl.HandleRequest(ctx, &Request{
		Method:      req.Method,
		Path:        req.Path,
		OriginalURL: req.OriginalUrl,
		Headers:     req.Headers,
		Body:        req.Body,
		QueryParams: req.QueryParams,
		ClientIP:    req.ClientIp,
	})
	if err != nil {
		return nil, err
	}
	return &pb.HttpResponse{
		Action:     resp.Action,
		StatusCode: int32(resp.StatusCode),
		Headers:    resp.Headers,
		Body:       resp.Body,
	}, nil
}

func (s *GRPCServer) Configure(ctx context.Context, req *pb.PluginConfig) (*pb.ConfigureResponse, error) {
	result, err := s.Impl.Configure(ctx, &PluginSettings{Settings: req.Settings})
	if err != nil {
		return nil, err
	}
	return &pb.ConfigureResponse{
		Success:      result.Success,
		ErrorMessage: result.ErrorMessage,
	}, nil
}

func (s *GRPCServer) HealthCheck(ctx context.Context, _ *pb.Empty) (*pb.HealthCheckResponse, error) {
	result, err := s.Impl.HealthCheck(ctx)
	if err != nil {
		return nil, err
	}
	return &pb.HealthCheckResponse{
		Healthy: result.Healthy,
		Message: result.Message,
	}, nil
}
