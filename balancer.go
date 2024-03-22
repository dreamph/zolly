package zolly

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
)

type Config struct {
	Servers         []string
	OnRequest       fiber.Handler
	OnResponse      fiber.Handler
	Timeout         time.Duration
	ReadBufferSize  int
	WriteBufferSize int
	TlsConfig       *tls.Config
	Client          *fasthttp.LBClient
	Path            string
	SkipPath        bool
}

var ConfigDefault = Config{
	OnRequest:  nil,
	OnResponse: nil,
	Timeout:    time.Second * 10,
}

func configDefault(config ...Config) Config {
	if len(config) < 1 {
		return ConfigDefault
	}

	cfg := config[0]

	if cfg.Timeout <= 0 {
		cfg.Timeout = ConfigDefault.Timeout
	}

	if len(cfg.Servers) == 0 && cfg.Client == nil {
		panic("Servers cannot be empty")
	}
	return cfg
}

func NewBalancer(config Config) (fiber.Handler, error) {
	cfg := configDefault(config)

	lbClient := &fasthttp.LBClient{}
	if config.Client == nil {
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

				ReadBufferSize:  config.ReadBufferSize,
				WriteBufferSize: config.WriteBufferSize,

				TLSConfig: config.TlsConfig,
			})
		}
	} else {
		lbClient = config.Client
	}

	return func(c *fiber.Ctx) error {
		req := c.Request()
		resp := c.Response()

		// Don't proxy "Connection" header
		req.Header.Del(fiber.HeaderConnection)

		if cfg.OnRequest != nil {
			if err := cfg.OnRequest(c); err != nil {
				return err
			}
		}

		if config.SkipPath {
			req.SetRequestURI(strings.ReplaceAll(utils.UnsafeString(req.RequestURI()), config.Path, Slash))
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
	}, nil
}
