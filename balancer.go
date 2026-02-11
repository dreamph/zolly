package main

import (
	"crypto/tls"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"github.com/valyala/fasthttp"
)

const defaultTimeout = 10 * time.Second

type Config struct {
	Servers   []string
	Timeout   time.Duration
	TlsConfig *tls.Config
	Path      string
	StripPath bool
}

func NewBalancer(config Config) (fiber.Handler, error) {
	client := &fasthttp.Client{
		ReadTimeout:  defaultTimeout,
		WriteTimeout: defaultTimeout,
		TLSConfig:    config.TlsConfig,
	}

	if config.Timeout > 0 {
		client.ReadTimeout = config.Timeout
		client.WriteTimeout = config.Timeout
	}

	roundRobin := &RoundRobin{
		Pool: config.Servers,
	}

	return func(c *fiber.Ctx) error {
		backend := roundRobin.Get()
		targetURL := backend + c.OriginalURL()
		if config.StripPath {
			targetURL = backend + strings.Replace(c.OriginalURL(), config.Path, "", 1)
		}

		return proxy.Do(c, targetURL, client)
	}, nil
}
