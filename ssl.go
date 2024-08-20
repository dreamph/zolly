package main

import (
	"crypto"
	"crypto/tls"
	"os"
	gopkcs12 "software.sslmate.com/src/go-pkcs12"
)

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
