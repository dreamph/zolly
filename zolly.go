package main

import (
	"crypto/tls"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/mattn/go-colorable"
	"time"
)

func Start(config *GatewayConfig) error {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})
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

	out := colorable.NewColorableStdout()
	mainLogo := logo
	mainLogo += fmt.Sprintf("Server : http://127.0.0.1:%s", config.Server.Port)
	_, _ = fmt.Fprintln(out, mainLogo)

	err := app.Listen(":" + config.Server.Port)
	if err != nil {
		return err
	}

	return nil
}
