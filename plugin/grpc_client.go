package plugin

import (
	"context"

	pb "github.com/dreamph/zolly/plugin/proto"
)

// GRPCClient is the gateway-side implementation of MiddlewarePlugin.
type GRPCClient struct {
	client pb.MiddlewarePluginClient
}

func (c *GRPCClient) HandleRequest(ctx context.Context, req *Request) (*Response, error) {
	pbResp, err := c.client.HandleRequest(ctx, &pb.HttpRequest{
		Method:      req.Method,
		Path:        req.Path,
		OriginalUrl: req.OriginalURL,
		Headers:     req.Headers,
		Body:        req.Body,
		QueryParams: req.QueryParams,
		ClientIp:    req.ClientIP,
	})
	if err != nil {
		return nil, err
	}

	return &Response{
		Action:     pbResp.Action,
		StatusCode: int(pbResp.StatusCode),
		Headers:    pbResp.Headers,
		Body:       pbResp.Body,
	}, nil
}

func (c *GRPCClient) Configure(ctx context.Context, settings *PluginSettings) (*ConfigureResult, error) {
	pbResp, err := c.client.Configure(ctx, &pb.PluginConfig{
		Settings: settings.Settings,
	})
	if err != nil {
		return nil, err
	}
	return &ConfigureResult{
		Success:      pbResp.Success,
		ErrorMessage: pbResp.ErrorMessage,
	}, nil
}

func (c *GRPCClient) HealthCheck(ctx context.Context) (*HealthCheckResult, error) {
	pbResp, err := c.client.HealthCheck(ctx, &pb.Empty{})
	if err != nil {
		return nil, err
	}
	return &HealthCheckResult{
		Healthy: pbResp.Healthy,
		Message: pbResp.Message,
	}, nil
}
