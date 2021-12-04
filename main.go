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
	basePath := flag.String("base-path", "./data", "Directory that contains buckets")
	memMapFs := flag.Bool("mem-map-fs", false, "Use memory backed filesystem")
	bindAddr := flag.String("bind-addr", ":8080", "Use specified network interface")
	caCert := flag.String("ca-cert", "", "File that contains list of trusted SSL Certificate Authorities")
	clientCert := flag.String("client-cert", "", "File that contains X.509 certificate")
	clientKey := flag.String("client-key", "", "File that contains X.509 key")
	flag.Parse()

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	logger.Info("objects-service", zap.String("version", Version))

	tcpAddr, err := net.ResolveTCPAddr("tcp", *bindAddr)
	if err != nil {
		logger.Fatal(err.Error())
	}

	var tlsConfig *tls.Config

	if *caCert != "" {
		pool := x509.NewCertPool()

		data, err := os.ReadFile(*caCert)
		if err != nil {
			logger.Fatal(err.Error())
		}

		if !pool.AppendCertsFromPEM(data) {
			logger.Fatal("failed to parse CA certificate")
		}

		tlsConfig = &tls.Config{
			RootCAs:    pool,
			ClientAuth: tls.RequireAndVerifyClientCert,
			MinVersion: tls.VersionTLS12,
		}
	}

	srv := InitializeService(*basePath, *memMapFs, tcpAddr, tlsConfig, logger)
	srv.Start(*clientCert, *clientKey)
}
