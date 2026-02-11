package main

import (
	"crypto"
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dreamph/zolly/key"
	zollyplugin "github.com/dreamph/zolly/plugin"
	"github.com/dreamph/zolly/utils"
	"github.com/fatih/color"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	gopkcs12 "software.sslmate.com/src/go-pkcs12"
)

func Start(config *GatewayConfig) error {
	// Initialize plugin manager
	var pluginMgr *zollyplugin.Manager
	if len(config.Plugins) > 0 {
		pluginMgr = zollyplugin.NewManager()
		defer pluginMgr.Shutdown()

		for _, pd := range config.Plugins {
			err := pluginMgr.LoadPlugin(pd.Name, pd.Path, pd.Settings)
			if err != nil {
				return fmt.Errorf("failed to load plugin %q: %w", pd.Name, err)
			}
		}
	}

	app := fiber.New(fiber.Config{
		BodyLimit:             config.Server.BodyLimit,
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

		if pluginMgr != nil && len(configService.Plugins) > 0 {
			middlewares, err := zollyplugin.NewFiberMiddlewareChain(
				pluginMgr,
				configService.Plugins,
				time.Duration(configService.Timeout)*time.Second,
			)
			if err != nil {
				return fmt.Errorf("failed to create plugin chain for service %q: %w", configService.Path, err)
			}

			handlers := make([]fiber.Handler, 0, len(middlewares)+1)
			handlers = append(handlers, middlewares...)
			handlers = append(handlers, balancerHandler)
			app.Group(configService.Path, handlers...)
		} else {
			app.Group(configService.Path, balancerHandler)
		}
	}

	welcomeInfo(config)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-quit
		log.Println("Shutting down gateway...")
		if pluginMgr != nil {
			pluginMgr.Shutdown()
		}
		if err := app.Shutdown(); err != nil {
			log.Printf("Server shutdown error: %v\n", err)
		}
	}()

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

	if utils.FileExists(config.Server.SSL.GenerateKey.KeyConfig.File) {
		keyData, err := os.ReadFile(config.Server.SSL.GenerateKey.KeyConfig.File)
		if err != nil {
			return nil, "", err
		}
		return keyData, config.Server.SSL.GenerateKey.KeyConfig.Password, nil
	}

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
	color.Red(logo)
	color.Cyan(fmt.Sprintf("Version : %s\n", CurrentVersion))
	if config.Server.SSL != nil && config.Server.SSL.Enable {
		color.Cyan(fmt.Sprintf("Server : https://localhost:%s", config.Server.Port))
	} else {
		color.Cyan(fmt.Sprintf("Server : http://localhost:%s", config.Server.Port))
	}
}
