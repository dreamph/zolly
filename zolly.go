package main

import (
	"crypto/tls"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/mattn/go-colorable"
	"time"
)

func Start(config *GatewayConfig) error {
	app := fiber.New(fiber.Config{
		BodyLimit:             -1,
		DisableStartupMessage: true,
	})

	app.Use(cors.New())

	for _, configService := range config.Services {
		balancerHandler, err := NewBalancer(Config{
			Path:      configService.Path,
			StripPath: configService.StripPath,
			Timeout:   time.Duration(configService.Timeout) * time.Second,
			Servers:   configService.Servers,
			TlsConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		})
		if err != nil {
			return err
		}

		app.Group(configService.Path, balancerHandler)
	}

	out := colorable.NewColorableStdout()
	mainLogo := logo
	mainLogo += fmt.Sprintf("Version : %s\n", CurrentVersion)
	mainLogo += fmt.Sprintf("Server : http://127.0.0.1:%s", config.Server.Port)
	_, _ = fmt.Fprintln(out, mainLogo)

	err := app.Listen(":" + config.Server.Port)
	if err != nil {
		return err
	}

	return nil
}
