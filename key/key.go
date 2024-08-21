package key

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
	gopkcs12 "software.sslmate.com/src/go-pkcs12"
	"time"
)

const (
	DefaultKeySize = 2048 // NIST recommendation
	SerialNumSize  = 128
)

type CreateKeyRequest struct {
	Bits       int
	CommonName string
	Password   string
	ExpireDate time.Time
}

type CreateKeyResponse struct {
	SerialNumber *big.Int
	KeyBytes     []byte
}

func CreateKey(request *CreateKeyRequest) (*CreateKeyResponse, error) {
	keyBytes, err := rsa.GenerateKey(rand.Reader, request.Bits)
	if err != nil {
		return nil, err
	}

	err = keyBytes.Validate()
	if err != nil {
		return nil, err
	}

	subject := pkix.Name{
		CommonName: request.CommonName,
	}

	serialNumber, err := generateSerial(SerialNumSize)
	if err != nil {
		return nil, err
	}

	template := &x509.Certificate{
		SignatureAlgorithm:    x509.SHA256WithRSA,
		SerialNumber:          serialNumber,
		Subject:               subject,
		NotBefore:             time.Now(),
		NotAfter:              request.ExpireDate,
		BasicConstraintsValid: true,
		IsCA:                  false,
	}

	caCert := template
	caKey := keyBytes
	var caCerts []*x509.Certificate

	signedCert, err := x509.CreateCertificate(rand.Reader, template, caCert, &keyBytes.PublicKey, caKey)
	if err != nil {
		return nil, err
	}

	cert, err := x509.ParseCertificate(signedCert)
	if err != nil {
		return nil, err
	}

	p12Bytes, err := gopkcs12.Encode(rand.Reader, keyBytes, cert, caCerts, request.Password)
	if err != nil {
		return nil, err
	}

	_, _, _, err = gopkcs12.DecodeChain(p12Bytes, request.Password)
	if err != nil {
		return nil, err
	}

	var p12PublicBytes bytes.Buffer
	err = pem.Encode(&p12PublicBytes, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: signedCert,
	})
	if err != nil {
		return nil, err
	}

	return &CreateKeyResponse{
		KeyBytes:     p12Bytes,
		SerialNumber: serialNumber,
	}, nil
}

func generateSerial(bits int) (*big.Int, error) {
	var payload = make([]byte, bits/8)
	_, err := io.ReadFull(rand.Reader, payload)
	if err != nil {
		return nil, fmt.Errorf("read random source: %w", err)
	}
	v := new(big.Int)
	v.SetBytes(payload)
	return v, nil
}
