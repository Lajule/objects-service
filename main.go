package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"net"
	"os"

	"go.uber.org/zap"
)

// Version contains the program version.
var Version = "development"

func main() {
	basePath := flag.String("b", "./data", "Store base path")
	memory := flag.Bool("m", false, "Use memory backed filesystem")
	addr := flag.String("addr", ":8080", "TCP address")
	caCert := flag.String("ca-cert", "", "File that contains list of trusted SSL Certificate Authorities")
	clientCert := flag.String("client-cert", "", "File that contains X.509 certificate")
	clientKey := flag.String("client-key", "", "File that contains X.509 key")
	flag.Parse()

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	logger.Info("objects-service", zap.String("version", Version))

	tcpAddr, err := net.ResolveTCPAddr("tcp", *addr)
	if err != nil {
		logger.Fatal(err.Error())
	}

	var tlsConfig *tls.Config
	var pool *x509.CertPool

	if *clientCert != "" && *clientKey != "" {
		if *caCert != "" {
			pool = x509.NewCertPool()

			data, err := os.ReadFile(*caCert)
			if err != nil {
				logger.Fatal(err.Error())
			}

			if !pool.AppendCertsFromPEM(data) {
				logger.Fatal("failed to parse CA certificate")
			}
		}

		cert, err := tls.LoadX509KeyPair(*clientCert, *clientKey)
		if err != nil {
			logger.Fatal(err.Error())
		}

		tlsConfig = &tls.Config{
			RootCAs:      pool,
			ClientAuth:   tls.RequireAndVerifyClientCert,
			MinVersion:   tls.VersionTLS12,
			Certificates: []tls.Certificate{cert},
		}
	}

	srv := InitializeService(*basePath, *memory, tcpAddr, tlsConfig, logger)
	srv.Start()
}
