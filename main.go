package main

import (
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"github.com/caarlos0/env/v6"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"os"
	"github.com/pavel-v-chernykh/keystore-go/v4"
	"software.sslmate.com/src/go-pkcs12"
	"time"
)

type config struct {
	Password          string `env:"PASSWORD" envDefault:"password"`
	Mode              bool   `env:"FILE_MODE" envDefault:false`
	Key               string `env:"KAFKA_KEY"`
	CaCertificate     string `env:"KAFKA_CA"`
	Certificate       string `env:"KAFKA_CERT"`
	KeyFile           string `env:"KAFKA_KEY_FILE,file"`
	CertificateFile   string `env:"KAFKA_CERT_FILE,file"`
	CaCertificateFile string `env:"KAFKA_CA_FILE"`
	OutputP12         string `env:"OUTPUT_P12" envDefault:"/var/run/secrets/truststore.p12"`
	OutputJKS         string `env:"OUTPUT_JKS" envDefault:"/var/run/secrets/truststore.jks"`
}

func main() {
	zerolog.DurationFieldUnit = time.Second
	if err := run(); err != nil {
		log.Fatal().Err(err).Msg("Failed to run")
	}
	log.Info().Msg("Gracefully exiting")
}

func run() error {
	cfg := config{}

	if err := env.Parse(&cfg); err != nil {
		return err
	}

	var key, certificate, cacertificate string

	if cfg.Mode == true {
		key = cfg.KeyFile
		certificate = cfg.CertificateFile
		cacertificate = cfg.CaCertificateFile
	} else {
		key = cfg.Key
		certificate = cfg.Certificate
		cacertificate = cfg.CaCertificate
	}

	pemPrivateKey, err := readPem("PRIVATE KEY", key)
	if err != nil {
		return err
	}

	pemCertificate, err := readPem("CERTIFICATE", certificate)
	if err != nil {
		return err
	}

	pemCaCertificate, err := readPem("CERTIFICATE", cacertificate)
	if err != nil {
		return err
	}

	crt, err := x509.ParseCertificate(pemCertificate)
	if err != nil {
		panic(err)
	}

	priKey, err := x509.ParsePKCS8PrivateKey(pemPrivateKey)
	if err != nil {
		return err
	}

	// create p12
	pfxBytes, err := pkcs12.Encode(rand.Reader, priKey, crt, nil, cfg.Password)
	if err != nil {
		return err
	}

	// create keystore
	ks := keystore.New()

	tce := keystore.TrustedCertificateEntry{
		CreationTime: time.Now(),
		Certificate: keystore.Certificate{
			Type:    "X509",
			Content: pemCaCertificate,
		},
	}

	if err := ks.SetTrustedCertificateEntry("alias", tce); err != nil {
		return err
	}

	if err := writeKeyStore(ks, cfg.OutputJKS, []byte(cfg.Password)); err != nil {
		return err
	}


	// validate output
	_, _, _, err = pkcs12.DecodeChain(pfxBytes, cfg.Password)
	if err != nil {
		return (err)
	}

	// write output
	if err := ioutil.WriteFile(
		cfg.OutputP12,
		pfxBytes,
		os.ModePerm,
	); err != nil {
		return (err)
	}

	return nil
}

func writeKeyStore(ks keystore.KeyStore, filename string, password []byte) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}

	if err = ks.Store(f, password); err != nil {
		f.Close()
		return err
	}

	f.Close()
	return nil
}

func readPem(expectedType string, data string) ([]byte, error) {
	b, _ := pem.Decode([]byte(data))
	if b == nil {
		return nil, errors.New("should have at least one pem block")
	}

	if b.Type != expectedType {
		return nil, errors.New("should be a " + expectedType)
	}

	return b.Bytes, nil
}
