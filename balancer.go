package main

import (
	"crypto/tls"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/valyala/fasthttp"
	"net/url"
	"strings"
	"time"
)

const (
	Slash = "/"
	Empty = ""
)

type Config struct {
	Servers    []string
	OnRequest  fiber.Handler
	OnResponse fiber.Handler
	Timeout    time.Duration
	//ReadBufferSize  int
	//WriteBufferSize int
	TlsConfig *tls.Config
	Client    *fasthttp.LBClient
	Path      string
	StripPath bool
}

func configDefault(config Config) Config {
	if config.Timeout <= 0 {
		config.Timeout = time.Second * 10
	}

	if len(config.Servers) == 0 && config.Client == nil {
		panic("Servers cannot be empty")
	}
	return config
}

func NewBalancer(config Config) (fiber.Handler, error) {
	cfg := configDefault(config)

	lbClient, err := intiLBClient(cfg)
	if err != nil {
		return nil, err
	}

	return func(c *fiber.Ctx) error {
		return handleRequest(c, cfg, lbClient)
	}, nil
}

func handleRequest(c *fiber.Ctx, cfg Config, lbClient *fasthttp.LBClient) error {
	req := c.Request()
	resp := c.Response()

	// Don't proxy "Connection" header
	req.Header.Del(fiber.HeaderConnection)

	if cfg.OnRequest != nil {
		if err := cfg.OnRequest(c); err != nil {
			return err
		}
	}

	if cfg.StripPath {
		req.SetRequestURI(strings.Replace(utils.UnsafeString(req.RequestURI()), cfg.Path, Empty, 1))
	} else {
		req.SetRequestURI(utils.UnsafeString(req.RequestURI()))
	}

	err := lbClient.Do(req, resp)
	if err != nil {
		return err
	}

	resp.Header.Del(fiber.HeaderConnection)

	if cfg.OnResponse != nil {
		if err := cfg.OnResponse(c); err != nil {
			return err
		}
	}

	return nil
}

func intiLBClient(cfg Config) (*fasthttp.LBClient, error) {
	lbClient := &fasthttp.LBClient{}
	if cfg.Client == nil {
		lbClient.Timeout = cfg.Timeout
		for _, server := range cfg.Servers {
			if !strings.HasPrefix(server, "http") {
				server = "http://" + server
			}

			u, err := url.Parse(server)
			if err != nil {
				return nil, err
			}

			lbClient.Clients = append(lbClient.Clients, &fasthttp.HostClient{
				NoDefaultUserAgentHeader: true,
				DisablePathNormalizing:   true,
				Addr:                     u.Host,
				TLSConfig:                cfg.TlsConfig,
			})
		}
	} else {
		lbClient = cfg.Client
	}

	return lbClient, nil
}
