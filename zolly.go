package zolly

import (
	"crypto/tls"
	"github.com/gofiber/fiber/v2"
	"time"
)

func Start(config *GatewayConfig) error {
	app := fiber.New()
	for _, s := range config.Services {
		balancerHandler, err := NewBalancer(Config{
			Path:     s.Path,
			SkipPath: s.StripPath,
			Timeout:  time.Duration(s.Timeout) * time.Second,
			Servers:  s.Servers,
			TlsConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		})
		if err != nil {
			return err
		}

		app.Group(s.Path, balancerHandler)
	}

	err := app.Listen(":" + config.Server.Port)
	if err != nil {
		return err
	}

	return nil
}
