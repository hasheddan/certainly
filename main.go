package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"time"
)

func main() {
	org := os.Args[1]
	countArg := os.Args[2]
	count, err := strconv.ParseInt(countArg, 10, 64)
	if err != nil {
		panic(err)
	}
	rootTmpl := &x509.Certificate{
		SerialNumber: big.NewInt(2023),
		Subject: pkix.Name{
			Organization:  []string{org},
			Country:       []string{"US"},
			Province:      []string{""},
			Locality:      []string{"Cloud"},
			StreetAddress: []string{"IP"},
			PostalCode:    []string{"00000"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}

	ca, err := x509.CreateCertificate(rand.Reader, rootTmpl, rootTmpl, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}

	caPEM := new(bytes.Buffer)
	if err := pem.Encode(caPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: ca,
	}); err != nil {
		panic(err)
	}

	if err := os.WriteFile("root.crt.pem", caPEM.Bytes(), 0o666); err != nil {
		panic(err)
	}

	for i := 0; i < int(count); i++ {
		devTmpl := &x509.Certificate{
			SerialNumber: big.NewInt(2023),
			Subject: pkix.Name{
				Organization:  []string{org},
				CommonName:    fmt.Sprintf("device-%d", i),
				Country:       []string{"US"},
				Province:      []string{""},
				Locality:      []string{"Cloud"},
				StreetAddress: []string{"IP"},
				PostalCode:    []string{"00000"},
			},
			NotBefore:             time.Now(),
			NotAfter:              time.Now().AddDate(10, 0, 0),
			IsCA:                  true,
			ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
			KeyUsage:              x509.KeyUsageDigitalSignature,
			BasicConstraintsValid: true,
		}

		devKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			panic(err)
		}

		devKeyBytes, err := x509.MarshalECPrivateKey(devKey)
		if err != nil {
			panic(err)
		}

		devKeyPEM := new(bytes.Buffer)
		if err := pem.Encode(devKeyPEM, &pem.Block{
			Type:  "EC PRIVATE KEY",
			Bytes: devKeyBytes,
		}); err != nil {
			panic(err)
		}

		if err := os.WriteFile(fmt.Sprintf("dev-%d.key.pem", i), devKeyPEM.Bytes(), 0o666); err != nil {
			panic(err)
		}

		devCert, err := x509.CreateCertificate(rand.Reader, devTmpl, rootTmpl, &devKey.PublicKey, key)
		if err != nil {
			panic(err)
		}

		devPEM := new(bytes.Buffer)
		if err := pem.Encode(devPEM, &pem.Block{
			Type:  "CERTIFICATE",
			Bytes: devCert,
		}); err != nil {
			panic(err)
		}

		if err := os.WriteFile(fmt.Sprintf("dev-%d.crt.pem", i), devPEM.Bytes(), 0o666); err != nil {
			panic(err)
		}

	}
}
