package plugin

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
)

const DefaultPluginTimeout = 5 * time.Second

// NewFiberMiddleware creates a fiber.Handler that calls the given plugin for every request.
func NewFiberMiddleware(p MiddlewarePlugin, timeout time.Duration) fiber.Handler {
	if timeout == 0 {
		timeout = DefaultPluginTimeout
	}

	return func(c *fiber.Ctx) error {
		req := &Request{
			Method:      c.Method(),
			Path:        c.Path(),
			OriginalURL: c.OriginalURL(),
			Headers:     extractHeaders(c),
			Body:        c.Body(),
			QueryParams: extractQueryParams(c),
			ClientIP:    c.IP(),
		}

		ctx, cancel := context.WithTimeout(c.UserContext(), timeout)
		defer cancel()

		resp, err := p.HandleRequest(ctx, req)
		if err != nil {
			return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
				"error": "plugin error: " + err.Error(),
			})
		}

		switch resp.Action {
		case ActionAbort:
			for k, v := range resp.Headers {
				c.Set(k, v)
			}
			return c.Status(resp.StatusCode).Send(resp.Body)

		case ActionContinue:
			if resp.Headers != nil {
				for k, v := range resp.Headers {
					c.Request().Header.Set(k, v)
				}
			}
			if len(resp.Body) > 0 {
				c.Request().SetBody(resp.Body)
			}
			return c.Next()

		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "plugin returned unknown action: " + resp.Action,
			})
		}
	}
}

// NewFiberMiddlewareChain creates an ordered slice of fiber.Handler from plugin names.
func NewFiberMiddlewareChain(mgr *Manager, pluginNames []string, timeout time.Duration) ([]fiber.Handler, error) {
	handlers := make([]fiber.Handler, 0, len(pluginNames))
	for _, name := range pluginNames {
		p, err := mgr.GetPlugin(name)
		if err != nil {
			return nil, err
		}
		handlers = append(handlers, NewFiberMiddleware(p, timeout))
	}
	return handlers, nil
}

func extractHeaders(c *fiber.Ctx) map[string]string {
	headers := make(map[string]string)
	c.Request().Header.VisitAll(func(key, value []byte) {
		headers[string(key)] = string(value)
	})
	return headers
}

func extractQueryParams(c *fiber.Ctx) map[string]string {
	params := make(map[string]string)
	c.Request().URI().QueryArgs().VisitAll(func(key, value []byte) {
		params[string(key)] = string(value)
	})
	return params
}
