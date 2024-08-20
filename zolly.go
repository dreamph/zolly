package main

import (
	"crypto"
	"crypto/tls"
	"fmt"
	"os"
	gopkcs12 "software.sslmate.com/src/go-pkcs12"

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

		app.All(configService.Path+"/*", balancerHandler)
	}

	welcomeInfo(config)

	if config.Server.SSL != nil && config.Server.SSL.Enable {
		tlsCert, err := initTLSConfig(config.Server.SSL.P12KeyFile, config.Server.SSL.P12KeyPassword)
		if err != nil {
			return err
		}
		ln, err := tls.Listen("tcp", ":"+config.Server.Port,
			&tls.Config{
				Certificates: []tls.Certificate{*tlsCert},
			},
		)
		if err != nil {
			return err
		}
		err = app.Listener(ln)
		if err != nil {
			return err
		}
	} else {
		err := app.Listen(":" + config.Server.Port)
		if err != nil {
			return err
		}
	}

	return nil
}

func initTLSConfig(path string, password string) (*tls.Certificate, error) {
	pkcs12Data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	key, cert, err := gopkcs12.Decode(pkcs12Data, password)
	if err != nil {
		return nil, err
	}

	tlsCert := tls.Certificate{
		Certificate: [][]byte{cert.Raw},
		PrivateKey:  key.(crypto.PrivateKey),
		Leaf:        cert,
	}

	return &tlsCert, nil
}

func welcomeInfo(config *GatewayConfig) {
	out := colorable.NewColorableStdout()
	mainLogo := logo
	mainLogo += fmt.Sprintf("Version : %s\n", CurrentVersion)
	if config.Server.SSL != nil && config.Server.SSL.Enable {
		mainLogo += fmt.Sprintf("Server : https://localhost:%s", config.Server.Port)
	} else {
		mainLogo += fmt.Sprintf("Server : http://localhost:%s", config.Server.Port)
	}
	_, _ = fmt.Fprintln(out, mainLogo)
}
