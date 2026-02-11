package plugin

import (
	"context"

	goplugin "github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"

	pb "github.com/dreamph/zolly/plugin/proto"
)

const (
	PluginName = "middleware"

	ActionContinue = "continue"
	ActionAbort    = "abort"
)

var Handshake = goplugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "ZOLLY_PLUGIN",
	MagicCookieValue: "zolly-middleware-v1",
}

var PluginMap = map[string]goplugin.Plugin{
	PluginName: &MiddlewareGRPCPlugin{},
}

// Request represents an incoming HTTP request passed to a plugin.
type Request struct {
	Method      string
	Path        string
	OriginalURL string
	Headers     map[string]string
	Body        []byte
	QueryParams map[string]string
	ClientIP    string
}

// Response represents the plugin's decision.
type Response struct {
	Action     string
	StatusCode int
	Headers    map[string]string
	Body       []byte
}

// PluginSettings holds arbitrary config for a plugin.
type PluginSettings struct {
	Settings map[string]string
}

// ConfigureResult is returned from Configure.
type ConfigureResult struct {
	Success      bool
	ErrorMessage string
}

// HealthCheckResult is returned from HealthCheck.
type HealthCheckResult struct {
	Healthy bool
	Message string
}

// MiddlewarePlugin is the interface that all middleware plugins must implement.
type MiddlewarePlugin interface {
	HandleRequest(ctx context.Context, req *Request) (*Response, error)
	Configure(ctx context.Context, settings *PluginSettings) (*ConfigureResult, error)
	HealthCheck(ctx context.Context) (*HealthCheckResult, error)
}

// MiddlewareGRPCPlugin implements goplugin.GRPCPlugin.
type MiddlewareGRPCPlugin struct {
	goplugin.Plugin
	Impl MiddlewarePlugin
}

func (p *MiddlewareGRPCPlugin) GRPCServer(broker *goplugin.GRPCBroker, s *grpc.Server) error {
	pb.RegisterMiddlewarePluginServer(s, &GRPCServer{Impl: p.Impl})
	return nil
}

func (p *MiddlewareGRPCPlugin) GRPCClient(ctx context.Context, broker *goplugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCClient{client: pb.NewMiddlewarePluginClient(c)}, nil
}

// Serve starts the plugin process. Plugin authors call this from main().
func Serve(impl MiddlewarePlugin) {
	goplugin.Serve(&goplugin.ServeConfig{
		HandshakeConfig: Handshake,
		Plugins: map[string]goplugin.Plugin{
			PluginName: &MiddlewareGRPCPlugin{
				Impl: impl,
			},
		},
		GRPCServer: goplugin.DefaultGRPCServer,
	})
}
