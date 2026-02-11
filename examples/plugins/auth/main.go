package main

import (
	"context"
	"log"

	zollyplugin "github.com/dreamph/zolly/plugin"
)

type AuthPlugin struct {
	apiKey     string
	headerName string
}

func (p *AuthPlugin) Configure(_ context.Context, settings *zollyplugin.PluginSettings) (*zollyplugin.ConfigureResult, error) {
	p.apiKey = settings.Settings["api_key"]
	p.headerName = settings.Settings["header_name"]
	if p.headerName == "" {
		p.headerName = "X-API-Key"
	}
	if p.apiKey == "" {
		return &zollyplugin.ConfigureResult{
			Success:      false,
			ErrorMessage: "api_key setting is required",
		}, nil
	}
	log.Printf("[auth-plugin] Configured with header=%s\n", p.headerName)
	return &zollyplugin.ConfigureResult{Success: true}, nil
}

func (p *AuthPlugin) HandleRequest(_ context.Context, req *zollyplugin.Request) (*zollyplugin.Response, error) {
	providedKey := req.Headers[p.headerName]
	if providedKey != p.apiKey {
		return &zollyplugin.Response{
			Action:     zollyplugin.ActionAbort,
			StatusCode: 401,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       []byte(`{"error":"unauthorized","message":"invalid or missing API key"}`),
		}, nil
	}

	return &zollyplugin.Response{
		Action: zollyplugin.ActionContinue,
	}, nil
}

func (p *AuthPlugin) HealthCheck(_ context.Context) (*zollyplugin.HealthCheckResult, error) {
	return &zollyplugin.HealthCheckResult{Healthy: true, Message: "ok"}, nil
}

func main() {
	zollyplugin.Serve(&AuthPlugin{})
}
