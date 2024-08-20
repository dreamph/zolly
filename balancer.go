package main

import (
	"crypto/tls"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"github.com/valyala/fasthttp"
	"strings"
	"time"
)

const (
	Empty = ""
)

type Config struct {
	Servers   []string
	Timeout   time.Duration
	TlsConfig *tls.Config
	Path      string
	StripPath bool
}

func NewBalancer(config Config) (fiber.Handler, error) {
	//maxIdleConnDuration := time.Hour * 1
	//tcpDialer := fasthttp.TCPDialer{
	//	Concurrency:      4096,
	//	DNSCacheDuration: time.Hour,
	//}
	client := &fasthttp.Client{
		ReadTimeout:  time.Duration(10) * time.Second,
		WriteTimeout: time.Duration(10) * time.Second,
		TLSConfig:    config.TlsConfig,

		//MaxIdleConnDuration:           maxIdleConnDuration,
		//NoDefaultUserAgentHeader:      true, // Don't send: User-Agent: fasthttp
		//DisableHeaderNamesNormalizing: true, // If you set the case on your headers correctly you can enable this
		//DisablePathNormalizing:        true,
		//Dial:                          tcpDialer.Dial,
	}

	if config.Timeout > 0 {
		client.ReadTimeout = config.Timeout
		client.WriteTimeout = config.Timeout
	}

	roundRobin := &RoundRobin{
		Current: 0,
		Pool:    config.Servers,
	}

	return func(c *fiber.Ctx) error {
		backend := roundRobin.Get()
		targetURL := backend + c.OriginalURL()
		if config.StripPath {
			targetURL = backend + strings.Replace(c.OriginalURL(), config.Path, Empty, 1)
		}

		return proxy.Do(c, targetURL, client)
	}, nil
}
