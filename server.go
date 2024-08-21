package main

import (
	"crypto"
	"crypto/tls"
	"fmt"
	"github.com/dreamph/zolly/key"
	"github.com/dreamph/zolly/utils"
	"github.com/gofiber/fiber/v2/middleware/logger"
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

	if config.Server.Cors != nil && config.Server.Cors.Enable {
		app.Use(cors.New())
	}

	if config.Server.Log != nil && config.Server.Log.Enable {
		app.Use(logger.New())
	}

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

	welcomeInfo(config)

	if config.Server.SSL != nil && config.Server.SSL.Enable {
		keyBytes, password, err := initKey(config)
		if err != nil {
			return err
		}

		tlsCert, err := initTLSConfig(keyBytes, password)
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

func initKey(config *GatewayConfig) ([]byte, string, error) {
	if config.Server.SSL.GenerateKey == nil || !config.Server.SSL.GenerateKey.Enable {
		keyData, err := os.ReadFile(config.Server.SSL.Key.File)
		if err != nil {
			return nil, "", err
		}

		return keyData, config.Server.SSL.Key.Password, nil
	}

	if !utils.FileExists(config.Server.SSL.GenerateKey.KeyConfig.File) {
		keyResponse, err := key.CreateKey(&key.CreateKeyRequest{
			Bits:       key.DefaultKeySize,
			CommonName: config.Server.SSL.GenerateKey.KeyConfig.CommonName,
			Password:   config.Server.SSL.GenerateKey.KeyConfig.Password,
			ExpireDate: time.Now().AddDate(1, 0, 0),
		})
		if err != nil {
			return nil, "", err
		}

		err = utils.WriteFile(config.Server.SSL.GenerateKey.KeyConfig.File, keyResponse.KeyBytes)
		if err != nil {
			return nil, "", err
		}

		return keyResponse.KeyBytes, config.Server.SSL.GenerateKey.KeyConfig.Password, nil
	} else {
		keyData, err := os.ReadFile(config.Server.SSL.GenerateKey.KeyConfig.File)
		if err != nil {
			return nil, "", err
		}

		return keyData, config.Server.SSL.GenerateKey.KeyConfig.Password, nil
	}
}

func initTLSConfig(pkcs12Data []byte, password string) (*tls.Certificate, error) {
	privateKey, cert, err := gopkcs12.Decode(pkcs12Data, password)
	if err != nil {
		return nil, err
	}

	tlsCert := tls.Certificate{
		Certificate: [][]byte{cert.Raw},
		PrivateKey:  privateKey.(crypto.PrivateKey),
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
